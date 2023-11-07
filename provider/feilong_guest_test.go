package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFeilongGuest(t *testing.T) {
	t.Skip("resource not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeilongGuest,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"feilong_guest.foo", "sample_attribute", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccFeilongGuest = `
resource "feilong_guest" "foo" {
  sample_attribute = "bar"
}
`
