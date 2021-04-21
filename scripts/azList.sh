#!/bin/bash -eo pipefail
az account set -s $AZURE_SUBSCRIPTION_ID
az resource list -g $TEST_RESOURCE_GROUP