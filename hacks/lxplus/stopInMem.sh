#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
IFS=$'\n\t\v'
cd `dirname "${BASH_SOURCE[0]:-$0}"`

# kill all nomad jobs
! /tmp/dune-rc-lxplus/bin/nomad status | awk 'NR>1{print $1}' | xargs -iID /tmp/dune-rc-lxplus/bin/nomad stop ID
! pkill nomad
! pkill consul
# in case you bricked nomad
! killall daq_application