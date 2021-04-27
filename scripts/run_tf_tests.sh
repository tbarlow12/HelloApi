cd terraform
terraform init
terraform workspace select default
# Unit tests
go test -v $(go list ./... | grep unit)
