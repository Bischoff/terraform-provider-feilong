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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const tmpdir string = "/tmp/terraform-provider-feilong/"

func feilongCloudinitParams() *schema.Resource {
	return &schema.Resource{
		Description:	"Feilong cloud-init parameters resource",

		CreateContext:	feilongCloudinitParamsCreate,
		ReadContext:	feilongCloudinitParamsRead,
		UpdateContext:	feilongCloudinitParamsUpdate,
		DeleteContext:	feilongCloudinitParamsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:	"Arbitrary name for the resource",
				Type:		schema.TypeString,
				Required:	true,
			},
			"hostname": {
				Description:	"Fully-qualified domain name of the guest",
				Type:		schema.TypeString,
				Required:	true,
			},
			"public_key": {
				Description:	"SSH public key for the default user on the guest",
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

func feilongCloudinitParamsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Compute values written to the disk but not part of the data model
	resourceName := d.Get("name").(string)
	hostname := d.Get("hostname").(string)
	publicKey := d.Get("public_key").(string)

	// Customize the parameters on the cloud-init disk
	err := createCloudinitTempDir(resourceName)
	if err != nil {
		return diag.Errorf("Temporary Directory Creation Error: %s", err)
	}
	err = createMetadataTempfile(resourceName, hostname, publicKey)
	if err != nil {
		return diag.Errorf("Temporary Files Creation Error: %s", err)
	}
	err = makeIsoDrive(resourceName)
	if err != nil {
		return diag.Errorf("ISO 9660 Creation Error: %s", err)
	}

	// Register the result
	err = d.Set("file", tmpdir + resourceName + "/cfgdrive.iso")
	if err != nil {
		return diag.Errorf("Cloud-init Parameters Disk Registration Error: %s", err)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a cloud-init parameters disk resource")

	// Set resource identifier
	d.SetId(resourceName)

	return nil
}

func feilongCloudinitParamsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

	// return diag.Errorf("not implemented")
	return nil
}

func feilongCloudinitParamsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// client := meta.(*apiClient).Client

	// return diag.Errorf("not implemented")
	return nil
}

func feilongCloudinitParamsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	resourceName := d.Get("name").(string)

	err := removeCloudinitTempFiles(resourceName)
	if err != nil {
		return diag.Errorf("Temporary Files Removal Error: %s", err)
	}

	return nil
}

// For internal use

func createCloudinitTempDir(resourceName string) error {
	path := tmpdir + resourceName + "/cfgdrive/openstack/latest"
	return os.MkdirAll(path, 0755)
}

func makeIsoDrive(resourceName string) error {
	sourcePath := tmpdir + resourceName + "/cfgdrive/"
	destPath := tmpdir + resourceName + "/cfgdrive.iso"
	cmd := exec.Command(
		"/usr/bin/mkisofs", "-o", destPath,
		"-quiet", "-ldots", "-allow-lowercase", "-allow-multidot", "-l", "-J", "-r",
		"-V", "config-2", sourcePath)
	_, err := cmd.Output()
	return err
}

func removeCloudinitTempFiles(resourceName string) error {
	cmd := exec.Command(
		"rm", "-r", tmpdir + resourceName)
	_, err := cmd.Output()
	return err
}

//go:embed files/cfgdrive/openstack/latest/meta_data.json
var file_metadata string

func createMetadataTempfile(resourceName string, hostname string, publicKey string) error {
	path := tmpdir + resourceName + "/cfgdrive/openstack/latest/meta_data.json"
	shortname, _, _ := strings.Cut(hostname, ".")
	s := file_metadata
	s = strings.Replace(s, "HOSTNAME", hostname, -1)
	s = strings.Replace(s, "NAME", shortname, -1)
	s = strings.Replace(s, "PUBLIC_KEY", publicKey, -1)
	return os.WriteFile(path, []byte(s), 0644)
}
