package utils

import (
	"encoding/json"
	"testing"
	"os"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// AsMap parses block of JSON into a Go Map. Fails the test if the JSON is invalid.
func AsMap(t *testing.T, jsonString string) map[string]interface{} {
	var theMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &theMap); err != nil {
		t.Fatal(err)
	}
	return theMap
}

func SetTfVar(key string, value string) {
	os.Setenv("TF_VAR_" + key, value)
}

func UnsetTfVar(key string) {
	os.Unsetenv("TF_VAR_" + key)
}

// TfOptions TF options that can be used by all tests
var TfOptions = &terraform.Options{
	TerraformDir: "../../",
	Upgrade:      true,
	Vars:         map[string]interface{}{},
}