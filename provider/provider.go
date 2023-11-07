package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Bischoff/feilong_api"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		return &schema.Provider {
			Schema: map[string]*schema.Schema{
				"connector": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("ZVM_CONNECTOR", nil),
					Description: "Domain name or address of the z/VM cloud connector",
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

	return &apiClient{*client}, nil
}
