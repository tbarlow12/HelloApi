locals {
  sql_version      = "12.0"
  sql_db_name      = "sample-db"
  sql_db_collation = "SQL_Latin1_General_CP1_CI_AS"
  sql_db_edition   = "Standard"
  sql_db_sku       = "S2"
  sql_auth_type    = "SQL"
}

resource "azurerm_sql_server" "primary" {
  name                         = "sql-hadr-sample-westus"
  resource_group_name          = azurerm_resource_group.myterraformgroup.name
  location                     = "westus"
  version                      = local.sql_version
  administrator_login          = var.sql_admin_user
  administrator_login_password = var.sql_admin_password

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_mssql_server_extended_auditing_policy" "primary" {
  server_id                               = azurerm_sql_server.primary.id
  storage_endpoint                        = azurerm_storage_account.storage_server.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.storage_server.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 7
}

resource "azurerm_sql_database" "primary" {
  name                             = "sample-db"
  resource_group_name              = azurerm_resource_group.myterraformgroup.name
  location                         = azurerm_sql_server.primary.location
  server_name                      = azurerm_sql_server.primary.name
  collation                        = local.sql_db_collation
  edition                          = local.sql_db_edition
  requested_service_objective_name = local.sql_db_sku

  depends_on = [
    azurerm_sql_server.primary
  ]
}

resource "azurerm_mssql_database_extended_auditing_policy" "primary" {
  database_id                             = azurerm_sql_database.primary.id
  storage_endpoint                        = azurerm_storage_account.storage_server.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.storage_server.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 7
}

resource "azurerm_sql_server" "secondary" {
  name                         = "sql-hadr-sample-eastus"
  resource_group_name          = azurerm_resource_group.myterraformgroup.name
  location                     = "eastus"
  version                      = local.sql_version
  administrator_login          = var.sql_admin_user
  administrator_login_password = var.sql_admin_password

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_sql_server.primary
  ]
}

resource "azurerm_mssql_server_extended_auditing_policy" "secondary" {
  server_id                               = azurerm_sql_server.secondary.id
  storage_endpoint                        = azurerm_storage_account.storage_server.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.storage_server.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 7
}

resource "azurerm_sql_database" "secondary" {
  name                             = "sample-db"
  resource_group_name              = azurerm_resource_group.myterraformgroup.name
  location                         = azurerm_sql_server.secondary.location
  server_name                      = azurerm_sql_server.secondary.name
  collation                        = local.sql_db_collation
  edition                          = local.sql_db_edition
  requested_service_objective_name = local.sql_db_sku

  create_mode        = "OnlineSecondary"
  source_database_id = azurerm_sql_database.primary.id

  depends_on = [
    azurerm_sql_server.secondary
  ]
}

resource "azurerm_sql_failover_group" "main" {
  name                = "sfg-hadr-sample"
  resource_group_name = azurerm_resource_group.myterraformgroup.name
  server_name         = azurerm_sql_server.primary.name
  databases           = [azurerm_sql_database.primary.id]

  partner_servers {
    id = azurerm_sql_server.secondary.id
  }

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }

  depends_on = [
    azurerm_sql_database.secondary
  ]
}
