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

func TestAccCloudStackTags_basic(t *testing.T) {
	var vm cloudstack.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackTagsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackTags_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackInstanceExists(
						"cloudstack_instance.foobar", &vm),
					testAccCheckCloudStackTagsExists(
						"cloudstack_tags.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_tags.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr(
						"cloudstack_tags.foo", "tags.terraform-tag", "true"),
					resource.TestCheckResourceAttr(
						"cloudstack_tags.foo", "tags.tag2", "value2"),
				),
			},
		},
	})
}

func TestAccCloudStackTags_update(t *testing.T) {
	var vm cloudstack.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackTagsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackTags_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackInstanceExists(
						"cloudstack_instance.foobar", &vm),
					testAccCheckCloudStackTagsExists(
						"cloudstack_tags.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_tags.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr(
						"cloudstack_tags.foo", "tags.terraform-tag", "true"),
					resource.TestCheckResourceAttr(
						"cloudstack_tags.foo", "tags.tag2", "value2"),
				),
			},
			{
				Config: testAccCloudStackTags_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackInstanceExists(
						"cloudstack_instance.foobar", &vm),
					testAccCheckCloudStackTagsExists(
						"cloudstack_tags.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_tags.foo", "tags.%", "3"),
					resource.TestCheckResourceAttr(
						"cloudstack_tags.foo", "tags.terraform-tag", "true"),
					resource.TestCheckResourceAttr(
						"cloudstack_tags.foo", "tags.tag2", "value2-updated"),
					resource.TestCheckResourceAttr(
						"cloudstack_tags.foo", "tags.tag3", "value3"),
				),
			},
		},
	})
}

func testAccCheckCloudStackTagsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No tags ID is set")
		}

		return nil
	}
}

func testAccCheckCloudStackTagsDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_tags" {
			continue
		}

		// Tags are automatically destroyed when the resource is destroyed
		// so there's nothing to check here
	}

	return nil
}

var testAccCloudStackTags_basic = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-test"
  service_offering= "Small Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "cloudstack_tags" "foo" {
  resource_ids  = [cloudstack_instance.foobar.id]
  resource_type = "UserVm"
  
  tags = {
    terraform-tag = "true"
    tag2          = "value2"
  }
}
`

var testAccCloudStackTags_update = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-test"
  service_offering= "Small Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "cloudstack_tags" "foo" {
  resource_ids  = [cloudstack_instance.foobar.id]
  resource_type = "UserVm"
  
  tags = {
    terraform-tag = "true"
    tag2          = "value2-updated"
    tag3          = "value3"
  }
}
`
