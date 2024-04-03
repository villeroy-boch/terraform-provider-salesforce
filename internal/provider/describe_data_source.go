package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/villeroy-boch/tf-provider-salesforce/internal/salesforce"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &descriptionDataSource{}
	_ datasource.DataSourceWithConfigure = &descriptionDataSource{}
)

// NewDescriptionDataSource is a helper function to simplify the provider implementation.
func NewDescriptionDataSource() datasource.DataSource {
	return &descriptionDataSource{}
}

// descriptionDataSource is the data source implementation.
type descriptionDataSource struct {
	client *salesforce.Client
}

// Configure adds the provider configured client to the data source.
func (d *descriptionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Salesforce Description data source")

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*salesforce.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *salesforce.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client

	tflog.Info(ctx, "Configured Salesforce Description data source", map[string]any{"success": true})
}

// Metadata returns the data source type name.
func (d *descriptionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_description"
}

// Schema defines the schema for the data source.
func (d *descriptionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Fetches a description.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the described object.",
				Required:    true,
			},
			"label": schema.StringAttribute{
				Description: "Label of the described object.",
				Computed:    true,
			},
			"fields": schema.ListNestedAttribute{
				Description: "Fields of the described object.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field.",
							Computed:    true,
						},
						"label": schema.StringAttribute{
							Description: "Label of the field.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the field.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// descriptionDataSourceModel maps the data source schema data.
type descriptionDataSourceModel struct {
	ID     types.String            `tfsdk:"id"`
	Name   types.String            `tfsdk:"name"`
	Label  types.String            `tfsdk:"label"`
	Fields []descriptionFieldModel `tfsdk:"fields"`
}

type descriptionFieldModel struct {
	Name  types.String `tfsdk:"name"`
	Label types.String `tfsdk:"label"`
	Type  types.String `tfsdk:"type"`
}

// Read refreshes the Terraform state with the latest data.
func (d *descriptionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state descriptionDataSourceModel

	// log req
	req.Config.Get(ctx, &state)

	tflog.Info(ctx, "Reading Salesforce Description data source", map[string]any{
		"input": fmt.Sprintf("%+v", state),
	})

	description, err := d.client.GetDescription(
		state.Name.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Salesforce descriptions",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, field := range description.Fields {
		col := descriptionFieldModel{
			Name:  types.StringValue(field.Name),
			Label: types.StringValue(field.Label),
			Type:  types.StringValue(field.Type),
		}

		state.Fields = append(state.Fields, col)
	}

	state.ID = types.StringValue("placeholder")
	state.Name = types.StringValue(description.Name)
	state.Label = types.StringValue(description.Label)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
