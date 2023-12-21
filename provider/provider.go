/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Bischoff/feilong-client-go"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		return &schema.Provider {
			Schema: map[string]*schema.Schema{
				"connector": {
					Type:		schema.TypeString,
					Required:	true,
					Description:	"URL of the z/VM connector",
				},
				"admin_token": {
					Type:		schema.TypeString,
					Optional:	true,
					Description:	"Shared secret to authenticate the client",
				},
				"local_user": {
					Type:		schema.TypeString,
					Optional:	true,
					Description:	"Where parameter files are uploaded from",

				},
			},

			// DataSourcesMap: map[string]*schema.Resource {
			//	"feilong_data_source": dataSourceFeilong(),
			//},

			ResourcesMap: map[string]*schema.Resource {
				"feilong_cloudinit_params": feilongCloudinitParams(),
				"feilong_guest": feilongGuest(),
				"feilong_vswitch": feilongVSwitch(),
			},

			ConfigureFunc: providerConfigure,
		}
	}
}

type apiClient struct {
	Client		feilong.Client
	LocalUser	string
}

func providerConfigure(d *schema.ResourceData) (any, error) {
	connector := d.Get("connector").(string)
	client := feilong.NewClient(&connector, nil)
	adminToken := d.Get("admin_token").(string)
	localUser := d.Get("local_user").(string)

	// If needed, create an authentication token
	if adminToken != "" {
		err := client.CreateToken(adminToken)
		if err != nil {
			return nil, fmt.Errorf("Unable to create authentication token, got error: %s", err)
		}
	}

	// Check that the API is of expected version
	result, err := client.GetFeilongVersion()
	if err != nil {
		return nil, fmt.Errorf("Unable to contact z/VM connector, got error: %s", err)
	}
	if result.Output.APIVersion != "1.0" {
		return nil, fmt.Errorf("Expected Feilong API version 1.0, got: %s", result.Output.APIVersion)
	}

	return &apiClient{*client, localUser}, nil
}
