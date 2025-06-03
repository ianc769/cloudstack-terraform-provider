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
    dns2 = "8.8.4.4"
    internaldns2 = "8.8.8.8"
    guestcidraddress = "10.1.1.0/24"
    localstorageenabled = true
    securitygroupenabled = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the zone.
* `dns1` - (Required) The DNS server  1 for the zone.
* `internal_dns1` - (Required) The internal DNS server  1 for the zone.
* `network_type` - (Required, ForceNew) The type of network to use for the zone.
* `allocationstate` - (Optional) Allocation state of this Zone for allocation of new resources.
* `dns2` - (Optional) The second DNS for the Zone.
* `domain` - (Optional) Network domain name for the networks in the zone.
* `domainid` - (Optional, ForceNew) The ID of the containing domain, null for public zones.
* `guestcidraddress` - (Optional) The guest CIDR address for the Zone.
* `internaldns2` - (Optional) The second internal DNS for the Zone.
* `ip6dns1` - (Optional) The first DNS for IPv6 network in the Zone.
* `ip6dns2` - (Optional) The second DNS for IPv6 network in the Zone.
* `isedge` - (Optional, ForceNew) True if the zone is an edge zone, false otherwise.
* `localstorageenabled` - (Optional) True if local storage offering enabled, false otherwise.
* `securitygroupenabled` - (Optional, ForceNew) True if network is security group enabled, false otherwise.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the zone.
* `name` - The name of the zone.
* `dns1` - The DNS server  1 for the zone.
* `internal_dns1` - The internal DNS server  1 for the zone.
* `network_type` - The type of network to use for the zone.
* `allocationstate` - Allocation state of this Zone for allocation of new resources.
* `dns2` - The second DNS for the Zone.
* `domain` - Network domain name for the networks in the zone.
* `domainid` - The ID of the containing domain, null for public zones.
* `guestcidraddress` - The guest CIDR address for the Zone.
* `internaldns2` - The second internal DNS for the Zone.
* `ip6dns1` - The first DNS for IPv6 network in the Zone.
* `ip6dns2` - The second DNS for IPv6 network in the Zone.
* `localstorageenabled` - True if local storage offering enabled, false otherwise.
* `securitygroupenabled` - True if network is security group enabled, false otherwise.

## Import

Zones can be imported; use `<ZONEID>` as the import ID. For example:

```shell
terraform import cloudstack_zone.example <ZONEID>
```
