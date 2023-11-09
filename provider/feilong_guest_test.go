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
	resource.UnitTest(t, resource.TestCase {
		PreCheck:		func() { testAccPreCheck(t) },
		ProviderFactories:	providerFactories,
		Steps:			[]resource.TestStep {
			{
				Config: testConfig,
				Check:	resource.ComposeTestCheckFunc(
// check that guest is really created!
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "name", "leap"),		// required
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "userid", "LEAP"),		// derived from name: all caps, 8 chars max
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "vcpus", "1"),			// default value
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "memory", "512M"),		// default value
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "disk", "10G"),		// default value
					resource.TestCheckResourceAttr("feilong_guest.opensuse", "image", "opensuse155"),	// required
				),
			},
		},
	})
}
