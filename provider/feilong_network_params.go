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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func feilongNetworkParams() *schema.Resource {
	return &schema.Resource{
		Description:	"Feilong network parameters resource",

		CreateContext:	feilongNetworkParamsCreate,
		ReadContext:	feilongNetworkParamsRead,
		UpdateContext:	feilongNetworkParamsUpdate,
		DeleteContext:	feilongNetworkParamsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:	"Arbitrary name for the resource",
				Type:		schema.TypeString,
				Required:	true,
			},
			"os_distro": {
				Description:	"OS and distro of the network parameters",
				Type:		schema.TypeString,
				Required:	true,
			},
			"file": {
				Description:	"Path to the created resource",
				Type:		schema.TypeString,
				Computed:	true,
			},
		},
	}
}

func feilongNetworkParamsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Compute values written to the disk but not part of the data model
	resourceName := d.Get("name").(string)
	osDistro := d.Get("os_distro").(string)

	// Customize the parameters in the network tarball
	err := createNetworkTempDir(resourceName, osDistro)
	if err != nil {
		return diag.Errorf("Temporary Directory Creation Error: %s", err)
	}
	err = createNetworkTempFiles(resourceName, osDistro)
	if err != nil {
		return diag.Errorf("Temporary Files Creation Error: %s", err)
	}
	err = tarNetworkConfig(resourceName, osDistro)
	if err != nil {
		return diag.Errorf("Archive Creation Error: %s", err)
	}

	// Register the result
	err = d.Set("file", tmpdir + resourceName + "/network.doscript")
	if err != nil {
		return diag.Errorf("Network Parameters Tarball Registration Error: %s", err)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a network parameters tarball resource")

	// Set resource identifier
	d.SetId(resourceName)

	return nil
}

func feilongNetworkParamsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

	// return diag.Errorf("not implemented")
	return nil
}

func feilongNetworkParamsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

	// return diag.Errorf("not implemented")
	return nil
}

func feilongNetworkParamsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	resourceName := d.Get("name").(string)

	err := removeNetworkTempFiles(resourceName)
	if err != nil {
		return diag.Errorf("Temporary Files Removal Error: %s", err)
	}

	return nil
}

// For internal use

func createNetworkTempDir(resourceName string, osDistro string) error {
	path := tmpdir + resourceName + "/network.config/" + osDistro
	return os.MkdirAll(path, 0755)
}

func createNetworkTempFiles(resourceName string, osDistro string) error {
	err := create0000Tempfile(resourceName, osDistro)
	if err != nil {
		return err
	}

	err = create0001Tempfile(resourceName, osDistro)
	if err != nil {
		return err
	}

	err = create0002Tempfile(resourceName, osDistro)
	if err != nil {
		return err
	}

	err = create0003Tempfile(resourceName, osDistro)
	if err != nil {
		return err
	}

	return createInvokescriptTempfile(resourceName, osDistro)
}

func tarNetworkConfig(resourceName string, osDistro string) error {
	sourcePath := tmpdir + resourceName + "/network.config/" + osDistro + "/"
	destPath := tmpdir + resourceName + "/network.doscript"
	cmd := exec.Command(
		"/usr/bin/tar", "-C", sourcePath, "-cf", destPath,
		"0000", "0001", "0002", "0003", "invokeScript.sh")
	_, err := cmd.Output()
	return err
}

func removeNetworkTempFiles(resourceName string) error {
	cmd := exec.Command(
		"rm", "-r", tmpdir + resourceName)
	_, err := cmd.Output()
	return err
}

//go:embed files/network.config/sles15/0000
var file_0000 string

func create0000Tempfile(resourceName string, osDistro string) error {
	path := tmpdir + resourceName + "/network.config/" + osDistro + "/0000"
	return os.WriteFile(path, []byte(file_0000), 0644)
}

//go:embed files/network.config/sles15/0001
var file_0001 string

func create0001Tempfile(resourceName string, osDistro string) error {
	path := tmpdir + resourceName + "/network.config/" + osDistro + "/0001"
	return os.WriteFile(path, []byte(file_0001), 0644)
}

//go:embed files/network.config/sles15/0002
var file_0002 string

func create0002Tempfile(resourceName string, osDistro string) error {
	path := tmpdir + resourceName + "/network.config/" + osDistro + "/0002"
	return os.WriteFile(path, []byte(file_0002), 0644)
}

//go:embed files/network.config/sles15/0003
var file_0003 string

func create0003Tempfile(resourceName string, osDistro string) error {
	path := tmpdir + resourceName + "/network.config/" + osDistro + "/0003"
	return os.WriteFile(path, []byte(file_0003), 0644)
}

//go:embed files/network.config/sles15/invokeScript.sh
var file_invokescript string

func createInvokescriptTempfile(resourceName string, osDistro string) error {
	path := tmpdir + resourceName + "/network.config/" + osDistro + "/invokeScript.sh"
	return os.WriteFile(path, []byte(file_invokescript), 0644)
}
