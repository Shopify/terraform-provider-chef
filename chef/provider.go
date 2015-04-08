package chef

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Chef client name.",
			},
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Absolute path to chef pem key.",
			},
			"server_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Chef server url.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"chef_node": resourceChefNode(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Name:    d.Get("client_name").(string),
		Key:     d.Get("key").(string),
		BaseURL: d.Get("server_url").(string),
	}
	return config.Client()
}
