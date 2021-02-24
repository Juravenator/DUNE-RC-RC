#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
IFS=$'\n\t\v'
cd `dirname "${BASH_SOURCE[0]:-$0}"`

curl 'http://localhost:8500/v1/kv/daq-applications/my-first-daq-app?raw=' 2>/dev/null | jq