package mongodb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func resourceMongoDBDocument() *schema.Resource {
	return &schema.Resource{
		Create:   resourceMongoDBDocumentCreate,
		Read:     resourceMongoDBDocumentRead,
		Update:   resourceMongoDBDocumentUpdate,
		Delete:   resourceMongoDBDocumentDelete,
		Exists:   resourceMongoDBDocumentExists,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},

		Schema: map[string]*schema.Schema{
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database for the MongoDB document",
			},
			"collection": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The collection for the MongoDB document",
			},
			"document_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The _id value for the MongoDB document",
			},
			"json": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The raw JSON of the MongoDB document without _id",
				Default:     "{}",
			},
		},
	}
}

func resourceMongoDBDocumentCreate(d *schema.ResourceData, v interface{}) error {
	ctx := context.TODO()
	client := v.(*mongo.Client)

	database := d.Get("database").(string)
	collection := d.Get("collection").(string)
	docid := d.Get("document_id").(string)
	data := []byte(d.Get("json").(string))

	doc := bson.M{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	if docid != "" {
		doc["_id"] = docid
	}

	if res, err := client.Database(database).Collection(collection).InsertOne(ctx, doc); err != nil {
		return err
	} else if id, err := encodeID(res.InsertedID); err != nil {
		return err
	} else {
		d.SetId(id)
	}

	return nil
}

func resourceMongoDBDocumentRead(d *schema.ResourceData, v interface{}) error {
	ctx := context.TODO()
	client := v.(*mongo.Client)

	database := d.Get("database").(string)
	collection := d.Get("collection").(string)
	filter := bson.M{"_id": decodeID(d.Id())}

	var doc bson.M
	if err := client.Database(database).Collection(collection).FindOne(ctx, filter).Decode(&doc); err != nil {
		return err
	}
	delete(doc, "_id")

	if data, err := json.Marshal(doc); err != nil {
		return err
	} else {
		d.Set("json", string(data))
	}

	return nil
}

func resourceMongoDBDocumentUpdate(d *schema.ResourceData, v interface{}) error {
	ctx := context.TODO()
	client := v.(*mongo.Client)

	database := d.Get("database").(string)
	collection := d.Get("collection").(string)
	filter := bson.M{"_id": decodeID(d.Id())}

	var doc bson.M
	if err := json.Unmarshal([]byte(d.Get("json").(string)), &doc); err != nil {
		return err
	}
	delete(doc, "_id")

	if _, err := client.Database(database).Collection(collection).ReplaceOne(ctx, filter, doc); err != nil {
		return err
	}

	return nil
}

func resourceMongoDBDocumentDelete(d *schema.ResourceData, v interface{}) error {
	ctx := context.TODO()
	client := v.(*mongo.Client)

	database := d.Get("database").(string)
	collection := d.Get("collection").(string)
	filter := bson.M{"_id": decodeID(d.Id())}

	if _, err := client.Database(database).Collection(collection).DeleteOne(ctx, filter); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceMongoDBDocumentExists(d *schema.ResourceData, v interface{}) (bool, error) {
	ctx := context.TODO()
	client := v.(*mongo.Client)

	database := d.Get("database").(string)
	collection := d.Get("collection").(string)
	filter := bson.M{"_id": decodeID(d.Id())}

	var m bson.M
	if err := client.Database(database).Collection(collection).FindOne(ctx, filter).Decode(&m); err != nil {
		return false, nil
	}

	return true, nil
}

func encodeID(v interface{}) (string, error) {
	switch value := v.(type) {
	case primitive.ObjectID:
		return value.Hex(), nil
	case string:
		return value, nil
	case fmt.Stringer:
		return value.String(), nil
	default:
		return "", fmt.Errorf("unsupported id type (%[1]T): %[1]v", value)
	}
}

func decodeID(id string) interface{} {
	oid, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		return oid
	}
	return id
}
