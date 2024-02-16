package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/mondata-dev/terraform-provider-salesforce/internal/salesforce"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &salesforceProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &salesforceProvider{
			version: version,
		}
	}
}

// salesforceProvider is the provider implementation.
type salesforceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// salesforceProviderModel maps provider schema data to a Go type.
type salesforceProviderModel struct {
	ApiHost      types.String `tfsdk:"api_host"`
	ApiVersion   types.String `tfsdk:"api_version"`
	AuthHost     types.String `tfsdk:"auth_host"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	GrantType    types.String `tfsdk:"grant_type"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *salesforceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "salesforce"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *salesforceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Salesforce",
		Attributes: map[string]schema.Attribute{
			"api_host": schema.StringAttribute{
				Optional:    true,
				Description: "URI for Salesforce API. May also be provided via SALESFORCE_API_HOST environment variable.",
			},
			"api_version": schema.StringAttribute{
				Optional:    true,
				Description: "Version for Salesforce API. May also be provided via SALESFORCE_API_VERSION environment variable.",
			},
			"auth_host": schema.StringAttribute{
				Optional:    true,
				Description: "URI for Salesforce API Authentication. May also be provided via SALESFORCE_AUTH_HOST environment variable.",
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Client ID for Salesforce API. May also be provided via SALESFORCE_CLIENT_ID environment variable.",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Client Secret for Salesforce API. May also be provided via SALESFORCE_CLIENT_SECRET environment variable.",
			},
			"grant_type": schema.StringAttribute{
				Optional:    true,
				Description: "Grant type for Salesforce API. May also be provided via SALESFORCE_GRANT_TYPE environment variable.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username for Salesforce API. May also be provided via SALESFORCE_USERNAME environment variable.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password for Salesforce API. May also be provided via SALESFORCE_PASSWORD environment variable.",
			},
		},
	}
}

// Configure prepares a Salesforce API client for data sources and resources.
func (p *salesforceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Salesforce client")

	// Retrieve provider data from configuration
	var config salesforceProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ApiHost.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_host"),
			"Unknown Salesforce API Host",
			"The provider cannot create the Salesforce API client as there is an unknown configuration value for the Salesforce API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SALESFORCE_API_HOST environment variable.",
		)
	}

	if config.ApiVersion.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_version"),
			"Unknown Salesforce API Version",
			"The provider cannot create the Salesforce API client as there is an unknown configuration value for the Salesforce API version. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SALESFORCE_API_VERSION environment variable.",
		)
	}

	if config.AuthHost.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_host"),
			"Unknown Salesforce Auth Host",
			"The provider cannot create the Salesforce API client as there is an unknown configuration value for the Salesforce auth host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SALESFORCE_AUTH_HOST environment variable.",
		)
	}

	if config.ClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown Salesforce Client ID",
			"The provider cannot create the Salesforce API client as there is an unknown configuration value for the Salesforce client id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SALESFORCE_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown Salesforce Client Secret",
			"The provider cannot create the Salesforce API client as there is an unknown configuration value for the Salesforce client secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SALESFORCE_CLIENT_SECRET environment variable.",
		)
	}

	if config.GrantType.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("grant_type"),
			"Unknown Salesforce Grant Type",
			"The provider cannot create the Salesforce API client as there is an unknown configuration value for the Salesforce grant type. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SALESFORCE_GRANT_TYPE environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Salesforce Username",
			"The provider cannot create the Salesforce API client as there is an unknown configuration value for the Salesforce username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SALESFORCE_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Salesforce Password",
			"The provider cannot create the Salesforce API client as there is an unknown configuration value for the Salesforce password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SALESFORCE_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	apiHost := os.Getenv("SALESFORCE_API_HOST")
	apiVersion := os.Getenv("SALESFORCE_API_VERSION")
	authHost := os.Getenv("SALESFORCE_AUTH_HOST")
	clientID := os.Getenv("SALESFORCE_CLIENT_ID")
	clientSecret := os.Getenv("SALESFORCE_CLIENT_SECRET")
	grantType := os.Getenv("SALESFORCE_GRANT_TYPE")
	username := os.Getenv("SALESFORCE_USERNAME")
	password := os.Getenv("SALESFORCE_PASSWORD")

	if !config.ApiHost.IsNull() {
		apiHost = config.ApiHost.ValueString()
	}

	if !config.ApiVersion.IsNull() {
		apiVersion = config.ApiVersion.ValueString()
	}

	if !config.AuthHost.IsNull() {
		authHost = config.AuthHost.ValueString()
	}

	if !config.ClientID.IsNull() {
		clientID = config.ClientID.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	if !config.GrantType.IsNull() {
		grantType = config.GrantType.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiHost == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_host"),
			"Missing Salesforce API Host",
			"The provider cannot create the Salesforce API client as there is a missing or empty value for the Salesforce API host. "+
				"Set the host value in the configuration or use the SALESFORCE_API_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiVersion == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_version"),
			"Missing Salesforce API Version",
			"The provider cannot create the Salesforce API client as there is a missing or empty value for the Salesforce API version. "+
				"Set the host value in the configuration or use the SALESFORCE_API_VERSION environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if authHost == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_host"),
			"Missing Salesforce Auth Host",
			"The provider cannot create the Salesforce API client as there is a missing or empty value for the Salesforce auth host. "+
				"Set the host value in the configuration or use the SALESFORCE_AUTH_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Salesforce Client ID",
			"The provider cannot create the Salesforce API client as there is a missing or empty value for the Salesforce client id. "+
				"Set the host value in the configuration or use the SALESFORCE_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing Salesforce Client Secret",
			"The provider cannot create the Salesforce API client as there is a missing or empty value for the Salesforce client secret. "+
				"Set the host value in the configuration or use the SALESFORCE_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if grantType == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("grant_type"),
			"Missing Salesforce Grant Type",
			"The provider cannot create the Salesforce API client as there is a missing or empty value for the Salesforce grant type. "+
				"Set the host value in the configuration or use the SALESFORCE_GRANT_TYPE environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Salesforce Username",
			"The provider cannot create the Salesforce API client as there is a missing or empty value for the Salesforce username. "+
				"Set the username value in the configuration or use the SALESFORCE_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Salesforce Password",
			"The provider cannot create the Salesforce API client as there is a missing or empty value for the Salesforce password. "+
				"Set the password value in the configuration or use the SALESFORCE_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "salesforce_apiHost", apiHost)
	ctx = tflog.SetField(ctx, "salesforce_apiVersion", apiVersion)
	ctx = tflog.SetField(ctx, "salesforce_authHost", authHost)
	ctx = tflog.SetField(ctx, "salesforce_clientID", clientID)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "salesforce_clientID")
	ctx = tflog.SetField(ctx, "salesforce_clientSecret", clientSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "salesforce_clientSecret")
	ctx = tflog.SetField(ctx, "salesforce_grantType", grantType)
	ctx = tflog.SetField(ctx, "salesforce_username", username)
	ctx = tflog.SetField(ctx, "salesforce_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "salesforce_password")

	tflog.Debug(ctx, "Creating Salesforce client")

	// Create a new Salesforce client using the configuration values
	client, err := salesforce.NewClient(&apiHost, &apiVersion, &authHost, &clientID, &clientSecret, &grantType, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Salesforce API Client",
			"An unexpected error occurred when creating the Salesforce API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Salesforce Client Error: "+err.Error(),
		)
		return
	}

	// Make the Salesforce client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Salesforce client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *salesforceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDescriptionDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *salesforceProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
