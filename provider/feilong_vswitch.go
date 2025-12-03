/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"context"
	"fmt"
	"strings"
	"strconv"
	"golang.org/x/exp/maps"

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
	VSwitch		types.String	`tfsdk:"vswitch"`
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
				MarkdownDescription:	"Arbitrary name for the resource",
				Required:		true,
			},
			"vswitch": schema.StringAttribute {
				MarkdownDescription:	"Virtual switch name for z/VM",
				Optional:		true,
				Computed:		true,
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

	// Compute computed fields
	resourceName := data.Name.ValueString()
	vswitch := data.VSwitch.ValueString()
	if vswitch == "" {
		vswitch = strings.ToUpper(resourceName)
		if (len(vswitch) > 8) {
			vswitch = vswitch[:8]
		}
		data.VSwitch = types.StringValue(vswitch)
	}

	// Create the virtual switch
	client := guest.Client
	persist := data.Persist.ValueBool()
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
                createParams.Persist = &persist
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

	client := guest.Client
	vswitch := data.VSwitch.ValueString()

	// Obtain info about this vswitch
	vswitchDetails, err := client.GetVSwitchDetails(vswitch)
	if err != nil {
		resp.Diagnostics.AddError("VSwitch Querying Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Read real device
	devices := maps.Keys(vswitchDetails.Output.RealDevices)
	if len(devices) != 1 {
		resp.Diagnostics.AddError("VSwitch Querying Error", fmt.Sprintf("Unexpected number of real devices %d for vswitch %s", len(devices), vswitch))
		return
	}
	realDevice := devices[0]
	data.RealDevice = types.StringValue(realDevice)

	// Read controller
	controller := vswitchDetails.Output.RealDevices[realDevice].Controller
	if data.Controller.IsNull() && controller == "NONE" {
		tflog.Info(ctx, "Not replacing undeclared controller with default value NONE")
	} else {
		data.Controller = types.StringValue(controller)
	}

	// Read network type
	networkType := vswitchDetails.Output.TransportType
	if data.NetworkType.IsNull() && networkType == "ETHERNET" {
		tflog.Info(ctx, "Not replacing undeclared network type with default value ETHERNET")
	} else {
		data.NetworkType = types.StringValue(networkType)
	}

	// Read VLAN id
	vlanId, err := strconv.Atoi(vswitchDetails.Output.VLANId)
	if err != nil {
		resp.Diagnostics.AddError("VLAN Id Conversion Error: %s", fmt.Sprintf("Got error: %s", err))
		return
	}
	data.VLANId = types.Int64Value(int64(vlanId))

	// Read port type
	portType := vswitchDetails.Output.PortType
	data.PortType = types.StringValue(portType)

	// Read GVRP
	gvrp := vswitchDetails.Output.GVRPEnabledAttribute
	if data.GVRP.IsNull() && gvrp == "NOGVRP" {
		tflog.Info(ctx, "Not replacing undeclared GVRP with default value NOGVRP")
	} else {
		data.GVRP = types.StringValue(gvrp)
	}

	// Read queue memory
	queueMem, err := strconv.Atoi(vswitchDetails.Output.QueueMemoryLimit)
	if err != nil {
		resp.Diagnostics.AddError("Queue Memory Conversion Error: %s", fmt.Sprintf("Got error: %s", err))
		return
	}
	if data.QueueMem.IsNull() && queueMem == 8 {
		tflog.Info(ctx, "Not replacing undeclared queue memory with default value 8")
	} else {
		data.QueueMem = types.Int64Value(int64(queueMem))
	}

	// Read native VLAN id
	nativeVlanId, err := strconv.Atoi(vswitchDetails.Output.NativeVLANId)
	if err != nil {
		resp.Diagnostics.AddError("Native VLAN Id Conversion Error: %s", fmt.Sprintf("Got error: %s", err))
		return
	}
	if data.NativeVLANId.IsNull() && nativeVlanId == 1 {
		tflog.Info(ctx, "Not replacing undeclared native VLAN id with default value 1")
	} else {
		data.NativeVLANId = types.Int64Value(int64(nativeVlanId))
	}

	// CAVEATS:
	//  - the connection type cannot be determined after the deployment
	//  - the router cannot be determined after the deployment
	//  - the persist flag cannot be determined after the deployment

	// Write logs using the tflog package
	tflog.Trace(ctx, "Read characteristics of Feilong virtual switch resource")

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
