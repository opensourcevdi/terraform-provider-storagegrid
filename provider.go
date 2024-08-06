package main

import (
	"context"
	"terraform-provider-storagegrid/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Provider struct {
	ApiUrl    types.String `tfsdk:"api_url"`
	AccountId types.String `tfsdk:"account_id"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
	client    client.Client
}

func (p *Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "storagegrid"
}

func (p *Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Required: true,
			},
			"account_id": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	req.Config.Get(ctx, p)
	p.client = client.Client{
		ApiUrl:    p.ApiUrl.ValueString(),
		AccountId: p.AccountId.ValueString(),
		Username:  p.Username.ValueString(),
		Password:  p.Password.ValueString(),
	}
	resp.DataSourceData = p
	resp.ResourceData = p
}

func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &ResourceUser{} },
		func() resource.Resource { return &ResourceAccesskey{} },
	}
}

var _ provider.Provider = (*Provider)(nil)
