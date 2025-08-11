---
subcategory: "CloudStack"
layout: "cloudstack"
page_title: "CloudStack: cloudstack_tags"
description: |-
  Manages tags for CloudStack resources.
---

# cloudstack_tags

Manages tags for CloudStack resources.

## Example Usage

```hcl
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "cloudstack_instance" "web" {
  name = "server-1"
  display_name = "server-1"
  service_offering = "Small Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "cloudstack_tags" "web_tags" {
  resource_ids  = [cloudstack_instance.web.id]
  resource_type = "UserVm"
  
  tags = {
    environment = "production"
    role        = "webserver"
    managed-by  = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_ids` - (Required, ForceNew) List of resource IDs to tag.
* `resource_type` - (Required, ForceNew) Type of resource being tagged. Valid values include: UserVm, Template, ISO, Volume, Snapshot, Network, VPC, etc.
* `tags` - (Required) Key/value pairs of tags.
* `project` - (Optional, ForceNew) The project to tag resources in.

## Attributes Reference

No additional attributes are exported.

## Import

Tags cannot be imported at this time.
