#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
IFS=$'\n\t\v'
cd `dirname "${BASH_SOURCE[0]:-$0}"`

../../../cli/build/run-control apply daq-app-a.json daq-process-a.json my-first-config.json