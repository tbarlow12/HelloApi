terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.52.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_security_center_auto_provisioning" "security_center_auto_provision" {
  auto_provision = "On"
}

# Generate random text for a unique storage account name
resource "random_id" "randomId" {
  keepers = {
    # Generate a new ID only when a new resource group is defined
    resource_group = azurerm_resource_group.myterraformgroup.name
  }

  byte_length = 8

  depends_on = [
    azurerm_resource_group.myterraformgroup
  ]
}

resource "random_string" "fqdn" {
  length  = 6
  special = false
  upper   = false
  number  = false
}

