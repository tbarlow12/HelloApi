package unittest

import (
	"terraform/tests/utils"
	"testing"
	"fmt"
	"strings"

	"github.com/microsoft/terratest-abstraction/unit"
)

var test_variable_values = []struct {
	test_name string
	set_variables bool
	business_unit string
	app_name string
	stage string
	region string
}{
	{"input-variables", true, "bu", "an", "test", "centralus"},
}
func TestWellnessTraceVariables(t *testing.T) {
	for _, tt := range test_variable_values {
		t.Run(tt.test_name, func(t *testing.T) {
			
			business_unit := tt.business_unit
			app_name := tt.app_name
			stage := tt.stage
			region := tt.region

			if tt.set_variables {
				utils.SetTfVar("BUSINESS_UNIT", business_unit)
				utils.SetTfVar("APP_NAME", app_name)
				utils.SetTfVar("STAGE", stage)
				utils.SetTfVar("REGION", region)
			} else {
				utils.UnsetTfVar("BUSINESS_UNIT")
				utils.UnsetTfVar("APP_NAME")
				utils.UnsetTfVar("STAGE")
				utils.UnsetTfVar("REGION")
			}

			RunResourceValidation(business_unit, app_name, stage, region, t)
		})
	}
}

func RunResourceValidation(business_unit string, app_name string, stage string, region string, t *testing.T) {
	// Needed at runtime
	default_guid := "00000000-0000-0000-0000-000000000000"
	utils.SetTfVar("KUBERNETES_ADMIN_GROUP_ID", default_guid)
	utils.SetTfVar("MYSQL_ADMIN_GROUP_ID", default_guid)

	region_resource_suffix := fmt.Sprintf("%s-%s-%s-%s", business_unit, app_name, stage, region)
	region_resource_suffix_condensed := strings.Replace(region_resource_suffix, "-", "", -1)

	global_resource_suffix := fmt.Sprintf("%s-%s-%s-gbl", business_unit, app_name, stage)
	
	region_resource_group_name := "rg-" + region_resource_suffix
	global_resource_group_name := "rg-" + global_resource_suffix

	// This is a JSON representation of what the Terraform plan should produce if run against a "fresh" environment.
	// Some things to keep in mind:
	//	(1) The key will be the path to the resource. This will match the path found from `terraform state list`. For
	//		this template, it is simply `{RESOURCE_TYPE}.{RESOURCE_NAME}` but this will be more complex if you are using
	//		modules in your deployment!
	//	(2) The test will only verify that the data is specified here exists in the plan. If you omit fields from the
	//		description, they will be ignored.
	resourceDescription := unit.ResourceDescription{
		
		///////////////////////////
		// GLOBAL RESOURCE GROUP //
		///////////////////////////
		
		// Resource name: "traf-endpoint-aks-ge-wt-{stage}-{region}"
		"azurerm_traffic_manager_endpoint.traf-endpoint": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "traf-endpoint-aks-%s"
		}`, global_resource_group_name, region_resource_suffix)),

		"azurerm_traffic_manager_profile.traffic-manager": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "traf-%s"
		}`, global_resource_group_name, global_resource_suffix)),

		// Resource name: "appi-ge-wt-{stage}-{region}"
		"azurerm_application_insights.main": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"location": "%s",
			"name":     "appi-%s"
		}`, global_resource_group_name, region, global_resource_suffix)),
		
		// Resource name: "kv-ge-wt-{stage}-{region}"
		"azurerm_key_vault.main": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"location": "%s",
			"name":     "kv-%s"
		}`, global_resource_group_name, region, global_resource_suffix)),
		
		"azurerm_key_vault_access_policy.terraform": utils.AsMap(t, `{
			"secret_permissions": [
				"backup",
				"delete",
				"get",
				"list",
				"purge",
				"recover",
				"restore",
				"set"
			]
		}`),

		"azurerm_log_analytics_solution.container_insights": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"location": "%s"
		}`, global_resource_group_name, region)),
		
		/////////////////////////////
		// REGIONAL RESOURCE GROUP //
		/////////////////////////////

		// Resource name: "aks-ge-wt-{stage}-{region}"
		"azurerm_kubernetes_cluster.main": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"location": "%s",
			"name":     "aks-%s"
		}`, region_resource_group_name, region, region_resource_suffix)),
		
		//Server name: "mysql-ge-wt-{stage}-{region}"
		"azurerm_mysql_active_directory_administrator.main": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"server_name": "mysql-%s"
		}`, region_resource_group_name, region_resource_suffix)),
		
		//Server name: "mysql-ge-wt-{stage}-{region}"
		"azurerm_mysql_database.wellnesstrace": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"server_name": "mysql-%s",
			"name": "wellnesstrace"
		}`, region_resource_group_name, region_resource_suffix)),
		
		//Server name: "mysql-ge-wt-{stage}-{region}"
		"azurerm_mysql_firewall_rule.allow_all": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"server_name": "mysql-%s",
			"start_ip_address": "0.0.0.0",
			"end_ip_address": "255.255.255.255",
			"name": "AllowAll"
		}`, region_resource_group_name, region_resource_suffix)),
		
		//Server name: "mysql-ge-wt-{stage}-{region}"
		"azurerm_mysql_server.main": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "mysql-%s",
			"public_network_access_enabled": true
		}`, region_resource_group_name, region_resource_suffix)),
		
		// Resource name: "snet-kubernetes-ge-wt-{stage}-{region}"
		"azurerm_mysql_virtual_network_rule.subnet_kubernetes": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "snet-kubernetes-%s"
		}`, region_resource_group_name, region_resource_suffix)),
		
		// Regional resource group. Name "rg-ge-wt-{stage}-{region}"
		"azurerm_resource_group.region": utils.AsMap(t, fmt.Sprintf(`{
			"name": "%s"
		}`, region_resource_group_name)),
		
		// Resource name: "stintgewt{stage}{region}"
		"azurerm_storage_account.internal": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "stint%s"
		}`, region_resource_group_name, region_resource_suffix_condensed)),
		
		// Resource name: "stpubgewt{stage}{region}"
		"azurerm_storage_account.public": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"name": "stpub%s"
		}`, region_resource_group_name, region_resource_suffix_condensed)),
		
		// Resource name: "blobgewt{stage}{region}"
		// Storage account: "stintgewt{stage}{region}"
		"azurerm_storage_blob.internal": utils.AsMap(t, fmt.Sprintf(`{
			"name": "blob%s",
			"storage_account_name": "stint%s"
		}`, region_resource_suffix_condensed, region_resource_suffix_condensed)),
		
		// Resource name: "blobgewt{stage}{region}"
		// Storage account: "stpubgewt{stage}{region}"
		"azurerm_storage_blob.public": utils.AsMap(t, fmt.Sprintf(`{
			"name": "blob%s",
			"storage_account_name": "stpub%s"
		}`, region_resource_suffix_condensed, region_resource_suffix_condensed)),
		
		// Resource name: "contentgewt{stage}{region}"
		// Storage account: "stintgews{stage}{region}"
		"azurerm_storage_container.internal": utils.AsMap(t, fmt.Sprintf(`{
			"name": "content%s",
			"storage_account_name": "stint%s"
		}`, region_resource_suffix_condensed, region_resource_suffix_condensed)),
		
		// Resource name: "contentgewt{stage}{region}"
		// Storage account: "stpubgews{stage}{region}"
		"azurerm_storage_container.public": utils.AsMap(t, fmt.Sprintf(`{
			"name": "content%s",
			"storage_account_name": "stpub%s"
		}`, region_resource_suffix_condensed, region_resource_suffix_condensed)),
		
		// Virtual Network name: "vnet-ge-wt-{stage}-{region}"
		"azurerm_subnet.app_gateway": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"virtual_network_name": "vnet-%s",
			"name": "app-gateway"
		}`, region_resource_group_name, region_resource_suffix)),
		
		// Virtual Network name: "ge-wt-{stage}-{region}"
		"azurerm_subnet.kubernetes": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"virtual_network_name": "vnet-%s",
			"name": "kubernetes"
		}`, region_resource_group_name, region_resource_suffix)),
		
		// Resource name: "ge-wt-{stage}-{region}"
		"azurerm_virtual_network.main": utils.AsMap(t, fmt.Sprintf(`{
			"resource_group_name": "%s",
			"location": "%s",
			"name": "vnet-%s"
		}`, region_resource_group_name, region, region_resource_suffix)),
	}

	// This is the number of expected Terraform resources being provisioned.
	//
	// Note: There may be more Terraform resources provisioned than Azure resources provisioned!
	expectedTerraformResourceCount := 32

	testFixture := unit.UnitTestFixture{
		GoTest:                          t,
		TfOptions:                       utils.TfOptions,
		ExpectedResourceCount:           expectedTerraformResourceCount,
		ExpectedResourceAttributeValues: resourceDescription,
	}

	unit.RunUnitTests(&testFixture)
}
