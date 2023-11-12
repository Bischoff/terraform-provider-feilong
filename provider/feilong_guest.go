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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Description:	"System name for linux",
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
				Description:	"Memory with unit (G, M, k)",
				Type:		schema.TypeString,
				Optional:	true,
				Default:	"512M",
			},
			"disk": {
				Description:	"Disk size of first disk with unit (G, M, k)",
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
				Description:	"MAC address of first interface",
				Type:		schema.TypeString,
				Optional:	true,
				Default:	"",
			},
		},
	}
}

func feilongGuestCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Compute computed fields
	userid := d.Get("userid").(string)
	if userid == "" {
		name := d.Get("name").(string)
		userid = strings.ToUpper(name)
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
		return diag.Errorf("%s", err)
	}
	image := d.Get("image").(string)
	hostname := d.Get("name").(string)
	mac := d.Get("mac").(string)

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
		return diag.Errorf("%s", err)
	}

	// Create the first network interface
	createNICParams := feilong.CreateGuestNICParams {
		MACAddress:	mac,
	}
	err = client.CreateGuestNIC(userid, &createNICParams)
	if err != nil {
		return diag.Errorf("%s", err)
	}

	// Couple the first network interface to the virtual switch
	updateNICParams := feilong.UpdateGuestNICParams {
		Couple:		true,
		VSwitch:	"DEVNET",
	}
	err = client.UpdateGuestNIC(userid, "1000", &updateNICParams)
	if err != nil {
		return diag.Errorf("%s", err)
	}

	// Deploy the guest
	deployParams := feilong.DeployGuestParams {
		Image:          image,
		Hostname:       hostname,
	}
	err = client.DeployGuest(userid, &deployParams)
	if err != nil {
		return diag.Errorf("%s", err)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a Feilong guest resource")

	// Set resource identifier
	d.SetId(userid)

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
		return diag.Errorf("%s", err)
	}

	return nil
}

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
