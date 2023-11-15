/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"context"
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
var _ resource.Resource = &FeilongNetworkParams{}
var _ resource.ResourceWithImportState = &FeilongNetworkParams{}

func NewFeilongNetworkParams() resource.Resource {
	return &FeilongNetworkParams{}
}

// FeilongNetworkParams defines the resource implementation.
type FeilongNetworkParams struct {
	ResourceName	string
	OSDistro	string
}

// FeilongNetworkParamsModel describes the resource data model.
type FeilongNetworkParamsModel struct {
	Name		types.String	`tfsdk:"name"`
	OSDistro	types.String	`tfsdk:"os_distro"`
	File		types.String	`tfsdk:"file"`
}

func (params *FeilongNetworkParams) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_params"
}

func (params *FeilongNetworkParams) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema {
		MarkdownDescription: "Feilong network parameters resource",

		Attributes: map[string]schema.Attribute {
			"name": schema.StringAttribute {
				MarkdownDescription:	"Arbitrary name for the resource",
				Required:		true,
			},
			"os_distro": schema.StringAttribute {
				MarkdownDescription:	"OS and distro of the network parameters",
				Required:		true,
			},
			"file": schema.StringAttribute {
				MarkdownDescription:	"Path to the created resource",
				Computed:		true,
			},
		},
	}
}

func (params *FeilongNetworkParams) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// Nothing to do
}

func (params *FeilongNetworkParams) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeilongNetworkParamsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Compute implementation variables
	params.ResourceName = data.Name.ValueString()
	params.OSDistro = data.OSDistro.ValueString()

	// Customize the parameters in the network tarball
	err := params.createTempDir()
	if err != nil {
		resp.Diagnostics.AddError("Temporary Directory Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}
	err = params.createTempFiles()
	if err != nil {
		resp.Diagnostics.AddError("Temporary Files Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}
	err = params.tarNetworkConfig()
	if err != nil {
		resp.Diagnostics.AddError("Archive Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Register the result
	data.File = types.StringValue(tmpdir + params.ResourceName + "/network.doscript")

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a network parameters tarball resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (params *FeilongNetworkParams) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeilongNetworkParamsModel

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

func (params *FeilongNetworkParams) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FeilongNetworkParamsModel

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

func (params *FeilongNetworkParams) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FeilongNetworkParamsModel

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

func (params *FeilongNetworkParams) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// For internal use

func (params *FeilongNetworkParams) createTempDir() error {
	path := tmpdir + params.ResourceName + "/network.config/" + params.OSDistro
	return os.MkdirAll(path, 0755)
}

func (params *FeilongNetworkParams) createTempFiles() error {
	err := params.create0000Tempfile()
	if err != nil {
		return err
	}

	err = params.create0001Tempfile()
	if err != nil {
		return err
	}

	err = params.create0002Tempfile()
	if err != nil {
		return err
	}

	err = params.create0003Tempfile()
	if err != nil {
		return err
	}

	return params.createInvokescriptTempfile()
}

func (params *FeilongNetworkParams) tarNetworkConfig() error {
	sourcePath := tmpdir + params.ResourceName + "/network.config/" + params.OSDistro + "/"
	destPath := tmpdir + params.ResourceName + "/network.doscript"
	cmd := exec.Command(
		"/usr/bin/tar", "-C", sourcePath, "-cf", destPath,
		"0000", "0001", "0002", "0003", "invokeScript.sh")
	_, err := cmd.Output()
	return err
}

func (params *FeilongNetworkParams) removeTempFiles() error {
	cmd := exec.Command(
		"rm", "-r", tmpdir + params.ResourceName)
	_, err := cmd.Output()
	return err
}

//go:embed files/network.config/sles15/0000
var file_0000 string

func (params *FeilongNetworkParams) create0000Tempfile() error {
	path := tmpdir + params.ResourceName + "/network.config/" + params.OSDistro + "/0000"
	return os.WriteFile(path, []byte(file_0000), 0644)
}

//go:embed files/network.config/sles15/0001
var file_0001 string

func (params *FeilongNetworkParams) create0001Tempfile() error {
	path := tmpdir + params.ResourceName + "/network.config/" + params.OSDistro + "/0001"
	return os.WriteFile(path, []byte(file_0001), 0644)
}

//go:embed files/network.config/sles15/0002
var file_0002 string

func (params *FeilongNetworkParams) create0002Tempfile() error {
	path := tmpdir + params.ResourceName + "/network.config/" + params.OSDistro + "/0002"
	return os.WriteFile(path, []byte(file_0002), 0644)
}

//go:embed files/network.config/sles15/0003
var file_0003 string

func (params *FeilongNetworkParams) create0003Tempfile() error {
	path := tmpdir + params.ResourceName + "/network.config/" + params.OSDistro + "/0003"
	return os.WriteFile(path, []byte(file_0003), 0644)
}

//go:embed files/network.config/sles15/invokeScript.sh
var file_invokescript string

func (params *FeilongNetworkParams) createInvokescriptTempfile() error {
	path := tmpdir + params.ResourceName + "/network.config/" + params.OSDistro + "/invokeScript.sh"
	return os.WriteFile(path, []byte(file_invokescript), 0644)
}
