package mongodb

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/require"
)

var testProviders map[string]terraform.ResourceProvider
var testProvider *schema.Provider

func init() {
	testProvider = Provider().(*schema.Provider)
	testProviders = map[string]terraform.ResourceProvider{
		"mongodb": testProvider,
	}
}

func TestProvider(t *testing.T) {
	require.NoError(t, Provider().(*schema.Provider).InternalValidate())
}

func testProviderPreCheck(t *testing.T) {
	t.Helper()
	require.NoError(t, testProvider.Configure(terraform.NewResourceConfig(nil)))
}
