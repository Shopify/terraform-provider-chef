package chef

import (
	"log"
  "fmt"

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
				ForceNew: true,
			},

			"run_list": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
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
		Name:        d.Get("name").(string),
		Environment: d.Get("environment").(string),
		ChefType:    "node",
		JsonClass:   "Chef::Node",
		RunList:     run_list,
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
  
