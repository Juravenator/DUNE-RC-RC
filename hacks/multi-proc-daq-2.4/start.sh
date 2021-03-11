#!/usr/bin/env bash
# we can't use safe mode, shitty scripts being called
# set -o errexit -o nounset -o pipefail
shopt -s expand_aliases # script uses aliases, does't work in non-interactive shells unless explicitly enabled
# shitty scripts
# IFS=$'\n\t\v'
cd `dirname "${BASH_SOURCE[0]:-$0}"`

echo "starting setup"

# tmp_dir=$(mktemp -d -t ci-XXXXXXXXXX)
tmp_dir=/tmp/multiproc-rc-demo
cd $tmp_dir
echo "setting up daq application in $tmp_dir"

echo "sourcing daq built tools"
source daq-buildtools/dbt-setup-env.sh

echo "setting up project"
cd daq-app

echo "setting up runtime environment"
dbt-setup-runtime-environment

echo "starting daq_application $@"
set -o errexit
daq_application "$@"