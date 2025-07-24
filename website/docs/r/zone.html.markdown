---
layout: default
page_title: "CloudStack: cloudstack_zone"
sidebar_current: "docs-cloudstack-resource-zone"
description: |-
    Creates a Zone
---

# CloudStack: cloudstack_zone

A `cloudstack_zone` resource manages a zone within CloudStack.

## Example Usage

```hcl
resource "cloudstack_zone" "example" {
    name = "example-zone"
    dns1 = "8.8.8.8"
    internal_dns1 = "8.8.4.4"
    network_type = "Basic"
    
    # Optional parameters
    dns2 = "1.1.1.1"
    internal_dns2 = "1.0.0.1"
    guest_cidr_address = "10.1.1.0/24"
    local_storage_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the zone.
* `dns1` - (Required) The DNS server 1 for the zone.
* `internal_dns1` - (Required) The internal DNS server 1 for the zone.
* `network_type` - (Required, ForceNew) The type of network to use for the zone. Changing this forces a new resource to be created.
* `dns2` - (Optional) The DNS server 2 for the zone.
* `internal_dns2` - (Optional) The internal DNS server 2 for the zone.
* `ip6_dns1` - (Optional) The IPv6 DNS server 1 for the zone.
* `ip6_dns2` - (Optional) The IPv6 DNS server 2 for the zone.
* `guest_cidr_address` - (Optional) The guest CIDR address for the zone.
* `domain` - (Optional) The domain for the zone.
* `domain_id` - (Optional) The domain ID for the zone.
* `network_domain` - (Optional) The network domain for the zone.
* `local_storage_enabled` - (Optional) Whether local storage is enabled for the zone.
* `security_group_enabled` - (Optional, ForceNew) Whether security groups are enabled for the zone. Changing this forces a new resource to be created.
* `allocation_state` - (Optional) The allocation state of the zone.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the zone.
* `name` - The name of the zone.
* `dns1` - The DNS server 1 for the zone.
* `internal_dns1` - The internal DNS server 1 for the zone.
* `network_type` - The type of network to use for the zone.
* `dns2` - The DNS server 2 for the zone.
* `internal_dns2` - The internal DNS server 2 for the zone.
* `ip6_dns1` - The IPv6 DNS server 1 for the zone.
* `ip6_dns2` - The IPv6 DNS server 2 for the zone.
* `guest_cidr_address` - The guest CIDR address for the zone.
* `domain` - The domain for the zone.
* `domain_id` - The domain ID for the zone.
* `network_domain` - The network domain for the zone.
* `local_storage_enabled` - Whether local storage is enabled for the zone.
* `security_group_enabled` - Whether security groups are enabled for the zone.
* `allocation_state` - The allocation state of the zone.

## Import

Zones can be imported; use `<ZONEID>` as the import ID. For example:

```shell
$ terraform import cloudstack_zone.example <ZONEID>
```
