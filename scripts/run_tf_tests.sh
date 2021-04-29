#!/bin/bash

# Move into terraform directory
cd terraform
# Initialize workspace
terraform init
# Set default workspace to start - unit tests create their own workspace
terraform workspace select default
# Unit tests
go test -v $(go list ./... | grep unit)
