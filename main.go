package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/Bischoff/terraform-provider-feilong/provider"
)

var (
	version string = "0.0.1"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts {
		ProviderAddr:	"registry.terraform.io/bischoff/feilong",
		ProviderFunc:	provider.New(version),
		Debug:		debug,
	}

	plugin.Serve(opts)
}
