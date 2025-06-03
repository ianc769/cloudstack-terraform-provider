//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package cloudstack

import (
	"fmt"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackZone_basic(t *testing.T) {
	var zone cloudstack.Zone

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackZone_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackZoneExists(
						"cloudstack_zone.foo", &zone),
					testAccCheckCloudStackZoneAttributes(&zone),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.foo", "name", "terraform-zone"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.foo", "dns1", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.foo", "internal_dns1", "8.8.4.4"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.foo", "network_type", "Advanced"),
				),
			},
		},
	})
}

func TestAccCloudStackZone_update(t *testing.T) {
	var zone cloudstack.Zone

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackZone_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackZoneExists(
						"cloudstack_zone.foo", &zone),
					testAccCheckCloudStackZoneAttributes(&zone),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.foo", "name", "terraform-zone"),
				),
			},
			{
				Config: testAccCloudStackZone_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackZoneExists(
						"cloudstack_zone.foo", &zone),
					testAccCheckCloudStackZoneUpdatedAttributes(&zone),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.foo", "name", "terraform-zone-updated"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.foo", "dns2", "8.8.4.4"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.foo", "internaldns2", "8.8.8.8"),
				),
			},
		},
	})
}

func testAccCheckCloudStackZoneExists(
	n string, zone *cloudstack.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No zone ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		z, _, err := cs.Zone.GetZoneByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if z.Id != rs.Primary.ID {
			return fmt.Errorf("Zone not found")
		}

		*zone = *z

		return nil
	}
}

func testAccCheckCloudStackZoneAttributes(
	zone *cloudstack.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Name != "terraform-zone" {
			return fmt.Errorf("Bad name: %s", zone.Name)
		}

		if zone.Dns1 != "8.8.8.8" {
			return fmt.Errorf("Bad DNS1: %s", zone.Dns1)
		}

		if zone.Internaldns1 != "8.8.4.4" {
			return fmt.Errorf("Bad internal DNS1: %s", zone.Internaldns1)
		}

		if zone.Networktype != "Advanced" {
			return fmt.Errorf("Bad network type: %s", zone.Networktype)
		}

		return nil
	}
}

func testAccCheckCloudStackZoneUpdatedAttributes(
	zone *cloudstack.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Name != "terraform-zone-updated" {
			return fmt.Errorf("Bad name: %s", zone.Name)
		}

		if zone.Dns2 != "8.8.4.4" {
			return fmt.Errorf("Bad DNS2: %s", zone.Dns2)
		}

		if zone.Internaldns2 != "8.8.8.8" {
			return fmt.Errorf("Bad internal DNS2: %s", zone.Internaldns2)
		}

		return nil
	}
}

func testAccCheckCloudStackZoneDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_zone" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No zone ID is set")
		}

		_, _, err := cs.Zone.GetZoneByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Zone %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func TestAccCloudStackZone_extended(t *testing.T) {
	var zone cloudstack.Zone

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackZone_extended,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackZoneExists(
						"cloudstack_zone.extended", &zone),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.extended", "name", "terraform-zone-extended"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.extended", "dns1", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.extended", "dns2", "8.8.4.4"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.extended", "internal_dns1", "8.8.4.4"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.extended", "internaldns2", "8.8.8.8"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.extended", "network_type", "Advanced"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.extended", "guestcidraddress", "10.1.1.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.extended", "domain", "example.com"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.extended", "localstorageenabled", "true"),
					resource.TestCheckResourceAttr(
						"cloudstack_zone.extended", "securitygroupenabled", "true"),
				),
			},
		},
	})
}

const testAccCloudStackZone_basic = `
resource "cloudstack_zone" "foo" {
  name = "terraform-zone"
  dns1 = "8.8.8.8"
  internal_dns1 = "8.8.4.4"
  network_type = "Advanced"
}
`

const testAccCloudStackZone_update = `
resource "cloudstack_zone" "foo" {
  name = "terraform-zone-updated"
  dns1 = "8.8.8.8"
  dns2 = "8.8.4.4"
  internal_dns1 = "8.8.4.4"
  internaldns2 = "8.8.8.8"
  network_type = "Advanced"
}
`

const testAccCloudStackZone_extended = `
resource "cloudstack_zone" "extended" {
  name = "terraform-zone-extended"
  dns1 = "8.8.8.8"
  dns2 = "8.8.4.4"
  internal_dns1 = "8.8.4.4"
  internaldns2 = "8.8.8.8"
  network_type = "Advanced"
  guestcidraddress = "10.1.1.0/24"
  domain = "example.com"
  localstorageenabled = true
  securitygroupenabled = true
}
`
