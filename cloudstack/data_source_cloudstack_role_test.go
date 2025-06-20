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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudStackRole_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudStackRole_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.cloudstack_role.role", "name", "terraform-role"),
					resource.TestCheckResourceAttr(
						"data.cloudstack_role.role", "description", "terraform test role"),
					resource.TestCheckResourceAttr(
						"data.cloudstack_role.role", "is_public", "true"),
				),
			},
		},
	})
}

const testAccDataSourceCloudStackRole_basic = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  description = "terraform test role"
  is_public = true
  type = "User"
}

data "cloudstack_role" "role" {
  filter {
    name = "name"
    value = "${cloudstack_role.foo.name}"
  }
}
`
