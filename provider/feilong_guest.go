/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"context"
	"strings"
	"strconv"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Bischoff/feilong-client-go"
)

func feilongGuest() *schema.Resource {
	return &schema.Resource{
		Description:	"Feilong guest VM resource",

		CreateContext:	feilongGuestCreate,
		ReadContext:	feilongGuestRead,
		UpdateContext:	feilongGuestUpdate,
		DeleteContext:	feilongGuestDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:	"Arbitrary name for the resource",
				Type:		schema.TypeString,
				Required:	true,
			},
			"userid": {
				Description:	"System name for z/VM",
				Type:		schema.TypeString,
				Optional:	true,
				Computed:	true,
			},
			"vcpus": {
				Description:	"Virtual CPUs count",
				Type:		schema.TypeInt,
				Optional:	true,
				Default:	1,
			},
			"memory": {
				Description:	"Memory size with unit (G, M, K, B)",
				Type:		schema.TypeString,
				Optional:	true,
				Default:	"512M",
			},
			"disk": {
				Description:	"Disk size of first disk with unit (T, G, M, K, B)",
				Type:		schema.TypeString,
				Optional:	true,
				Default:	"10G",
				// we could use the size of the image file instead
			},
			"image": {
				Description:	"Image name",
				Type:		schema.TypeString,
				Required:	true,
			},
			"adapter_address": {
				Description:	"Desired virtual device of first interface",
				Type:		schema.TypeString,
				Optional:	true,
				Default:	"1000",
			},
			"mac": {
				Description:	"Desired MAC address of first interface",
				Type:		schema.TypeString,
				Optional:	true,
				Default:	"",
			},
			"vswitch": {
				Description:	"Name of virtual switch to connect to",
				Type:		schema.TypeString,
				Optional:	true,
				Default:	"DEVNET",
			},
			"cloudinit_params": {
				Description:	"Path to cloud-init parameters file",
				Type:		schema.TypeString,
				Optional:	true,
			},
			"mac_address": {
				Description:	"MAC address of first interface after deployment",
				Type:		schema.TypeString,
				Computed:	true,
			},
			"ip_address": {
				Description:	"IP address of first interface after deployment",
				Type:		schema.TypeString,
				Computed:	true,
			},
		},
	}
}

func feilongGuestCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Compute computed fields
	resourceName := d.Get("name").(string)
	userid := d.Get("userid").(string)
	if userid == "" {
		userid = strings.ToUpper(resourceName)
		if (len(userid) > 8) {
			userid = userid[:8]
		}
		d.Set("userid", userid)
	}

	// Compute values passed to Feilong API but not part of the data model
	size := d.Get("disk").(string)
	vcpus := d.Get("vcpus").(int)
	memory, err := convertToMegabytes(d.Get("memory").(string))
	if err != nil {
		return diag.Errorf("Conversion Error: %s", err)
	}
	image := d.Get("image").(string)
	adapterAddress := d.Get("adapter_address").(string)
	mac := d.Get("mac").(string)
	vswitch := d.Get("vswitch").(string)
	cloudinitParams := d.Get("cloudinit_params").(string)
	localUser := meta.(*apiClient).LocalUser

	// Create the guest
	client := meta.(*apiClient).Client
	diskList := []feilong.GuestDisk {
		{
			Size:		size,
			IsBootDisk:	true,
		},
	}
	createParams := feilong.CreateGuestParams {
		UserId:		userid,
		VCPUs:		vcpus,
		Memory:		memory,
		DiskList:	diskList,
	}
	_, err = client.CreateGuest(&createParams)
	if err != nil {
		return diag.Errorf("Creation Error: %s", err)
	}

	// Create the first network interface
	createNICParams := feilong.CreateGuestNICParams {
		VDev:           adapterAddress,
		MACAddress:	mac,
	}
	err = client.CreateGuestNIC(userid, &createNICParams)
	if err != nil {
		return diag.Errorf("NIC Creation Error: %s", err)
	}

	// Couple the first network interface to the virtual switch
	updateNICParams := feilong.UpdateGuestNICParams {
		Couple:		true,
		VSwitch:	vswitch,
	}
	err = client.UpdateGuestNIC(userid, adapterAddress, &updateNICParams)
	if err != nil {
		return diag.Errorf("NIC Coupling Error: %s", err)
	}

	// Deploy the guest
	deployParams := feilong.DeployGuestParams {
		Image:		image,
		TransportFiles:	cloudinitParams,
		RemoteHost:	localUser,
	}
	err = client.DeployGuest(userid, &deployParams)
	if err != nil {
		return diag.Errorf("Deployment Error: %s", err)
	}

	// Start the guest
	err = client.StartGuest(userid)
	if err != nil {
		return diag.Errorf("Startup Error: %s", err)
	}

	// Wait until the guest gets an IP address
	var macAddress string
	var ipAddress string
	err = waitForLease(ctx, client, userid, &macAddress, &ipAddress)
	if err != nil {
		return diag.Errorf("Error Waiting for an IP Address: %s", err)
	}
	err = d.Set("mac_address", macAddress)
	if err != nil {
		return diag.Errorf("MAC Address Registration Error: %s", err)
	}
	err = d.Set("ip_address", ipAddress)
	if err != nil {
		return diag.Errorf("IP Address Registration Error: %s", err)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a Feilong guest resource")

	// Set resource identifier
	d.SetId(resourceName)

	return nil
}

func feilongGuestRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiClient).Client
	userid := d.Get("userid").(string)

	// Obtain info about this guest
	guestInfo, err := client.GetGuestInfo(userid)
	if err != nil {
		return diag.Errorf("Guest Querying Error: %s", err)
	}

	// Read number of vCPUs
	err = d.Set("vcpus", guestInfo.Output.NumCPUs)
	if err != nil {
		return diag.Errorf("Virtual CPUs Setting Error: %s", err)
	}

	// Read memory
	declaredMemory, err := convertToMegabytes(d.Get("memory").(string))
	if err != nil {
		return diag.Errorf("Conversion Error: %s", err)
	}
	obtainedMemory := guestInfo.Output.MaxMemKB / 1_024
	if declaredMemory == obtainedMemory {
		// do not overwrite memory if value equal but a different unit
		tflog.Info(ctx, "Not replacing memory size " + d.Get("memory").(string) + " with equal value " + strconv.Itoa(obtainedMemory) + "M")
	} else {
		err = d.Set("memory", strconv.Itoa(obtainedMemory) + "M")
		if err != nil {
			return diag.Errorf("Memory Setting Error: %s", err)
		}
	}

	// Obtain first minidisk info
	minidisksInfo, err := client.GetGuestMinidisksInfo(userid)
	if err != nil {
		return diag.Errorf("Minidisks Querying Error: %s", err)
	}
	if len(minidisksInfo.Output.Minidisks) < 1 {
		return diag.Errorf("Minidisk Not Found Error")
	}
	firstMinidisk := minidisksInfo.Output.Minidisks[0]

	// Read disk size
	declaredDiskSize, err := convertToMegabytes(d.Get("disk").(string))
	if err != nil {
		return diag.Errorf("Conversion Error: %s", err)
	}
	if firstMinidisk.DeviceUnits != "Cylinders" {
		return diag.Errorf("Unknown Minidisk Unit Error: %s", firstMinidisk.DeviceUnits)
	}
	// tracks/cylinder=15  blocks/track=12  kilobytes/block=4  15*12*4=720
	obtainedDiskSize := (firstMinidisk.DeviceSize * 720) / 1_024
	if declaredDiskSize == obtainedDiskSize {
		// do not overwrite disk size if value equal but a different unit
		tflog.Info(ctx, "Not replacing disk size " + d.Get("disk").(string) + " with equal value " + strconv.Itoa(obtainedDiskSize) + "M")
	} else {
		err = d.Set("disk", strconv.Itoa(obtainedDiskSize) + "M")
		if err != nil {
			return diag.Errorf("Disk Setting Error: %s", err)
		}
	}

	// Obtain first network adapter info
	adaptersInfo, err := client.GetGuestAdaptersInfo(userid)
	if err != nil {
		return diag.Errorf("Network Adapter Info Querying Error: %s", err)
	}
	if len(adaptersInfo.Output.Adapters) < 1 {
		return diag.Errorf("Network Adapter Not Found Error")
	}
	firstAdapter := adaptersInfo.Output.Adapters[0]

	// Read virtual switch name
	err = d.Set("vswitch", firstAdapter.LANName)
	if err != nil {
		return diag.Errorf("VSwitch Setting Error: %s", err)
	}

	// TODO: read adapter virtual device address

	// Read MAC address
	declaredMAC := d.Get("mac").(string)
	obtainedMAC := firstAdapter.MACAddress
	if strings.ToLower(declaredMAC[8:]) == strings.ToLower(obtainedMAC[8:]) {
		// do not overwrite a MAC address if last 3 hex bytes are the same
		tflog.Info(ctx, "Not replacing MAC address " + declaredMAC + " with other MAC address with same last 3 hex bytes " + obtainedMAC)
	} else {
		err = d.Set("mac", obtainedMAC)
		if err != nil {
			return diag.Errorf("MAC Address Setting Error: %s", err)
		}
	}

	// Read IP address
	declaredIPAddress := d.Get("ip_address").(string)
	obtainedIPAddress := firstAdapter.IPAddress
	if firstAdapter.IPVersion == "6" && strings.HasPrefix(obtainedIPAddress, "fe80:") {
		// do not overwrite an IPv4 address with a link-local IPv6 address
		tflog.Info(ctx, "Not replacing IP address " + declaredIPAddress + " with an IPv6 link-local address " + obtainedIPAddress)
	} else {
		err = d.Set("ip_address", obtainedIPAddress)
		if err != nil {
			return diag.Errorf("IP Address Setting Error: %s", err)
		}
	}

	// CAVEATS:
	//  - the image used during the deployment cannot be determined after the deployment
	//  - the cloud init image used during the deployment cannot be determined after the deployment

	return nil
}

func feilongGuestUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiClient).Client
	userid := d.Get("userid").(string)

	// Address userid changes
	if d.HasChange("userid") {
		oldValue, _ := d.GetChange("userid")
		oldUserid := oldValue.(string)
		d.Set("userid", oldUserid)
		return diag.Errorf("Cannot change userid \"%s\" to \"%s\"", oldUserid, userid)
	}

	// Address vCPUs changes
	if d.HasChange("vcpus") {
		oldValue, newValue := d.GetChange("vcpus")
		oldvCPUs := oldValue.(int)
		newvCPUs := newValue.(int)
		if newvCPUs <= oldvCPUs {
			d.Set("vcpus", oldvCPUs)
			return diag.Errorf("Cannot decrease number of vCPUs from %d to %d", oldvCPUs, newvCPUs)
		}
		liveResizeCPUsParams := feilong.LiveResizeGuestCPUsParams {
			CPUCount: newvCPUs,
		}
		err := client.LiveResizeGuestCPUs(userid, &liveResizeCPUsParams)
		if err != nil {
			d.Set("vcpus", oldvCPUs)
			return diag.Errorf("CPUs resizing error: %s", err)
		}
		tflog.Info(ctx, "Increased number of vCPUs from " + strconv.Itoa(oldvCPUs) + " to " + strconv.Itoa(newvCPUs))
	}

	// Address memory changes
	if d.HasChanges("memory") {
		oldValue, newValue := d.GetChange("memory")
		oldMemoryMB, err := convertToMegabytes(oldValue.(string))
		if err != nil {
			d.Set("memory", oldValue.(string))
			return diag.Errorf("Conversion Error: %s", err)
		}
		newMemoryMB, err := convertToMegabytes(newValue.(string))
		if err != nil {
			d.Set("memory", oldValue.(string))
			return diag.Errorf("Conversion Error: %s", err)
		}
		if newMemoryMB < oldMemoryMB {
			d.Set("memory", oldValue.(string))
			return diag.Errorf("Cannot decrease memory size from %s to %s", oldValue.(string), newValue.(string))
		} else if newMemoryMB > oldMemoryMB {
			liveResizeMemoryParams := feilong.LiveResizeGuestMemoryParams {
				Size: newValue.(string),
			}
			err = client.LiveResizeGuestMemory(userid, &liveResizeMemoryParams)
			if err != nil {
				d.Set("memory", oldValue.(string))
				return diag.Errorf("Memory resizing error: %s", err)
			}
			tflog.Info(ctx, "Increased memory size from " + oldValue.(string) + " to " + newValue.(string))
		} else {
			tflog.Info(ctx, "Not replacing memory size " + oldValue.(string) + " with equal value " + newValue.(string))
		}
	}

	// Address main disk size changes
	if d.HasChanges("disk") {
		oldValue, newValue := d.GetChange("disk")
		oldDiskMB, err := convertToMegabytes(oldValue.(string))
		if err != nil {
			d.Set("disk", oldValue.(string))
			return diag.Errorf("Conversion Error: %s", err)
		}
		newDiskMB, err := convertToMegabytes(newValue.(string))
		if err != nil {
			d.Set("disk", oldValue.(string))
			return diag.Errorf("Conversion Error: %s", err)
		}
		if newDiskMB != oldDiskMB {
			// we could remove the old disk and recreate a new one, but then all the user data would be lost
			d.Set("disk", oldValue.(string))
			return diag.Errorf("Cannot change main disk size from %s to %s", oldValue.(string), newValue.(string))
		} else {
			// do not change disk size if value equal but a different unit
			tflog.Info(ctx, "Not replacing disk size " + oldValue.(string) + " with equal value " + newValue.(string))
		}
	}

	// Address image changes
	if d.HasChange("image") {
		oldValue, newValue := d.GetChange("image")
		oldImage := oldValue.(string)
		newImage := newValue.(string)
		// we could reapply a different image, but then all the user data would be lost
		d.Set("image", oldImage)
		return diag.Errorf("Cannot change image used to install the system from \"%s\" to \"%s\"", oldImage, newImage)
	}

	// TODO: address adapter VDev changes

	// Address desired MAC changes
	if d.HasChange("mac") {
		oldValue, newValue := d.GetChange("mac")
		oldMAC := oldValue.(string)
		newMAC := newValue.(string)
		if strings.ToLower(newMAC[8:]) != strings.ToLower(oldMAC[8:]) {
			// we could delete the network interface and create a new one, but then it would not be same interface anymore
			d.Set("mac", oldMAC)
			return diag.Errorf("Cannot change MAC address of main interface from \"%s\" to \"%s\"", oldMAC, newMAC)
		} else {
			tflog.Info(ctx, "Not replacing MAC address " + oldMAC + " with other MAC address with same last 3 hex bytes " + newMAC)
		}
	}

	// Address desired vswitch changes
	if d.HasChange("vswitch") {
		oldValue, newValue := d.GetChange("vswitch")
		oldVSwitch := oldValue.(string)
		newVSwitch := newValue.(string)
		// we could delete the network interface and create a new one, but then it would not be same interface anymore
		d.Set("vswitch", oldVSwitch)
		return diag.Errorf("Cannot change virtual switch of main interface from \"%s\" to \"%s\"", oldVSwitch, newVSwitch)
	}

	// Address cloud-init parameter changes
	if d.HasChange("cloudinit_params") {
		oldValue, newValue := d.GetChange("cloudinit_params")
		oldCloudinitParams := oldValue.(string)
		newCloudinitParams := newValue.(string)
		// we could redeploy, but then all the user data would be lost
		d.Set("cloudinit_params", oldCloudinitParams)
		return diag.Errorf("Cannot change cloud-init parameters used to install the system from \"%s\" to \"%s\"", oldCloudinitParams, newCloudinitParams)
	}

	return nil
}

func feilongGuestDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiClient).Client
	userid := d.Get("userid").(string)

	// Delete the guest
	err := client.DeleteGuest(userid)
	if err != nil {
		return diag.Errorf("Deletion Error: %s", err)
	}

	return nil
}

// For internal use

func convertToMegabytes(sizeWithUnit string) (int, error) {
	lastButOne := len(sizeWithUnit) - 1

	size, err := strconv.Atoi(sizeWithUnit[:lastButOne])
	if (err != nil) {
		return 0, err
	}

	unit := sizeWithUnit[lastButOne:]

	switch unit {
		case "B":
			size = size / 1_048_576
		case "K":
			size = size / 1_024
		case "M":
			// (nothing to do)
		case "G":
			size = size * 1_024
		case "T":
			size = size * 1_048_576
		default:
			return 0, errors.New("Unit must be one of B K M G T")
	}
	return size, nil
}

const waitingMsg string = "Still waiting for IP address"
const obtainedMsg string = "IP address obtained"

func waitForLease(ctx context.Context, client feilong.Client, userid string, macAddress *string, ipAddress *string) error {
	waitFunction := func() (interface{}, string, error) {
		err := getAddresses(client, userid, macAddress, ipAddress)
		if err != nil {
			return false, "", err
		}
		if *ipAddress == "" {
			return false, waitingMsg, nil
		}
		return true, obtainedMsg, nil
	}

	stateConf := &resource.StateChangeConf {
		Pending:	[]string { waitingMsg },
		Target:		[]string { obtainedMsg },
		Refresh:	waitFunction,
		Timeout:	1 * time.Minute,
		MinTimeout:	3 * time.Second,
		Delay:		5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func getAddresses(client feilong.Client, userid string, macAddress *string, ipAddress *string) error {
	result, err := client.GetGuestAdaptersInfo(userid)
	if err != nil {
		return err
	}

	if len(result.Output.Adapters) == 0 {
		return errors.New("Adapters Query Error")
	}
	*macAddress = result.Output.Adapters[0].MACAddress
	*ipAddress = result.Output.Adapters[0].IPAddress
	return nil
}
