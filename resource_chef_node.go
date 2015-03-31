package chef

import (
	"fmt"
	"log"

	chefGo "github.com/go-chef/chef"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func resourceChefNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceChefNodeCreate,
		Update: resourceChefNodeUpdate,
		Read:   resourceChefNodeRead,
		Delete: resourceChefNodeDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"environment": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"run_list": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceChefNodeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefGo.Client)

	var run_list []string
	schema_run_list := d.Get("run_list").(interface{})

	if err := mapstructure.Decode(schema_run_list, &run_list); err != nil {
		return err
	}

	node := chefGo.Node{
		Name:                d.Get("name").(string),
		Environment:         d.Get("environment").(string),
		AutomaticAttributes: map[string]interface{}{},
		NormalAttributes:    map[string]interface{}{},
		DefaultAttributes:   map[string]interface{}{},
		OverrideAttributes:  map[string]interface{}{},
		ChefType:            "node",
		JsonClass:           "Chef::Node",
		RunList:             run_list,
	}

	log.Printf("[DEBUG] node create configuration: %#v", node)

	_, err := client.Nodes.Post(node)
	if err != nil {
		return fmt.Errorf("Error creating chef node: %s", err)
	}

	d.SetId(node.Name)
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": node.Name,
	})

	return resourceChefNodeRead(d, meta)
}

func resourceChefNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefGo.Client)

	node, err := client.Nodes.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error reading chef node: %s", err)
	}

	schema_run_list := make([]interface{}, 0)
	if err := mapstructure.Decode(node.RunList, &schema_run_list); err != nil {
		return err
	}

	d.Set("name", node.Name)
	d.Set("environment", node.Environment)
	d.Set("run_list", schema_run_list)
	return nil
}

func resourceChefNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefGo.Client)

	updateNode := &chefGo.Node{
		AutomaticAttributes: map[string]interface{}{},
		NormalAttributes:    map[string]interface{}{},
		DefaultAttributes:   map[string]interface{}{},
		OverrideAttributes:  map[string]interface{}{},
		ChefType:            "node",
		JsonClass:           "Chef::Node",
  }

	if attr, ok := d.GetOk("name"); ok {
		updateNode.Name = attr.(string)
	}
	if attr, ok := d.GetOk("environment"); ok {
		updateNode.Environment = attr.(string)
	}
	if attr, ok := d.GetOk("run_list"); ok {
		var run_list []string
		if err := mapstructure.Decode(attr.(interface{}), &run_list); err != nil {
			return err
		}
		updateNode.RunList = run_list
	}

	log.Printf("[DEBUG] node update configuration: %#v", updateNode)

	_, err := client.Nodes.Put(*updateNode)
	if err != nil {
		return fmt.Errorf("Failed to update node: %s", err)
	}

	return resourceChefNodeRead(d, meta)
}

func resourceChefNodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefGo.Client)

	log.Printf("[INFO] Deleting node: %s", d.Id())

	err := client.Nodes.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting node: %s", err)
	}

  err = client.magicRequestDecoder("DELETE", "clients/"+d.Id(), nil, nil)
  if err != nil {
    return fmt.Errorf("Error deleting node: %s", err)
  }

	return nil
}
