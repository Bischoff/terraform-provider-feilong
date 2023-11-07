package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
				Description:	"System name for Linux",
				Type:		schema.TypeString,
				Required:	true,
			},
			"userid": {
				Description:	"System name for system/Z",
				Type:		schema.TypeString,
				Optional:	true,
				Computed:	true,
// compute userid dynamically here, if not provided?
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
				Description:	"Disk size with unit (G, M, k)",
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
		},
	}
}

func feilongGuestCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

// Do the real creation here

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a Feilong guest resource")

	d.SetId("foobar")
	// ideally, we should get the ID from feilong API

	return nil
}

func feilongGuestRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

	return diag.Errorf("not implemented")
}

func feilongGuestUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

	return diag.Errorf("not implemented")
}

func feilongGuestDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

	return diag.Errorf("not implemented")
}
