package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/shopify/terraform-provider-chef/chef"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: chef.Provider,
	})
}
