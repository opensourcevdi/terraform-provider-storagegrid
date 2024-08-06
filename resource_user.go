package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceUser struct {
	Provider *Provider
	State    struct {
		Name    types.String `tfsdk:"name"`
		Id      types.String `tfsdk:"id"`
		UserUrn types.String `tfsdk:"user_urn"`
	}
}

func (r *ResourceUser) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *ResourceUser) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"user_urn": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ResourceUser) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	req.Plan.Get(ctx, &r.State)
	user, err := r.Provider.client.CreateUser(r.State.Name.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic(err.Error(), ""))
		return
	}
	r.State.Id = types.StringValue(user.Id)
	r.State.UserUrn = types.StringValue(user.UserURN)
	resp.State.Set(ctx, &r.State)
}

func (r *ResourceUser) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	req.State.Get(ctx, &r.State)
	user, err := r.Provider.client.GetUser(r.State.Id.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diag.NewWarningDiagnostic(err.Error(), ""))
		resp.State.RemoveResource(ctx)
		return
	}
	r.State.Id = types.StringValue(user.Id)
	r.State.UserUrn = types.StringValue(user.UserURN)
	resp.State.Set(ctx, &r.State)
}

func (r *ResourceUser) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	req.State.Get(ctx, &r.State)
	err := r.Provider.client.DeleteUser(r.State.Id.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic(err.Error(), ""))
		return
	}
}

var _ resource.ResourceWithConfigure = (*ResourceUser)(nil)
var _ resource.ResourceWithModifyPlan = (*ResourceUser)(nil)

func (r *ResourceUser) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.Provider = req.ProviderData.(*Provider)
	}
}
func (r *ResourceUser) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	resp.RequiresReplace = resp.RequiresReplace.Append(path.Empty())
}

func (r *ResourceUser) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Will never be called due to `RequiresReplace` in `ModifyPlan()`
}
