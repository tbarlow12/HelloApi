package unittest

import (
	"fmt"
	"terraform/tests/utils"
	"testing"

	"github.com/microsoft/terratest-abstraction/unit"
)

var test_variable_values = []struct {
	test_name           string
	set_variables       bool
	location            string
	resource_group_name string
	admin_user          string
	admin_password      string
	sql_admin_user      string
	sql_admin_password  string
	tags                string
}{
	{"input-variables", true, "westus", "my_rg", "admin", "admin_password", "sql-admin", "sql_admin_password", "{}"},
}

func TestVariables(t *testing.T) {
	for _, tt := range test_variable_values {
		t.Run(tt.test_name, func(t *testing.T) {

			location := tt.location
			resource_group_name := tt.resource_group_name
			admin_user := tt.admin_user
			admin_password := tt.admin_password
			sql_admin_user := tt.sql_admin_user
			sql_admin_password := tt.sql_admin_password
			tags := tt.tags

			if tt.set_variables {
				utils.SetTfVar("location", location)
				utils.SetTfVar("resource_group_name", resource_group_name)
				utils.SetTfVar("admin_user", admin_user)
				utils.SetTfVar("admin_password", admin_password)
				utils.SetTfVar("sql_admin_user", sql_admin_user)
				utils.SetTfVar("sql_admin_password", sql_admin_password)
				utils.SetTfVar("tags", tags)
			} else {
				utils.UnsetTfVar("location")
				utils.UnsetTfVar("resource_group_name")
				utils.UnsetTfVar("admin_user")
				utils.UnsetTfVar("admin_password")
				utils.UnsetTfVar("sql_admin_user")
				utils.UnsetTfVar("sql_admin_password")
				utils.UnsetTfVar("tags")
			}

			RunResourceValidation(location, resource_group_name, admin_user, admin_password, sql_admin_user, sql_admin_password, t)
		})
	}
}

func RunResourceValidation(location string, resource_group_name string, admin_user string, admin_password string, sql_admin_user string, sql_admin_password string, t *testing.T) {

	resourceDescription := unit.ResourceDescription{

		"azurerm_subnet.subnet_server": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "serverSubnet"
		}`, resource_group_name)),

		"azurerm_network_security_group.myterraformnsg": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "projectNSG"
		}`, resource_group_name)),

		"azurerm_network_security_rule.in80": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "inbound80"
		}`, resource_group_name)),

		"azurerm_network_security_rule.in8172": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "inbound8172"
		}`, resource_group_name)),

		"azurerm_network_security_rule.out80": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "outbound80"
		}`, resource_group_name)),

		"azurerm_network_security_rule.out8172": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "outbound8172"
		}`, resource_group_name)),

		"azurerm_public_ip.serverPublicIP": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "serverPublicIP"
		}`, resource_group_name)),

		"azurerm_storage_account.storage_server": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s"
		}`, resource_group_name)),

		"azurerm_network_interface.nic_server": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "serverNIC"
		}`, resource_group_name)),

		"azurerm_advanced_threat_protection.storageThreatDetection": utils.AsMap(t, `{
			"enabled": true
		}`),

		"azurerm_windows_virtual_machine.vm_server": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "serverVM"
		}`, resource_group_name)),

		"azurerm_security_center_assessment_policy.vmAP": utils.AsMap(t, `{
			"display_name": "VM Access Policy"
		}`),

		// "azurerm_key_vault.projectKeyVault": utils.AsMap(t, fmt.Sprintf(`{
		// 	"resource_group_name": "%s",
		// 	"name": "projectKeyVault"
		// }`, resource_group_name)),

		"azurerm_virtual_network.myterraformnetwork": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "projectVNet"
		}`, resource_group_name)),
	}

	// This is the number of expected Terraform resources being provisioned.
	//
	// Note: There may be more Terraform resources provisioned than Azure resources provisioned!
	expectedTerraformResourceCount := 30

	testFixture := unit.UnitTestFixture{
		GoTest:                          t,
		TfOptions:                       utils.TfOptions,
		ExpectedResourceCount:           expectedTerraformResourceCount,
		ExpectedResourceAttributeValues: resourceDescription,
	}

	unit.RunUnitTests(&testFixture)
}
