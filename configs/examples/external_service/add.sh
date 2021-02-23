#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
IFS=$'\n\t\v'
cd `dirname "${BASH_SOURCE[0]:-$0}"`

# https://learn.hashicorp.com/tutorials/consul/service-registration-external-services#register-an-external-service-with-a-health-check

echo "registering:"
curl --request PUT --data @datastorage-service.json localhost:8500/v1/catalog/register
echo ""
echo "getting:"
curl localhost:8500/v1/catalog/service/datastorage