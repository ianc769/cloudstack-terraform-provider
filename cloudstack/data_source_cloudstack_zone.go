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
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudStackZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackZoneRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns1": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_dns1": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allocationstate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns2": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domainid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guestcidraddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internaldns2": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip6dns1": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip6dns2": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"localstorageenabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"securitygroupenabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudstackZoneRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Zone.NewListZonesParams()
	csZones, err := cs.Zone.ListZones(p)

	if err != nil {
		return fmt.Errorf("Failed to list zones: %s", err)
	}
	filters := d.Get("filter")
	var zone *cloudstack.Zone

	for _, z := range csZones.Zones {
		match, err := applyZoneFilters(z, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			zone = z
		}
	}

	if zone == nil {
		return fmt.Errorf("No zone is matching with the specified regex")
	}
	log.Printf("[DEBUG] Selected zone: %s\n", zone.Name)

	return zoneDescriptionAttributes(d, zone)
}

func zoneDescriptionAttributes(d *schema.ResourceData, zone *cloudstack.Zone) error {
	d.SetId(zone.Id)
	d.Set("name", zone.Name)
	d.Set("dns1", zone.Dns1)
	d.Set("internal_dns1", zone.Internaldns1)
	d.Set("network_type", zone.Networktype)

	// Set optional fields
	d.Set("allocationstate", zone.Allocationstate)
	d.Set("dns2", zone.Dns2)
	d.Set("domain", zone.Domain)
	d.Set("domainid", zone.Domainid)
	d.Set("guestcidraddress", zone.Guestcidraddress)
	d.Set("internaldns2", zone.Internaldns2)
	d.Set("ip6dns1", zone.Ip6dns1)
	d.Set("ip6dns2", zone.Ip6dns2)
	d.Set("localstorageenabled", zone.Localstorageenabled)
	d.Set("securitygroupenabled", zone.Securitygroupsenabled)

	return nil
}

func applyZoneFilters(zone *cloudstack.Zone, filters *schema.Set) (bool, error) {
	var zoneJSON map[string]interface{}
	k, _ := json.Marshal(zone)
	err := json.Unmarshal(k, &zoneJSON)
	if err != nil {
		return false, err
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})
		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("Invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(m["name"].(string), "_", "")
		zoneField := zoneJSON[updatedName].(string)
		if !r.MatchString(zoneField) {
			return false, nil
		}

	}
	return true, nil
}
