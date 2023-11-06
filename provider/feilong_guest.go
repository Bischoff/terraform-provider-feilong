package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Bischoff/feilong_api"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FeilongGuest{}
var _ resource.ResourceWithImportState = &FeilongGuest{}

func NewFeilongGuest() resource.Resource {
	return &FeilongGuest{}
}

// FeilongGuest defines the resource implementation.
type FeilongGuest struct {
	client *feilong_api.Client
}

// FeilongGuestModel describes the resource data model.
type FeilongGuestModel struct {
	Name	types.String	`tfsdk:"name"`
	UserID	types.String	`tfsdk:"userid"`
	VCPUs	types.Int64	`tfsdk:"vcpus"`
	Memory	types.String	`tfsdk:"memory"`
	Disk	types.String	`tfsdk:"disk"`
	Image	types.String	`tfsdk:"image"`
}

func (r *FeilongGuest) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guest"
}

func (r *FeilongGuest) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema {
		MarkdownDescription: "Feilong guest VM resource",

		Attributes: map[string]schema.Attribute {
			"name": schema.StringAttribute {
				MarkdownDescription: "System name for Linux",
				Required:            true,
			},
			"userid": schema.StringAttribute {
				MarkdownDescription: "System name for system/Z",
				Optional:            true,
				Computed:            true,
			},
			"vcpus": schema.Int64Attribute {
				MarkdownDescription: "Virtual CPUs count",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"memory": schema.StringAttribute {
				MarkdownDescription: "Memory with unit (G, M, k)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1G"),
			},
			"disk": schema.StringAttribute {
				MarkdownDescription: "Disk size with unit (G, M, k)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("10G"),
			},
			"image": schema.StringAttribute {
				MarkdownDescription: "Virtual CPUs count",
				Required:            true,
			},
		},
	}
}

func (r *FeilongGuest) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client := req.ProviderData.(*feilong_api.Client)

	r.client = client
}

func (r *FeilongGuest) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeilongGuestModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check that the z/VM connector answers and is of expected version
	httpResp, err := r.client.GetZvmCloudConnectorVersion()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to contact z/VM connector, got error: %s", err))
		return
	}
	if httpResp.Output.Version != "1.6.6" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Expected z/VM connector version 1.6.6, got: %s", httpResp.Output.Version))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	// data.Name = types.StringValue("foobar")
// compute userid dynamically here, if not provided? or above?

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a Feilong guest resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FeilongGuest) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeilongGuestModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FeilongGuest) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FeilongGuestModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FeilongGuest) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FeilongGuestModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *FeilongGuest) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
