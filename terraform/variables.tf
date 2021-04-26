variable "BUSINESS_UNIT" {
  type = string
  default = "az"
}

variable "APP_NAME" {
  type = string
  default = "tf"
}

variable "REGIONS" {
    type = list(string)
    default = ["westus2","eastus2"]
}