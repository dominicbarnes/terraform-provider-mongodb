package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/dominicbarnes/terraform-provider-mongodb/mongodb"
	"github.com/hashicorp/terraform/terraform"
)

var tmpl = flag.String("template", "", "template file to render with")
var output = flag.String("output", "", "destination file")

func init() {
	flag.Parse()
}

func main() {
	p := mongodb.Provider()

	// TODO: figure out if/how to avoid specifying all of this by hand
	req := terraform.ProviderSchemaRequest{
		ResourceTypes: []string{"mongodb_document"},
	}
	s, err := p.GetSchema(&req)
	if err != nil {
		log.Fatal(err)
	}

	t := template.Must(template.ParseFiles(*tmpl))

	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(f, s); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("rendered %s to %s\n", *tmpl, *output)
}
