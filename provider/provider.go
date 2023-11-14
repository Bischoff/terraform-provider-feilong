/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"context"
	"os"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/Bischoff/feilong-client-go"
)

// Ensure FeilongProvider satisfies various provider interfaces.
var _ provider.Provider = &FeilongProvider{}

// FeilongProvider defines the provider implementation.
type FeilongProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// FeilongProviderModel describes the provider data model.
type FeilongProviderModel struct {
	Connector types.String `tfsdk:"connector"`
}

func (p *FeilongProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "feilong"
	resp.Version = p.version
}

func (p *FeilongProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"connector": schema.StringAttribute{
				MarkdownDescription:	"Domain name or address of the z/VM cloud connector",
				Optional:		true,
			},
		},
	}
}

func (p *FeilongProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config FeilongProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Connector.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("connector"),
			"Unknown z/VM cloud connector",
			"The provider cannot create the Feilong client as there is an unknown configuration value for the z/VM cloud connector. " +
			"Please provide the value in the configuration, or use the ZVM_CONNECTOR environment variable.",
		)
		return
	}

	connector := os.Getenv("ZVM_CONNECTOR")
	if !config.Connector.IsNull() {
		connector = config.Connector.ValueString()
	}

	if connector == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("connector"),
			"Missing z/VM cloud connector",
			"The provider cannot create the Feilong client as there is a missing or empty value for the z/VM cloud connector. " +
			"Please make sure the value in the configuration, or of the ZVM_CONNECTOR environment variable, is not empty.",
		)
		return
	}

	// Create a new Feilong client using the configuration values
	client := feilong.NewClient(&connector, nil)

	// Check that the z/VM connector answers and that the API is of expected version
	result, err := client.GetFeilongVersion()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to contact z/VM connector, got error: %s", err))
		return
	}
	if result.Output.APIVersion != "1.0" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Expected Feilong API version 1.0, got: %s", result.Output.APIVersion))
		return
	}

	// Make the Feilong client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *FeilongProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	// return []func() datasource.DataSource{
	// 		NewFeilongDataSource,
	// }
	return nil
}

func (p *FeilongProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
//		NewFeilongCloudinitDisk,
		NewFeilongGuest,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FeilongProvider{
			version: version,
		}
	}
}
