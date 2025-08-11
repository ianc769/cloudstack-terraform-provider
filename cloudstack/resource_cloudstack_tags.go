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
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackTags() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackTagsCreate,
		Read:   resourceCloudStackTagsRead,
		Update: resourceCloudStackTagsUpdate,
		Delete: resourceCloudStackTagsDelete,

		Schema: map[string]*schema.Schema{
			"resource_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of resources to tag",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of resource being tagged",
			},

			"tags": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Key/value pairs of tags",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The project to tag resources in",
			},
		},
	}
}

func resourceCloudStackTagsCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	resourceIdsRaw := d.Get("resource_ids").([]interface{})
	resourceIds := make([]string, len(resourceIdsRaw))
	for i, v := range resourceIdsRaw {
		resourceIds[i] = v.(string)
	}

	resourceType := d.Get("resource_type").(string)
	tags := d.Get("tags").(map[string]interface{})

	// Create a new parameter struct
	p := cs.Resourcetags.NewCreateTagsParams(
		resourceIds,
		resourceType,
		tagsFromSchema(tags),
	)

	// Note: CreateTagsParams doesn't implement ProjectIDSetter interface
	// so we can't use setProjectid helper function

	_, err := cs.Resourcetags.CreateTags(p)
	if err != nil {
		return fmt.Errorf("Error creating tags: %s", err)
	}

	// Set the ID to a unique identifier for this tag set
	d.SetId(fmt.Sprintf("%s-%s", resourceType, strings.Join(resourceIds, "-")))

	return resourceCloudStackTagsRead(d, meta)
}

func resourceCloudStackTagsRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	resourceIdsRaw := d.Get("resource_ids").([]interface{})
	if len(resourceIdsRaw) == 0 {
		return fmt.Errorf("no resource IDs found")
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := resourceIdsRaw[0].(string)

	p := cs.Resourcetags.NewListTagsParams()
	p.SetResourceid(resourceId)
	p.SetResourcetype(resourceType)

	// Note: ListTagsParams doesn't have a SetProjectid method
	// We're not using project ID for listing tags as it's not supported by the API

	r, err := cs.Resourcetags.ListTags(p)
	if err != nil {
		return fmt.Errorf("Error listing tags: %s", err)
	}

	if r.Count == 0 {
		log.Printf("[DEBUG] No tags found for %s with ID %s", resourceType, resourceId)
		d.SetId("")
		return nil
	}

	tags := make(map[string]interface{})
	for _, tag := range r.Tags {
		tags[tag.Key] = tag.Value
	}

	d.Set("tags", tags)

	return nil
}

func resourceCloudStackTagsUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	if d.HasChange("tags") {
		resourceIdsRaw := d.Get("resource_ids").([]interface{})
		resourceIds := make([]string, len(resourceIdsRaw))
		for i, v := range resourceIdsRaw {
			resourceIds[i] = v.(string)
		}

		resourceType := d.Get("resource_type").(string)
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})

		remove, create := diffTags(tagsFromSchema(o), tagsFromSchema(n))
		log.Printf("[DEBUG] tags to remove: %v", remove)
		log.Printf("[DEBUG] tags to create: %v", create)

		// First remove any obsolete tags
		if len(remove) > 0 {
			log.Printf("[DEBUG] Removing tags: %v from %s", remove, resourceIds)
			p := cs.Resourcetags.NewDeleteTagsParams(resourceIds, resourceType)
			p.SetTags(remove)

			// Note: DeleteTagsParams doesn't implement ProjectIDSetter interface
			// so we can't use setProjectid helper function

			_, err := cs.Resourcetags.DeleteTags(p)
			if err != nil {
				return fmt.Errorf("Error deleting tags: %s", err)
			}
		}

		// Then add any new tags
		if len(create) > 0 {
			log.Printf("[DEBUG] Creating tags: %v for %s", create, resourceIds)
			p := cs.Resourcetags.NewCreateTagsParams(resourceIds, resourceType, create)

			// Note: CreateTagsParams doesn't implement ProjectIDSetter interface
			// so we can't use setProjectid helper function

			_, err := cs.Resourcetags.CreateTags(p)
			if err != nil {
				return fmt.Errorf("Error creating tags: %s", err)
			}
		}
	}

	return resourceCloudStackTagsRead(d, meta)
}

func resourceCloudStackTagsDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	resourceIdsRaw := d.Get("resource_ids").([]interface{})
	resourceIds := make([]string, len(resourceIdsRaw))
	for i, v := range resourceIdsRaw {
		resourceIds[i] = v.(string)
	}

	resourceType := d.Get("resource_type").(string)
	tags := d.Get("tags").(map[string]interface{})

	p := cs.Resourcetags.NewDeleteTagsParams(resourceIds, resourceType)
	p.SetTags(tagsFromSchema(tags))

	// Note: DeleteTagsParams doesn't implement ProjectIDSetter interface
	// so we can't use setProjectid helper function

	_, err := cs.Resourcetags.DeleteTags(p)
	if err != nil {
		return fmt.Errorf("Error deleting tags: %s", err)
	}

	return nil
}
