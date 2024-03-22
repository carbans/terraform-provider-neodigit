// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure NeodigitProvider satisfies various provider interfaces.
var _ provider.Provider = &NeodigitProvider{}

// ScaffoldingProvider defines the provider implementation.
type NeodigitProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version  string
	Endpoint types.String `tfsdk:"endpoint"`
	Api_key  types.String `tfsdk:"api_key"`
}

// ScaffoldingProviderModel describes the provider data model.
type NeodigitProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Api_key  types.String `tfsdk:"api_key"`
}

func (p *NeodigitProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "scaffolding"
	resp.Version = p.version
}

func (p *NeodigitProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Provider endpoint",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *NeodigitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	// Retrieve the configuration values.
	var config NeodigitProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }
	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing endpoint required for provider configuration.",
			"The provider cannot create the client without the endpoint.",
		)
	}

	if config.Api_key.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing api_key required for provider configuration.",
			"The provider cannot create the client without the api_key.",
		)
	}

	endpoint := os.Getenv("NEODIGIT_ENDPOINT")
	api_key := os.Getenv("NEODIGIT_API_KEY")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.Api_key.IsNull() {
		api_key = config.Api_key.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Set Terraform environment variable NEODIGIT_ENDPOINT",
			"If either is already set, ensure the value is not empty.",
		)
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Set Terraform environment variable NEODIGIT_API_KEY",
			"If either is already set, ensure the value is not empty.",
		)
	}

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *NeodigitProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *NeodigitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NeodigitProvider{
			version: version,
		}
	}
}
