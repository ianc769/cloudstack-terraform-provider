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
	"log"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackZoneCreate,
		Read:   resourceCloudStackZoneRead,
		Update: resourceCloudStackZoneUpdate,
		Delete: resourceCloudStackZoneDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_dns1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Optional parameters
			"dns2": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"internal_dns2": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip6_dns1": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip6_dns2": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"guest_cidr_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"local_storage_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"security_group_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allocation_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackZoneCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	dns1 := d.Get("dns1").(string)
	internal_dns1 := d.Get("internal_dns1").(string)
	network_type := d.Get("network_type").(string)

	// Create a new parameter struct
	p := cs.Zone.NewCreateZoneParams(dns1, internal_dns1, name, network_type)

	// Set optional parameters
	if dns2, ok := d.GetOk("dns2"); ok {
		p.SetDns2(dns2.(string))
	}

	if internal_dns2, ok := d.GetOk("internal_dns2"); ok {
		p.SetInternaldns2(internal_dns2.(string))
	}

	if ip6_dns1, ok := d.GetOk("ip6_dns1"); ok {
		p.SetIp6dns1(ip6_dns1.(string))
	}

	if ip6_dns2, ok := d.GetOk("ip6_dns2"); ok {
		p.SetIp6dns2(ip6_dns2.(string))
	}

	if guest_cidr_address, ok := d.GetOk("guest_cidr_address"); ok {
		p.SetGuestcidraddress(guest_cidr_address.(string))
	}

	if domain, ok := d.GetOk("domain"); ok {
		p.SetDomain(domain.(string))
	}

	if domain_id, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(domain_id.(string))
	}

	// Note: The SetNetworkdomain method is not available in the SDK
	// We'll need to check the SDK for the correct method name or parameter

	if local_storage_enabled, ok := d.GetOkExists("local_storage_enabled"); ok {
		p.SetLocalstorageenabled(local_storage_enabled.(bool))
	}

	if security_group_enabled, ok := d.GetOkExists("security_group_enabled"); ok {
		p.SetSecuritygroupenabled(security_group_enabled.(bool))
	}

	if allocation_state, ok := d.GetOk("allocation_state"); ok {
		p.SetAllocationstate(allocation_state.(string))
	}

	log.Printf("[DEBUG] Creating Zone %s", name)
	n, err := cs.Zone.CreateZone(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Zone %s successfully created", name)
	d.SetId(n.Id)

	return resourceCloudStackZoneRead(d, meta)
}

func resourceCloudStackZoneRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Zone %s", d.Get("name").(string))

	// Get the Zone details
	z, count, err := cs.Zone.GetZoneByName(d.Get("name").(string))

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Zone %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(z.Id)
	d.Set("name", z.Name)
	d.Set("dns1", z.Dns1)
	d.Set("internal_dns1", z.Internaldns1)
	d.Set("network_type", z.Networktype)

	// Set optional parameters
	d.Set("dns2", z.Dns2)
	d.Set("internal_dns2", z.Internaldns2)
	d.Set("ip6_dns1", z.Ip6dns1)
	d.Set("ip6_dns2", z.Ip6dns2)
	d.Set("guest_cidr_address", z.Guestcidraddress)
	d.Set("domain", z.Domain)
	d.Set("domain_id", z.Domainid)
	// Note: The Networkdomain field is not available in the SDK
	d.Set("local_storage_enabled", z.Localstorageenabled)
	d.Set("security_group_enabled", z.Securitygroupsenabled)
	d.Set("allocation_state", z.Allocationstate)

	return nil
}

func resourceCloudStackZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Zone.NewUpdateZoneParams(d.Id())

	// Check if the name has changed
	if d.HasChange("name") {
		p.SetName(d.Get("name").(string))
	}

	// Check if dns1 has changed
	if d.HasChange("dns1") {
		p.SetDns1(d.Get("dns1").(string))
	}

	// Check if dns2 has changed
	if d.HasChange("dns2") {
		p.SetDns2(d.Get("dns2").(string))
	}

	// Check if internal_dns1 has changed
	if d.HasChange("internal_dns1") {
		p.SetInternaldns1(d.Get("internal_dns1").(string))
	}

	// Check if internal_dns2 has changed
	if d.HasChange("internal_dns2") {
		p.SetInternaldns2(d.Get("internal_dns2").(string))
	}

	// Check if ip6_dns1 has changed
	if d.HasChange("ip6_dns1") {
		p.SetIp6dns1(d.Get("ip6_dns1").(string))
	}

	// Check if ip6_dns2 has changed
	if d.HasChange("ip6_dns2") {
		p.SetIp6dns2(d.Get("ip6_dns2").(string))
	}

	// Check if guest_cidr_address has changed
	if d.HasChange("guest_cidr_address") {
		p.SetGuestcidraddress(d.Get("guest_cidr_address").(string))
	}

	// Check if allocation_state has changed
	if d.HasChange("allocation_state") {
		p.SetAllocationstate(d.Get("allocation_state").(string))
	}

	// Update the zone
	log.Printf("[DEBUG] Updating zone %s", d.Get("name").(string))
	_, err := cs.Zone.UpdateZone(p)
	if err != nil {
		return fmt.Errorf("Error updating zone %s: %s", d.Get("name").(string), err)
	}

	return resourceCloudStackZoneRead(d, meta)
}

func resourceCloudStackZoneDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Zone.NewDeleteZoneParams(d.Id())
	_, err := cs.Zone.DeleteZone(p)

	if err != nil {
		return fmt.Errorf("Error deleting Zone: %s", err)
	}

	return nil
}
