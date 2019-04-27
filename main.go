package main

import (
	"github.com/dominicbarnes/terraform-provider-mongodb/mongodb"
	"github.com/hashicorp/terraform/plugin"
)

//go:generate go run docs/main.go -template=docs/readme.tmpl -output=README.md

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: mongodb.Provider})
}
