package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Bischoff/feilong_api"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		return &schema.Provider {
			Schema: map[string]*schema.Schema{
				"connector": {
					Type:		schema.TypeString,
					Required:	true,
					DefaultFunc:	schema.EnvDefaultFunc("ZVM_CONNECTOR", nil),
					Description:	"Domain name or address of the z/VM cloud connector",
				},
			},

			// DataSourcesMap: map[string]*schema.Resource {
			//	"feilong_data_source": dataSourceFeilong(),
			//},

			ResourcesMap: map[string]*schema.Resource {
				"feilong_guest": feilongGuest(),
			},

			ConfigureFunc: providerConfigure,
		}
	}
}

type apiClient struct {
	Client feilong_api.Client
}

func providerConfigure(d *schema.ResourceData) (any, error) {
	connector := d.Get("connector").(string)

	client, err := feilong_api.NewClient(&connector)
	if err != nil {
		return nil, err
	}

	// Check that the z/VM connector answers and is of expected version
	httpResp, err := client.GetZvmCloudConnectorVersion()
	if err != nil {
		return nil, fmt.Errorf("Unable to contact z/VM connector, got error: %s", err)
	}
	if httpResp.Output.Version != "1.6.6" {
		return nil, fmt.Errorf("Expected z/VM connector version 1.6.6, got: %s", httpResp.Output.Version)
	}

	return &apiClient{*client}, nil
}
