#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
IFS=$'\n\t\v'
cd `dirname "${BASH_SOURCE[0]:-$0}"`

# there is no constraint that these resources are added in any particular order
# but the order presented here makes the most sense

# to demonstrate how services work, the example config includes a reference to
# a service called 'datastorage'
# register this service (stolen from another example) to see it resolve
../external_service/add.sh

# schedule a daq application with nomad (process manager & scheduler)
curl -X POST --fail --data @ru-job.json http://localhost:4646/v1/jobs
echo ""

# put the desired daq config template in the raft key-value store
curl -X PUT --fail --data @my-first-config.json 'http://localhost:8500/v1/kv/daq-applications/configs/my-first-config'
echo ""

# register a config with daq-app-manager
curl -X PUT --fail --data @my-first-daq-app.json 'http://localhost:8500/v1/kv/daq-applications/my-first-daq-app'
echo ""
