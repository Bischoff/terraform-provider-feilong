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
			"network_params": {
				Description:	"Path to network parameters file",
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
	networkParams := d.Get("network_params").(string)
	cloudinitParams := d.Get("cloudinit_params").(string)
	localUser := meta.(*apiClient).LocalUser
	transportFiles := ""
	remoteHost := ""
	if networkParams != "" {
		if cloudinitParams != "" {
			transportFiles = networkParams + "," + cloudinitParams
			remoteHost = localUser
		} else {
			transportFiles = networkParams
			remoteHost = localUser
		}
	} else {
		if cloudinitParams != "" {
			transportFiles = cloudinitParams
			remoteHost = localUser
		}
	}

	// Create the guest
	client := meta.(*apiClient).Client
	diskList := []feilong.CreateGuestDisk {
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
		TransportFiles:	transportFiles,
		RemoteHost:	remoteHost,
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
	// client := meta.(*apiClient).Client

	// return diag.Errorf("not implemented")
	return nil
}

func feilongGuestUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

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
		Pending:    []string { waitingMsg },
		Target:     []string { obtainedMsg },
		Refresh:    waitFunction,
		Timeout:    1 * time.Minute,
		MinTimeout: 3 * time.Second,
		Delay:      5 * time.Second,
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
