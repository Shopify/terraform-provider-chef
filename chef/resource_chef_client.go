package chef

import (
	"crypto/sha256"
	"fmt"
	"log"

	chefGo "github.com/go-chef/chef"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceChefClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceChefClientCreate,
		Update: resourceChefClientUpdate,
		Read:   resourceChefClientRead,
		Delete: resourceChefClientDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				StateFunc: func(key interface{}) string {
					switch key.(type) {
					case string:
						return fmt.Sprintf("%x", sha256.Sum256([]byte(key.(string))))
					default:
						return ""
					}
				},
			},
		},
	}
}

func resourceChefClientCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefGo.Client)

	name := d.Get("name").(string)
	admin := d.Get("admin").(bool)

	log.Printf("[DEBUG] create chef client: %s", name)

	api_client, err := client.Clients.Create(name, admin)
	if err != nil {
		return fmt.Errorf("Error creating chef client: %s", err)
	}

	d.SetId(name)
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": name,
	})
	d.Set("private_key", api_client.PrivateKey)

	return resourceChefClientRead(d, meta)
}

func resourceChefClientRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefGo.Client)

	api_client, err := client.Clients.Get(d.Id())
	if err != nil {
		if err.(*chefGo.ErrorResponse).Response.StatusCode == 404 {
			// If the client doesn't exist, that's okay! Set the Id to an empty string
			// and terraform will happily recreate it on the next apply
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("Error reading chef client: %s", err)
		}
	}

	d.Set("name", api_client.Name)
	d.Set("admin", api_client.Admin)
	return nil
}

func resourceChefClientUpdate(d *schema.ResourceData, meta interface{}) error {
	// TODO implement when chef-go adds PUT method for clients

	return resourceChefClientRead(d, meta)
}

func resourceChefClientDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefGo.Client)

	log.Printf("[INFO] Deleting chef client: %s", d.Id())

	if err := client.Clients.Delete(d.Id()); err != nil {
		return fmt.Errorf("Error deleting chef client: %s", err)
	}
	return nil
}
