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
					DefaultFunc:	schema.EnvDefaultFunc("ZVM_CONNECTOR", nil),
					Description:	"Domain name or address of the z/VM connector",
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
				"feilong_network_params": feilongNetworkParams(),
				"feilong_guest": feilongGuest(),
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
	localUser := d.Get("local_user").(string)

	// Check that the z/VM connector answers and that the API is of expected version
	result, err := client.GetFeilongVersion()
	if err != nil {
		return nil, fmt.Errorf("Unable to contact z/VM connector, got error: %s", err)
	}
	if result.Output.APIVersion != "1.0" {
		return nil, fmt.Errorf("Expected Feilong API version 1.0, got: %s", result.Output.APIVersion)
	}

	return &apiClient{*client, localUser}, nil
}
