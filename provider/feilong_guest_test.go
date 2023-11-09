/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testConfig = `
resource "feilong_guest" "opensuse" {
  name = "leap"
  image = "opensuse155"
}
`
)

func TestAccFeilongGuest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfig,
				Check:	resource.ComposeAggregateTestCheckFunc(
// check that guest is really created!
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "name", "leap"),		// required
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "userid", "LEAP"),		// derived from name: all caps, 8 chars max
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "vcpus", "1"),			// default value
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "memory", "512M"),		// default value
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "disk", "10G"),		// default value
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "image", "opensuse155"),	// required
				),
			},
//			// ImportState testing
//			{
//				ResourceName:      "feilong_guest.opensuse",
//				ImportState:       true,
//				ImportStateVerify: true,
//				// This is not normally necessary, but is here because this
//				// example code does not have an actual upstream service.
//				// Once the Read method is able to refresh information from
//				// the upstream service, this can be removed.
//				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
//			},
//			// Update and Read testing
//			{
//				Config: ...
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("feilong_guest.test", "configurable_attribute", "two"),
//				),
//			},
//			// Delete testing automatically occurs in TestCase
		},
	})
}
