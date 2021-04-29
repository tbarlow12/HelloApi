#!/bin/bash

# Move into terraform directory
cd terraform
# Initialize workspace
terraform init \
    -backend-config="resource_group_name=${TF_TEST_RESOURCE_GROUP}" \
    -backend-config="storage_account_name=${TF_TEST_STORAGE_ACCOUNT}" \
    -backend-config="container_name=${TF_TEST_CONTAINER_NAME}" \
    -backend-config="key=${TF_TEST_KEY}"
# Set default workspace to start - unit tests create their own workspace
terraform workspace select default
# Unit tests
go test -v $(go list ./... | grep unit)
