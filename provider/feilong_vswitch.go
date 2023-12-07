/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Bischoff/feilong-client-go"
)

func feilongVSwitch() *schema.Resource {
	return &schema.Resource{
		Description:	"Feilong virtual switch resource",

		CreateContext:	feilongVSwitchCreate,
		ReadContext:	feilongVSwitchRead,
		UpdateContext:	feilongVSwitchUpdate,
		DeleteContext:	feilongVSwitchDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:	"Arbitrary name for the resource",
				Type:		schema.TypeString,
				Required:	true,
			},
			"vswitch": {
				Description:	"Virtual switch name for z/VM",
				Type:		schema.TypeString,
				Optional:	true,
				Computed:	true,
			},
			"real_device": {
				Description:	"Real device number",
				Type:		schema.TypeString,
				Optional:	true,
			},
			"controller": {
				Description:	"Controller",
				Type:		schema.TypeString,
				Optional:	true,
			},
			"connection_type": {
				Description:	"Connection type (CONNECT, DISCONNECT, or NOUPLINK)",
				Type:		schema.TypeString,
				Optional:	true,
			},
			"network_type": {
				Description:	"Network type (IP or ETHERNET)",
				Type:		schema.TypeString,
				Optional:	true,
			},
			"router": {
				Description:	"Router role (NONROUTER or PRIROUTER)",
				Type:		schema.TypeString,
				Optional:	true,
			},
			"vlan_id": {
				Description:	"VLAN identifier",
				Type:		schema.TypeInt,
				Optional:	true,
			},
			"port_type": {
				Description:	"Port type (ACCESS or TRUNK)",
				Type:		schema.TypeString,
				Optional:	true,
			},
			"gvrp": {
				Description:	"Whether to use GVRP protocol (GVRP or NOGVRP)",
				Type:		schema.TypeString,
				Optional:	true,
			},
			"queue_mem": {
				Description:	"QDIO buffer size in megabytes",
				Type:		schema.TypeInt,
				Optional:	true,
			},
			"native_vlan_id": {
				Description:	"Native VLAN identifier",
				Type:		schema.TypeInt,
				Optional:	true,
			},
			"persist": {
				Description:	"Whether virtual switch is permanent",
				Type:		schema.TypeBool,
				Optional:	true,
			},
		},
	}
}

func feilongVSwitchCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Compute computed fields
	resourceName := d.Get("name").(string)
	vswitch := d.Get("vswitch").(string)
	if vswitch == "" {
		vswitch = strings.ToUpper(resourceName)
		if (len(vswitch) > 8) {
			vswitch = vswitch[:8]
		}
		d.Set("vswitch", vswitch)
	}

	// Create the virtual switch
	client := meta.(*apiClient).Client
	createParams := feilong.CreateVSwitchParams { Name: vswitch }
	if d.Get("real_device") != nil {
		createParams.RealDev = d.Get("real_device").(string)
	}
	if d.Get("controller") != nil {
		createParams.Controller = d.Get("controller").(string)
	}
	if d.Get("connection_type") != nil {
		createParams.Connection = d.Get("connection_type").(string)
	}
	if d.Get("network_type") != nil {
		createParams.NetworkType = d.Get("network_type").(string)
	}
	if d.Get("router") != nil {
		createParams.Router = d.Get("router").(string)
	}
	if d.Get("vlan_id") != nil {
		createParams.VLANId = d.Get("vlan_id").(int)
	}
	if d.Get("port_type") != nil {
		createParams.PortType = d.Get("port_type").(string)
	}
	if d.Get("gvrp") != nil {
		createParams.GVRP = d.Get("gvrp").(string)
	}
	if d.Get("queue_mem") != nil {
		createParams.QueueMem = d.Get("queue_mem").(int)
	}
	if d.Get("native_vlan_id") != nil {
		createParams.NativeVLANId = d.Get("native_vlan_id").(int)
	}
	if d.Get("persist") != nil {
		createParams.Persist = d.Get("persist").(bool)
	}

	err := client.CreateVSwitch(&createParams)
	if err != nil {
		return diag.Errorf("Creation Error: %s", err)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a Feilong vswitch resource")

	// Set resource identifier
	d.SetId(resourceName)

	return nil
}

func feilongVSwitchRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

	// return diag.Errorf("not implemented")
	return nil
}

func feilongVSwitchUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

	// return diag.Errorf("not implemented")
	return nil
}

func feilongVSwitchDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiClient).Client

	name := d.Get("name").(string)

	err := client.DeleteVSwitch(name)
	if err != nil {
		return diag.Errorf("Deletion Error: %s", err)
	}

	return nil
}
