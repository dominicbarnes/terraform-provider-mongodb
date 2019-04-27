package mongodb

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_HOSTNAME", "localhost"),
				Description: "The hostname of MongoDB server to connect to (default 'localhost')",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_PORT", 27017),
				Description: "The port number of the MongoDB server to connect to (default 27017)",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_USERNAME", nil),
				Description: "The name of the user connecting to MongoDB",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODB_PASSWORD", nil),
				Description: "The password of the user connecting to MongoDB",
				Sensitive:   true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"mongodb_document": resourceMongoDBDocument(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	user := d.Get("username").(string)
	pass := d.Get("password").(string)
	app := fmt.Sprintf("Terraform v%s", terraform.VersionString())

	var ui *url.Userinfo
	if pass != "" {
		ui = url.UserPassword(user, pass)
	} else if user != "" {
		ui = url.User(user)
	}

	uri := url.URL{
		Scheme: "mongodb",
		Host:   fmt.Sprintf("%s:%d", d.Get("hostname"), d.Get("port")),
		User:   ui,
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri.String()).SetAppName(app))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a mongodb client")
	}

	if err := client.Connect(context.TODO()); err != nil {
		return nil, errors.Wrap(err, "failed to ping mongodb")
	}

	return client, nil
}
