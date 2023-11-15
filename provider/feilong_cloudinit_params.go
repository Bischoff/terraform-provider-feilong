/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"context"
	"strings"
	_ "embed"
	"os"
	"os/exec"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FeilongCloudinitParams{}
var _ resource.ResourceWithImportState = &FeilongCloudinitParams{}

const tmpdir string = "/tmp/terraform-provider-feilong/"

func NewFeilongCloudinitParams() resource.Resource {
	return &FeilongCloudinitParams{}
}

// FeilongCloudinitParams defines the resource implementation.
type FeilongCloudinitParams struct {
	ResourceName	string
	Hostname	string
	PublicKey	string
}

// FeilongCloudinitParamsModel describes the resource data model.
type FeilongCloudinitParamsModel struct {
	Name		types.String	`tfsdk:"name"`
	Hostname	types.String	`tfsdk:"hostname"`
	PublicKey	types.String	`tfsdk:"public_key"`
	File		types.String	`tfsdk:"file"`
}

func (params *FeilongCloudinitParams) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudinit_params"
}

func (params *FeilongCloudinitParams) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema {
		MarkdownDescription: "Feilong cloud-init parameters resource",

		Attributes: map[string]schema.Attribute {
			"name": schema.StringAttribute {
				MarkdownDescription:	"Arbitrary name for the resource",
				Required:		true,
			},
			"hostname": schema.StringAttribute {
				MarkdownDescription:	"Fully-qualified domain name of the guest",
				Required:		true,
			},
			"public_key": schema.StringAttribute {
				MarkdownDescription:	"SSH public key for the default user on the guest",
				Required:		true,
			},
			"file": schema.StringAttribute {
				MarkdownDescription:	"Path to the created resource",
				Computed:		true,
			},
		},
	}
}

func (params *FeilongCloudinitParams) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// Nothing to do
}

func (params *FeilongCloudinitParams) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeilongCloudinitParamsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Compute implementation variables
	params.ResourceName = data.Name.ValueString()
	params.Hostname = data.Hostname.ValueString()
	params.PublicKey = data.PublicKey.ValueString()

	// Customize the parameters on the cloud-init disk
	err := params.createTempDir()
	if err != nil {
		resp.Diagnostics.AddError("Temporary Directory Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}
	err = params.createMetadataTempfile()
	if err != nil {
		resp.Diagnostics.AddError("Temporary Files Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}
	err = params.makeIsoDrive()
	if err != nil {
		resp.Diagnostics.AddError("ISO 9660 Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Register the result
	data.File = types.StringValue(tmpdir + params.ResourceName + "/cfgdrive.iso")

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a cloud-init parameters disk resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (params *FeilongCloudinitParams) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeilongCloudinitParamsModel

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

func (params *FeilongCloudinitParams) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FeilongCloudinitParamsModel

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

func (params *FeilongCloudinitParams) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FeilongCloudinitParamsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Compute implementation variables
	params.ResourceName = data.Name.ValueString()

	err := params.removeTempFiles()
	if err != nil {
		resp.Diagnostics.AddError("Temporary Files Removal Error", fmt.Sprintf("Got error: %s", err))
		return
	}
}

func (params *FeilongCloudinitParams) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// For internal use

func (params *FeilongCloudinitParams) createTempDir() error {
	path := tmpdir + params.ResourceName + "/cfgdrive/openstack/latest"
	return os.MkdirAll(path, 0755)
}

func (params *FeilongCloudinitParams) makeIsoDrive() error {
	sourcePath := tmpdir + params.ResourceName + "/cfgdrive/"
	destPath := tmpdir + params.ResourceName + "/cfgdrive.iso"
	cmd := exec.Command(
		"/usr/bin/mkisofs", "-o", destPath,
		"-quiet", "-ldots", "-allow-lowercase", "-allow-multidot", "-l", "-J", "-r",
		"-V", "config-2", sourcePath)
	_, err := cmd.Output()
	return err
}

func (params *FeilongCloudinitParams) removeTempFiles() error {
	cmd := exec.Command(
		"rm", "-r", tmpdir + params.ResourceName)
	_, err := cmd.Output()
	return err
}

//go:embed files/cfgdrive/openstack/latest/meta_data.json
var file_metadata string

func (params *FeilongCloudinitParams) createMetadataTempfile() error {
	path := tmpdir + params.ResourceName + "/cfgdrive/openstack/latest/meta_data.json"
	shortname, _, _ := strings.Cut(params.Hostname, ".")
	s := file_metadata
	s = strings.Replace(s, "HOSTNAME", params.Hostname, -1)
	s = strings.Replace(s, "NAME", shortname, -1)
	s = strings.Replace(s, "PUBLIC_KEY", params.PublicKey, -1)
	return os.WriteFile(path, []byte(s), 0644)
}
