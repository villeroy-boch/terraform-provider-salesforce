package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Salesforce client is properly configured.
	// It is also possible to use the SALESFORCE_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "salesforce" {
  api_host      = "https://xyz.my.salesforce.com"
  api_version   = "v59.0"
  auth_host     = "https://login.salesforce.com/services/oauth2/token"
  client_id     = "idisjisjisjsfjs"
  client_secret = "sajfaspojfapopsa"
  grant_type    = "password"
  username      = "admin"
  password      = "password123"
}
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"salesforce": providerserver.NewProtocol6WithError(New("test")()),
	}
)
