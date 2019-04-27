package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/dominicbarnes/got"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type testStep struct {
	Config         string             `testdata:"config.tf"`
	TerraformState testTerraformState `testdata:"tfstate.json"`
}

type testTerraformState map[string]struct {
	AttributesSet   []string          `json:"attributes_set"`
	AttributesUnset []string          `json:"attributes_unset"`
	Attributes      map[string]string `json:"attributes"`
}

func TestProviderDocument(t *testing.T) {
	var steps []resource.TestStep
	eachdir(t, "testdata/resource_mongodb_document", func(fixture, dir string) {
		t.Run(fixture, func(t *testing.T) {
			eachdir(t, filepath.Join(dir, fixture), func(step, dir string) {
				var teststep testStep
				got.TestData(t, filepath.Join(dir, step), &teststep)

				var checks []resource.TestCheckFunc

				for id, r := range teststep.TerraformState {
					for _, k := range r.AttributesSet {
						checks = append(checks, resource.TestCheckResourceAttrSet(id, k))
					}

					for _, k := range r.AttributesUnset {
						checks = append(checks, resource.TestCheckNoResourceAttr(id, k))
					}

					for k, v := range r.Attributes {
						checks = append(checks, resource.TestCheckResourceAttr(id, k, v))
					}
				}

				checks = append(checks, checkMongoDB)

				steps = append(steps, resource.TestStep{
					Config: teststep.Config,
					Check:  resource.ComposeTestCheckFunc(checks...),
				})
			})

			resource.Test(t, resource.TestCase{
				PreCheck:  func() { testProviderPreCheck(t) },
				Providers: testProviders,
				Steps:     steps,
			})
		})
	})
}

func eachdir(t *testing.T, dir string, fn func(name, dir string)) {
	t.Helper()

	ff, err := ioutil.ReadDir(dir)
	require.NoError(t, err)

	for _, f := range ff {
		if f.IsDir() {
			fn(f.Name(), dir)
		}
	}
}

func checkMongoDB(state *terraform.State) error {
	ctx := context.TODO()
	mongodb := testProvider.Meta().(*mongo.Client)

	// find all the mongodb_document resources in the state and ensure that
	// they are properly formed in the database
	counts := make(map[[2]string]int64)
	for _, mod := range state.Modules {
		for _, res := range mod.Resources {
			if res.Type == "mongodb_document" {
				db := res.Primary.Attributes["database"]
				c := res.Primary.Attributes["collection"]
				id := decodeID(res.Primary.ID)

				expected := make(map[string]interface{})
				if err := json.Unmarshal([]byte(res.Primary.Attributes["json"]), &expected); err != nil {
					return err
				}
				actual := make(map[string]interface{})
				if err := mongodb.Database(db).Collection(c).FindOne(ctx, bson.M{"_id": id}).Decode(&actual); err != nil {
					return err
				}
				delete(actual, "_id")
				if !equalJSON(expected, actual) {
					return fmt.Errorf("document in mongo %+v does not match expected %+v", actual, expected)
				}

				counts[[2]string{db, c}]++
			}
		}
	}

	// make sure the number of documents in each collection is correct to
	// ensure that we aren't leaving documents behind or creating more than
	// we should.
	for dbc, expected := range counts {
		db := dbc[0]
		c := dbc[1]
		actual, err := mongodb.Database(db).Collection(c).CountDocuments(ctx, bson.M{})
		if err != nil {
			return err
		}
		if expected != actual {
			return fmt.Errorf("incorrect number of documents in %s.%s", db, c)
		}
	}

	return nil
}

func equalJSON(expected, actual interface{}) bool {
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		return false
	}

	actualJSON, err := json.Marshal(expected)
	if err != nil {
		return false
	}

	return string(expectedJSON) == string(actualJSON)
}
