/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Bischoff/feilong-client-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FeilongVSwitch{}
var _ resource.ResourceWithImportState = &FeilongVSwitch{}

func NewFeilongVSwitch() resource.Resource {
	return &FeilongVSwitch{}
}

// FeilongVSwitch defines the resource implementation.
type FeilongVSwitch struct {
	Client *feilong.Client
}

// FeilongVSwitchModel describes the resource data model.
type FeilongVSwitchModel struct {
	Name		types.String	`tfsdk:"name"`
	RealDevice	types.String	`tfsdk:"real_device"`
	Controller	types.String	`tfsdk:"controller"`
	ConnectionType	types.String	`tfsdk:"connection_type"`
	NetworkType	types.String	`tfsdk:"network_type"`
	Router		types.String	`tfsdk:"router"`
	VLANId		types.Int64	`tfsdk:"vlan_id"`
	PortType	types.String	`tfsdk:"port_type"`
	GVRP		types.String	`tfsdk:"gvrp"`
	QueueMem	types.Int64	`tfsdk:"queue_mem"`
	NativeVLANId	types.Int64	`tfsdk:"native_vlan_id"`
	Persist		types.Bool	`tfsdk:"persist"`
}

func (guest *FeilongVSwitch) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vswitch"
}

func (guest *FeilongVSwitch) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema {
		MarkdownDescription: "Feilong virtual switch resource",

		Attributes: map[string]schema.Attribute {
			"name": schema.StringAttribute {
				MarkdownDescription:	"Virtual switch name for z/VM",
				Required:		true,
			},
			"real_device": schema.StringAttribute {
				MarkdownDescription:	"Real device number",
				Optional:		true,
			},
			"controller": schema.StringAttribute {
				MarkdownDescription:	"Controller",
				Optional:		true,
			},
			"connection_type": schema.StringAttribute {
				MarkdownDescription:	"Connection type (CONNECT, DISCONNECT, or NOUPLINK)",
				Optional:		true,
			},
			"network_type": schema.StringAttribute {
				MarkdownDescription:	"Network type (IP or ETHERNET)",
				Optional:		true,
			},
			"router": schema.StringAttribute {
				MarkdownDescription:	"Router role (NONROUTER or PRIROUTER)",
				Optional:		true,
			},
			"vlan_id": schema.Int64Attribute {
				MarkdownDescription:	"VLAN identifier",
				Optional:		true,
			},
			"port_type": schema.StringAttribute {
				MarkdownDescription:	"Port type (ACCESS or TRUNK)",
				Optional:		true,
			},
			"gvrp": schema.StringAttribute {
				MarkdownDescription:	"Whether to use GVRP protocol (GVRP or NOGVRP)",
				Optional:		true,
			},
			"queue_mem": schema.Int64Attribute {
				MarkdownDescription:	"QDIO buffer size in megabytes",
				Optional:		true,
			},
			"native_vlan_id": schema.Int64Attribute {
				MarkdownDescription:	"Native VLAN identifier",
				Optional:		true,
			},
			"persist": schema.BoolAttribute {
				MarkdownDescription:	"Whether virtual switch is permanent",
				Optional:		true,
			},
		},
	}
}

func (guest *FeilongVSwitch) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	guest.Client = &req.ProviderData.(*apiClient).Client
}

func (guest *FeilongVSwitch) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeilongVSwitchModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the virtual switch
	client := guest.Client
	createParams := feilong.CreateVSwitchParams { Name: data.Name.ValueString() }
        if !data.RealDevice.IsNull() {
                createParams.RealDev = data.RealDevice.ValueString()
        }
        if !data.Controller.IsNull() {
                createParams.Controller = data.Controller.ValueString()
        }
        if !data.ConnectionType.IsNull() {
                createParams.Connection = data.ConnectionType.ValueString()
        }
        if !data.NetworkType.IsNull() {
                createParams.NetworkType = data.NetworkType.ValueString()
        }
        if !data.Router.IsNull() {
                createParams.Router = data.Router.ValueString()
        }
        if !data.VLANId.IsNull() {
                createParams.VLANId = data.VLANId.ValueInt64()
        }
        if !data.PortType.IsNull() {
                createParams.PortType = data.PortType.ValueString()
        }
        if !data.GVRP.IsNull() {
                createParams.GVRP = data.GVRP.ValueString()
        }
        if !data.QueueMem.IsNull() {
                createParams.QueueMem = int(data.QueueMem.ValueInt64())
        }
        if !data.NativeVLANId.IsNull() {
                createParams.NativeVLANId = int(data.NativeVLANId.ValueInt64())
        }
        if !data.Persist.IsNull() {
                createParams.Persist = data.Persist.ValueBool()
        }

	err := client.CreateVSwitch(&createParams)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a Feilong virtual switch resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (guest *FeilongVSwitch) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeilongVSwitchModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//	return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (guest *FeilongVSwitch) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FeilongVSwitchModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//	return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (guest *FeilongVSwitch) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FeilongVSwitchModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := guest.Client
	name := data.Name.ValueString()

	err := client.DeleteVSwitch(name)
	if err != nil {
		resp.Diagnostics.AddError("Deletion Error", fmt.Sprintf("Got error: %s", err))
		return
	}
}

func (guest *FeilongVSwitch) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
