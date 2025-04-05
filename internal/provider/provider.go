package provider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &LambdaconfigProvider{}

type LambdaconfigProvider struct {
	version string
}

type LambdaconfigProviderModel struct {
	Region types.String `tfsdk:"region"`
}

func (p *LambdaconfigProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lambdaconfig"
	resp.Version = p.version
}

func (p *LambdaconfigProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				MarkdownDescription: "The region to use.",
				Optional:            true,
			},
		},
	}
}

func (p *LambdaconfigProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data LambdaconfigProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	options := []func(*config.LoadOptions) error{}

	if !data.Region.IsNull() {
		options = append(options, config.WithRegion(data.Region.String()))
	}

	config, err := config.LoadDefaultConfig(ctx, options...)

	if err != nil {
		resp.Diagnostics.AddError("Unable to Load AWS config", err.Error())
		return
	}

	client := lambda.NewFromConfig(config)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *LambdaconfigProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewConcurrencyResource,
	}
}

func (p *LambdaconfigProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// No Data Sources
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LambdaconfigProvider{
			version: version,
		}
	}
}
