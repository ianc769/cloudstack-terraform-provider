---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_zone"
sidebar_current: "docs-cloudstack-cloudstack_zone"
description: |-
  Gets information about cloudstack zone.
---

# cloudstack_zone

Use this datasource to get information about a zone for use in other resources.

## Example Usage

```hcl
data "cloudstack_zone" "zone-data-source" {
  filter {
    name = "name"
    value = "TestZone"
  }
}

# Access zone attributes
output "zone_details" {
  value = {
    name = data.cloudstack_zone.zone-data-source.name
    dns1 = data.cloudstack_zone.zone-data-source.dns1
    dns2 = data.cloudstack_zone.zone-data-source.dns2
    guestcidraddress = data.cloudstack_zone.zone-data-source.guestcidraddress
    securitygroupenabled = data.cloudstack_zone.zone-data-source.securitygroupenabled
  }
}
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the zone.
* `dns1` - The first DNS for the Zone.
* `internal_dns1` - The first internal DNS for the Zone.
* `network_type` - The network type of the zone; can be Basic or Advanced.
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
