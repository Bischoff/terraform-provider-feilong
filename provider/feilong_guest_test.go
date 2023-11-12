/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"testing"
	"os"
	"regexp"
	"errors"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/Bischoff/feilong-client-go"
)

const (
	testConfig = `
resource "feilong_guest" "testacc" {
  name = "testacc"
  image = "testacc"
  mac = "12:34:56:78:9a:bc"
}
`
	rn = "feilong_guest.testacc"
)

// Check the z/VM guest definition
func testCheckZvmGuest() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := os.Getenv("ZVM_CONNECTOR")
		client := feilong.NewClient(&connector, nil)

		result, err := client.ShowGuestDefinition("TESTACC")
		if err != nil {
			return err
		}

		user := regexp.MustCompile(`^USER TESTACC .* 512M`)
		nUser := 0
		command := regexp.MustCompile(`^COMMAND DEFINE CPU`) // appears once per each vCPU
		nCommand := 0
		ipl := regexp.MustCompile(`^IPL 0100`)
		nIpl := 0
		mdisk := regexp.MustCompile(`^MDISK 0100 3390 .* 14564`) // 14564 cyl == 10 GiB
		nMdisk := 0
		nicdef := regexp.MustCompile(`^NICDEF 1000 .* MACID 789ABC`) // xx:xx:xx:78:9a:bc
		nNicdef := 0
		for _, s := range result.Output.UserDirect {
			if (user.MatchString(s)) { nUser++; }
			if (command.MatchString(s)) { nCommand++; }
			if (ipl.MatchString(s)) { nIpl++; }
			if (mdisk.MatchString(s)) { nMdisk++; }
			if (nicdef.MatchString(s)) { nNicdef++; }
		}
		if nUser != 1 || nCommand != 1 || nIpl != 1 || nMdisk != 1 || nNicdef != 1 {
			msg := "Invalid z/VM definition:\n"
			for _, s := range result.Output.UserDirect {
				msg += s + "\n"
			}
			return errors.New(msg)
		}

		return nil
	}
}

func TestAccFeilongGuest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:			func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories:	testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfig,
				Check:	resource.ComposeAggregateTestCheckFunc(				// check resource attributes
					resource.TestCheckResourceAttr(rn, "name", "testacc"),		// required
					resource.TestCheckResourceAttr(rn, "userid", "TESTACC"),	// derived from name: all caps, 8 chars max
					resource.TestCheckResourceAttr(rn, "vcpus", "1"),		// default value
					resource.TestCheckResourceAttr(rn, "memory", "512M"),		// default value
					resource.TestCheckResourceAttr(rn, "disk", "10G"),		// default value
					resource.TestCheckResourceAttr(rn, "image", "testacc"),		// required
					resource.TestCheckResourceAttr(rn, "mac", "12:34:56:78:9a:bc"),	// optional

					testCheckZvmGuest(),						// check that the guest is really created
				),
			},
//			// ImportState testing
//			{
//				ResourceName:		rn,
//				ImportState:		true,
//				ImportStateVerify:	true,
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
//					resource.TestCheckResourceAttr(rn, "configurable_attribute", "two"),
//				),
//			},
//			// Delete testing automatically occurs in TestCase
		},
	})
}
