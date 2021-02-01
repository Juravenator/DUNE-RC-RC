#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
IFS=$'\n\t\v'
cd `dirname "${BASH_SOURCE[0]:-$0}"`

trap "echo received signal && exit" SIGTERM SIGINT

echo "totally real DAQ application here"

while true; do
  echo "tell me what you want to do"
  read command
  echo ""
  echo "going to $command"
  sleep 1
  echo "beep"
  sleep 1
  echo "boop"
  sleep 1
  echo "2021-Feb-01 18:51:12,661 INFO [stdinCommandFacility::completionCallback(...) at /scratch/dingpf/dunedaq-v2.2.0-prep/workdir/sourcecode/cmdlib/plugins/stdinCommandFacility.cpp:91] Command execution resulted with: OK"
  echo ""
done