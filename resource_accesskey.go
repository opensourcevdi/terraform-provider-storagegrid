package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceAccesskey struct {
	Provider *Provider
	State    struct {
		UserId          types.String `tfsdk:"user_id"`
		Id              types.String `tfsdk:"id"`
		AccessKey       types.String `tfsdk:"access_key_id"`
		SecretAccessKey types.String `tfsdk:"secret_access_key"`
		UserUrn         types.String `tfsdk:"user_urn"`
	}
}

func (r *ResourceAccesskey) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_key"
}

func (r *ResourceAccesskey) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"access_key_id": schema.StringAttribute{
				Computed: true,
			},
			"secret_access_key": schema.StringAttribute{
				Computed: true,
			},
			"user_urn": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ResourceAccesskey) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	req.Plan.Get(ctx, &r.State)
	key, err := r.Provider.client.CreateAccessKey(r.State.UserId.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic(err.Error(), ""))
		return
	}
	r.State.Id = types.StringValue(key.Id)
	r.State.AccessKey = types.StringValue(key.AccessKey)
	r.State.SecretAccessKey = types.StringValue(key.SecretAccessKey)
	r.State.UserUrn = types.StringValue(key.UserURN)
	resp.State.Set(ctx, &r.State)
}

func (r *ResourceAccesskey) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	req.State.Get(ctx, &r.State)
	_, err := r.Provider.client.GetAccessKey(r.State.UserId.ValueString(), r.State.AccessKey.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diag.NewWarningDiagnostic(err.Error(), ""))
		resp.State.RemoveResource(ctx)
		return
	}
	resp.State.Set(ctx, &r.State)
}

func (r *ResourceAccesskey) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	req.State.Get(ctx, &r.State)
	err := r.Provider.client.DeleteAccessKey(r.State.UserId.ValueString(), r.State.Id.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic(err.Error(), ""))
		return
	}
}

var _ resource.ResourceWithConfigure = (*ResourceAccesskey)(nil)
var _ resource.ResourceWithModifyPlan = (*ResourceAccesskey)(nil)

func (r *ResourceAccesskey) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.Provider = req.ProviderData.(*Provider)
	}
}
func (r *ResourceAccesskey) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	resp.RequiresReplace = resp.RequiresReplace.Append(path.Empty())
}

func (r *ResourceAccesskey) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Will never be called due to `RequiresReplace` in `ModifyPlan()`
}
