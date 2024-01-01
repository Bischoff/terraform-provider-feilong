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
	err = client.UpdateGuestNIC(userid, "1000", &updateNICParams)
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
		tflog.Info(ctx, "Not replacing memory " +  d.Get("memory").(string) + " with equal value " + strconv.Itoa(obtainedMemory) + "M")
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
		tflog.Info(ctx, "Not replacing disk " + d.Get("disk").(string) + " with equal value " + strconv.Itoa(obtainedDiskSize) + "M")
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
	// client := meta.(*apiClient).Client

tflog.Info(ctx, "update function")
	// return diag.Errorf("not implemented")
	return nil
}

func feilongGuestDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiClient).Client

	userid := d.Get("userid").(string)

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
