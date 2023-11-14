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
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Bischoff/feilong-client-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FeilongGuest{}
var _ resource.ResourceWithImportState = &FeilongGuest{}

func NewFeilongGuest() resource.Resource {
	return &FeilongGuest{}
}

// FeilongGuest defines the resource implementation.
type FeilongGuest struct {
	client *feilong.Client
}

// FeilongGuestModel describes the resource data model.
type FeilongGuestModel struct {
	Name	types.String	`tfsdk:"name"`
	UserId	types.String	`tfsdk:"userid"`
	VCPUs	types.Int64	`tfsdk:"vcpus"`
	Memory	types.String	`tfsdk:"memory"`
	Disk	types.String	`tfsdk:"disk"`
	Image	types.String	`tfsdk:"image"`
	Mac	types.String	`tfsdk:"mac"`
}

func (r *FeilongGuest) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guest"
}

func (r *FeilongGuest) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema {
		MarkdownDescription: "Feilong guest VM resource",

		Attributes: map[string]schema.Attribute {
			"name": schema.StringAttribute {
				MarkdownDescription:	"System name for linux",
				Required:		true,
			},
			"userid": schema.StringAttribute {
				MarkdownDescription:	"System name for z/VM",
				Optional:		true,
				Computed:		true,
			},
			"vcpus": schema.Int64Attribute {
				MarkdownDescription:	"Virtual CPUs count",
				Optional:		true,
				Computed:		true,
				Default:		int64default.StaticInt64(1),
			},
			"memory": schema.StringAttribute {
				MarkdownDescription:	"Memory with unit (G, M, k)",
				Optional:		true,
				Computed:		true,
				Default:		stringdefault.StaticString("512M"),
			},
			"disk": schema.StringAttribute {
				MarkdownDescription:	"Disk size of first disk with unit (G, M, k)",
				Optional:		true,
				Computed:		true,
				Default:		stringdefault.StaticString("10G"),
				// we could use the size of the image file instead
			},
			"image": schema.StringAttribute {
				MarkdownDescription:	"Image name",
				Required:		true,
			},
			"mac": schema.StringAttribute {
				MarkdownDescription:	"MAC address of first interface",
				Optional:		true,
				Computed:		true,
				Default:		stringdefault.StaticString(""),
			},
		},
	}
}

func (guest *FeilongGuest) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client := req.ProviderData.(*feilong.Client)

	guest.client = client
}

func (guest *FeilongGuest) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeilongGuestModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Compute computed fields
	userid := data.UserId.ValueString()
	if userid == "" {
		name := data.Name.ValueString()
		userid = strings.ToUpper(name)
		if (len(userid) > 8) {
			userid = userid[:8]
		}
		data.UserId = types.StringValue(userid)
	}

	// Compute values passed to Feilong API but not part of the data model
	size := data.Disk.ValueString()
	vcpus := int(data.VCPUs.ValueInt64())
	memory, err := convertToMegabytes(data.Memory.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Conversion Error", fmt.Sprintf("Got error: %s", err))
		return
	}
	image := data.Image.ValueString()
	hostname := data.Name.ValueString()
	mac := data.Mac.ValueString()

	// Create the guest
	client := guest.client
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
		resp.Diagnostics.AddError("Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Create the first network interface
	createNICParams := feilong.CreateGuestNICParams {
		MACAddress:	mac,
	}
	err = client.CreateGuestNIC(userid, &createNICParams)
	if err != nil {
		resp.Diagnostics.AddError("NIC Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Couple the first network interface to the virtual switch
	updateNICParams := feilong.UpdateGuestNICParams {
		Couple:         true,
		VSwitch:        "DEVNET",
	}
	err = client.UpdateGuestNIC(userid, "1000", &updateNICParams)
	if err != nil {
		resp.Diagnostics.AddError("NIC Coupling Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Deploy the guest
	deployParams := feilong.DeployGuestParams {
		Image:		image,
		Hostname:	hostname,
	}
	err = client.DeployGuest(userid, &deployParams)
	if err != nil {
		resp.Diagnostics.AddError("Deployment Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Start the guest
	err = client.StartGuest(userid)
	if err != nil {
		resp.Diagnostics.AddError("Startup Error", fmt.Sprintf("Got error: %s", err))
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a Feilong guest resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (guest *FeilongGuest) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

func (guest *FeilongGuest) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func (guest *FeilongGuest) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FeilongGuestModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := guest.client
	userid := data.UserId.ValueString()

	err := client.DeleteGuest(userid)
	if err != nil {
		resp.Diagnostics.AddError("Deletion Error", fmt.Sprintf("Got error: %s", err))
		return
	}
}

func (guest *FeilongGuest) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
