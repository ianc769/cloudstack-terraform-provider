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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackZone_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackZoneExists(
						"cloudstack_zone.foo", &zone),
					testAccCheckCloudStackZoneBasicAttributes(&zone),
				),
			},
		},
	})
}

func TestAccCloudStackZone_update(t *testing.T) {
	var zone cloudstack.Zone

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackZone_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackZoneExists(
						"cloudstack_zone.foo", &zone),
					testAccCheckCloudStackZoneBasicAttributes(&zone),
				),
			},
			{
				Config: testAccCloudStackZone_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackZoneExists(
						"cloudstack_zone.foo", &zone),
					testAccCheckCloudStackZoneUpdatedAttributes(&zone),
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
		z, count, err := cs.Zone.GetZoneByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("Zone not found")
		}

		*zone = *z

		return nil
	}
}

func testAccCheckCloudStackZoneBasicAttributes(
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

		if zone.Networktype != "Basic" {
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

		if zone.Dns1 != "8.8.8.8" {
			return fmt.Errorf("Bad DNS1: %s", zone.Dns1)
		}

		if zone.Dns2 != "8.8.4.4" {
			return fmt.Errorf("Bad DNS2: %s", zone.Dns2)
		}

		if zone.Internaldns1 != "8.8.4.4" {
			return fmt.Errorf("Bad internal DNS1: %s", zone.Internaldns1)
		}

		if zone.Networktype != "Basic" {
			return fmt.Errorf("Bad network type: %s", zone.Networktype)
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

		_, count, err := cs.Zone.GetZoneByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("Zone %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackZone_basic = `
resource "cloudstack_zone" "foo" {
  name = "terraform-zone"
  dns1 = "8.8.8.8"
  internal_dns1 = "8.8.4.4"
  network_type = "Basic"
}`

const testAccCloudStackZone_update = `
resource "cloudstack_zone" "foo" {
  name = "terraform-zone-updated"
  dns1 = "8.8.8.8"
  dns2 = "8.8.4.4"
  internal_dns1 = "8.8.4.4"
  network_type = "Basic"
}`
