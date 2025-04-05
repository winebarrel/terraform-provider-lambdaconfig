package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.ResourceWithConfigure   = &ConcurrencyResource{}
	_ resource.ResourceWithImportState = &ConcurrencyResource{}
)

func NewConcurrencyResource() resource.Resource {
	return &ConcurrencyResource{}
}

type ConcurrencyResource struct {
	client *lambda.Client
}

type ConcurrencyResourceModel struct {
	FunctionName                 types.String `tfsdk:"function_name"`
	ReservedConcurrentExecutions types.Int32  `tfsdk:"reserved_concurrent_executions"`
}

func (r *ConcurrencyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_concurrency"
}

func (r *ConcurrencyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"function_name": schema.StringAttribute{
				MarkdownDescription: "Name of the lambda function.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"reserved_concurrent_executions": schema.Int32Attribute{
				MarkdownDescription: "Amount of reserved concurrent executions for this lambda function.",
				Required:            true,
			},
		},
	}
}

func (r *ConcurrencyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*lambda.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *lambda.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}

	r.client = client
}

func (r *ConcurrencyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConcurrencyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	output, err := r.client.PutFunctionConcurrency(ctx, &lambda.PutFunctionConcurrencyInput{
		FunctionName:                 plan.FunctionName.ValueStringPointer(),
		ReservedConcurrentExecutions: plan.ReservedConcurrentExecutions.ValueInt32Pointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Error Putting Lambda function concurrency", err.Error())
		return
	}

	plan.ReservedConcurrentExecutions = types.Int32Value(aws.ToInt32(output.ReservedConcurrentExecutions))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ConcurrencyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConcurrencyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	output, err := r.client.GetFunctionConcurrency(ctx, &lambda.GetFunctionConcurrencyInput{
		FunctionName: state.FunctionName.ValueStringPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Error Reading Lambda function concurrency", err.Error())
		return
	}

	state.ReservedConcurrentExecutions = types.Int32Value(aws.ToInt32(output.ReservedConcurrentExecutions))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ConcurrencyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConcurrencyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	output, err := r.client.PutFunctionConcurrency(ctx, &lambda.PutFunctionConcurrencyInput{
		FunctionName:                 plan.FunctionName.ValueStringPointer(),
		ReservedConcurrentExecutions: plan.ReservedConcurrentExecutions.ValueInt32Pointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Error Putting Lambda function concurrency", err.Error())
		return
	}

	plan.ReservedConcurrentExecutions = types.Int32Value(aws.ToInt32(output.ReservedConcurrentExecutions))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ConcurrencyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (r *ConcurrencyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("function_name"), req, resp)
}
