package unittest

import (
	"terraform/tests/utils"
	"testing"

	"github.com/microsoft/terratest-abstraction/unit"
)

var test_variable_values = []struct {
	test_name string
	set_variables bool
	location string
	resource_group_name string
	admin_user string
	admin_password string
	tags string
}{
	{"input-variables", true, "westus", "my_rg", "admin", "admin_password", "{}" },
}
func TestVariables(t *testing.T) {
	for _, tt := range test_variable_values {
		t.Run(tt.test_name, func(t *testing.T) {
			
			location := tt.location
			resource_group_name := tt.resource_group_name
			admin_user := tt.admin_user
			admin_password := tt.admin_password
			tags := tt.tags

			if tt.set_variables {
				utils.SetTfVar("location", location)
				utils.SetTfVar("resource_group_name", resource_group_name)
				utils.SetTfVar("admin_user", admin_user)
				utils.SetTfVar("admin_password", admin_password)
				utils.SetTfVar("tags", tags)
			} else {
				utils.UnsetTfVar("location")
				utils.UnsetTfVar("resource_group_name")
				utils.UnsetTfVar("admin_user")
				utils.UnsetTfVar("admin_password")
				utils.UnsetTfVar("tags")
			}

			RunResourceValidation(location, resource_group_name, admin_user, admin_password, t)
		})
	}
}

func RunResourceValidation(location string, resource_group_name string, admin_user string, admin_password string, t *testing.T) {
	// This is the number of expected Terraform resources being provisioned.
	//
	// Note: There may be more Terraform resources provisioned than Azure resources provisioned!
	expectedTerraformResourceCount := 15

	testFixture := unit.UnitTestFixture{
		GoTest:                          t,
		TfOptions:                       utils.TfOptions,
		ExpectedResourceCount:           expectedTerraformResourceCount,
	}

	unit.RunUnitTests(&testFixture)
}
