# Create subnet
resource "azurerm_subnet" "subnet_server" {
  name                 = "serverSubnet"
  resource_group_name  = azurerm_resource_group.myterraformgroup.name
  virtual_network_name = azurerm_virtual_network.myterraformnetwork.name
  address_prefixes     = ["10.0.1.0/24"]

  depends_on = [
    azurerm_resource_group.myterraformgroup,
    azurerm_virtual_network.myterraformnetwork
  ]
}

# Create Network Security Group
resource "azurerm_network_security_group" "myterraformnsg" {
  name                = "projectNSG"
  location            = var.location
  resource_group_name = azurerm_resource_group.myterraformgroup.name

  tags = var.tags

  depends_on = [
    azurerm_resource_group.myterraformgroup
  ]
}

# Enable inbound and outbound traffic for web deploy
resource "azurerm_network_security_rule" "in80" {
  name                        = "inbound80"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "TCP"
  source_port_range           = "*"
  destination_port_range      = "80"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.myterraformgroup.name
  network_security_group_name = azurerm_network_security_group.myterraformnsg.name
}
resource "azurerm_network_security_rule" "in8172" {
  name                        = "inbound8172"
  priority                    = 1010
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "8172"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.myterraformgroup.name
  network_security_group_name = azurerm_network_security_group.myterraformnsg.name
}
resource "azurerm_network_security_rule" "out80" {
  name                        = "outbound80"
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "TCP"
  source_port_range           = "*"
  destination_port_range      = "80"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.myterraformgroup.name
  network_security_group_name = azurerm_network_security_group.myterraformnsg.name
}
resource "azurerm_network_security_rule" "out8172" {
  name                        = "outbound8172"
  priority                    = 1010
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "8172"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.myterraformgroup.name
  network_security_group_name = azurerm_network_security_group.myterraformnsg.name
}

resource "azurerm_public_ip" "serverPublicIP" {
  name                = "serverPublicIP"
  location            = var.location
  resource_group_name = azurerm_resource_group.myterraformgroup.name
  allocation_method   = "Dynamic"
  domain_name_label   = random_string.fqdn.result

  tags = var.tags

  # TODO set port
}
output "public_ip" {
  value = azurerm_public_ip.serverPublicIP.ip_address
}
output "public_ip_fqdn" {
  value = azurerm_public_ip.serverPublicIP.fqdn
}

# Create network interface on the subnet
resource "azurerm_network_interface" "nic_server" {
  name                = "serverNIC"
  location            = var.location
  resource_group_name = azurerm_resource_group.myterraformgroup.name

  ip_configuration {
    name                          = "nicconfig_server"
    subnet_id                     = azurerm_subnet.subnet_server.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.serverPublicIP.id
  }

  tags = var.tags

  depends_on = [
    azurerm_resource_group.myterraformgroup,
    azurerm_subnet.subnet_server
  ]

  # Exported: applied_dns_servers, id, internal_domain_suffix, mac_address, [private_ip_address, private_ip_addresses], virtual_machine_id
}
output "private_ip_server" {
  value = azurerm_network_interface.nic_server.private_ip_addresses
}
output "public_ip_server" {
  value = azurerm_network_interface.nic_server.ip_configuration
}

# Connect the security group to the network interface
resource "azurerm_network_interface_security_group_association" "nsg2nic_server" {
  network_interface_id      = azurerm_network_interface.nic_server.id
  network_security_group_id = azurerm_network_security_group.myterraformnsg.id

  depends_on = [
    azurerm_network_interface.nic_server,
    azurerm_network_security_group.myterraformnsg
  ]
}

# Create storage account for boot diagnostics
resource "azurerm_storage_account" "storage_server" {
  name                     = "sdiag${random_id.randomId.hex}"
  resource_group_name      = azurerm_resource_group.myterraformgroup.name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = var.tags

  depends_on = [
    azurerm_resource_group.myterraformgroup
  ]
}
resource "azurerm_advanced_threat_protection" "storageThreatDetection" {
  target_resource_id = azurerm_storage_account.storage_server.id
  enabled            = true
}

# Create virtual machine
resource "azurerm_windows_virtual_machine" "vm_server" {
  name                  = "serverVM"
  location              = azurerm_network_interface.nic_server.location
  resource_group_name   = azurerm_resource_group.myterraformgroup.name
  network_interface_ids = [azurerm_network_interface.nic_server.id]
  size                  = "Standard_DS3_v2"

  os_disk {
    name                 = "myOsDiskServer"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  identity {
    type = "SystemAssigned" #SystemAssigned, UserAssigned
    #    identity_ids = []
  }

  computer_name  = "vmserver"
  admin_username = var.admin_user
  admin_password = var.admin_password

  boot_diagnostics {
    storage_account_uri = azurerm_storage_account.storage_server.primary_blob_endpoint
  }

  lifecycle {
    prevent_destroy = true
  }

  tags = var.tags

  depends_on = [
    azurerm_resource_group.myterraformgroup,
    azurerm_network_interface.nic_server,
    azurerm_storage_account.storage_server
  ]
}

resource "azurerm_security_center_assessment_policy" "vmAP" {
  display_name = "VM Access Policy"
  severity     = "Medium"
  description  = "Medium level access policy for the vm"
}

resource "azurerm_security_center_assessment" "vmSCA" {
  assessment_policy_id = azurerm_security_center_assessment_policy.vmAP.id
  target_resource_id   = azurerm_windows_virtual_machine.vm_server.id

  status {
    code = "Healthy"
  }
}

resource "azurerm_virtual_machine_extension" "vm_script_extension" {
  name                 = "hostname"
  virtual_machine_id   = azurerm_windows_virtual_machine.vm_server.id
  publisher            = "Microsoft.Azure.Extensions"
  type                 = "CustomScript"
  type_handler_version = "2.0"

  # https://docs.datadoghq.com/tracing/setup_overview/setup/dotnet-framework/?tab=environmentvariables#windows-services
  settings = <<SETTINGS
    {
        "commandToExecute": "SET COR_ENABLE_PROFILING=1 && SET COR_PROFILER={846F5F1C-F9AE-4B07-969E-05C26BC060D8}"
    }
SETTINGS
}
