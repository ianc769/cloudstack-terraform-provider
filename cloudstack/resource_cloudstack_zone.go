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
				ForceNew: true,
			},
			"allocationstate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domainid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"guestcidraddress": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internaldns2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip6dns1": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip6dns2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isedge": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"localstorageenabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"securitygroupenabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
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
	if allocationstate, ok := d.GetOk("allocationstate"); ok {
		p.SetAllocationstate(allocationstate.(string))
	}

	if dns2, ok := d.GetOk("dns2"); ok {
		p.SetDns2(dns2.(string))
	}

	if domain, ok := d.GetOk("domain"); ok {
		p.SetDomain(domain.(string))
	}

	if domainid, ok := d.GetOk("domainid"); ok {
		p.SetDomainid(domainid.(string))
	}

	if guestcidraddress, ok := d.GetOk("guestcidraddress"); ok {
		p.SetGuestcidraddress(guestcidraddress.(string))
	}

	if internaldns2, ok := d.GetOk("internaldns2"); ok {
		p.SetInternaldns2(internaldns2.(string))
	}

	if ip6dns1, ok := d.GetOk("ip6dns1"); ok {
		p.SetIp6dns1(ip6dns1.(string))
	}

	if ip6dns2, ok := d.GetOk("ip6dns2"); ok {
		p.SetIp6dns2(ip6dns2.(string))
	}

	if isedge, ok := d.GetOkExists("isedge"); ok {
		p.SetIsedge(isedge.(bool))
	}

	if localstorageenabled, ok := d.GetOkExists("localstorageenabled"); ok {
		p.SetLocalstorageenabled(localstorageenabled.(bool))
	}

	if securitygroupenabled, ok := d.GetOkExists("securitygroupenabled"); ok {
		p.SetSecuritygroupenabled(securitygroupenabled.(bool))
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

	// Set optional fields
	d.Set("allocationstate", z.Allocationstate)
	d.Set("dns2", z.Dns2)
	d.Set("domain", z.Domain)
	d.Set("domainid", z.Domainid)
	d.Set("guestcidraddress", z.Guestcidraddress)
	d.Set("internaldns2", z.Internaldns2)
	d.Set("ip6dns1", z.Ip6dns1)
	d.Set("ip6dns2", z.Ip6dns2)
	d.Set("localstorageenabled", z.Localstorageenabled)
	d.Set("securitygroupenabled", z.Securitygroupsenabled)

	return nil
}

func resourceCloudStackZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Zone.NewUpdateZoneParams(d.Id())

	// Check for changes and update parameters
	if d.HasChange("allocationstate") {
		p.SetAllocationstate(d.Get("allocationstate").(string))
	}

	if d.HasChange("dns1") {
		p.SetDns1(d.Get("dns1").(string))
	}

	if d.HasChange("dns2") {
		p.SetDns2(d.Get("dns2").(string))
	}

	if d.HasChange("domain") {
		p.SetDomain(d.Get("domain").(string))
	}

	if d.HasChange("guestcidraddress") {
		p.SetGuestcidraddress(d.Get("guestcidraddress").(string))
	}

	if d.HasChange("internal_dns1") {
		p.SetInternaldns1(d.Get("internal_dns1").(string))
	}

	if d.HasChange("internaldns2") {
		p.SetInternaldns2(d.Get("internaldns2").(string))
	}

	if d.HasChange("ip6dns1") {
		p.SetIp6dns1(d.Get("ip6dns1").(string))
	}

	if d.HasChange("ip6dns2") {
		p.SetIp6dns2(d.Get("ip6dns2").(string))
	}

	if d.HasChange("localstorageenabled") {
		p.SetLocalstorageenabled(d.Get("localstorageenabled").(bool))
	}

	if d.HasChange("name") {
		p.SetName(d.Get("name").(string))
	}

	// Update the zone
	log.Printf("[DEBUG] Updating Zone %s", d.Get("name").(string))
	_, err := cs.Zone.UpdateZone(p)

	if err != nil {
		return fmt.Errorf("Error updating Zone: %s", err)
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
