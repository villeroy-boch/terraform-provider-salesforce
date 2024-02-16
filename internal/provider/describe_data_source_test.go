package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDescriptionsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "salesforce_description" "test" {
					name = "test"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify metadata
					resource.TestCheckResourceAttr("data.salesforce_description.test", "name", "test"),
					resource.TestCheckResourceAttr("data.salesforce_description.test", "label", "Training Course"),

					// Verify number of columns returned
					resource.TestCheckResourceAttr("data.salesforce_description.test", "columns.#", "2"),
					// Verify the first column to ensure all attributes are set
					resource.TestCheckResourceAttr("data.salesforce_description.test", "columns.0.name", "OwnerId"),
					resource.TestCheckResourceAttr("data.salesforce_description.test", "columns.0.label", "Owner ID"),
					resource.TestCheckResourceAttr("data.salesforce_description.test", "columns.0.type", "reference"),

					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.salesforce_description.test", "id", "placeholder"),
				),
			},
		},
	})
}
