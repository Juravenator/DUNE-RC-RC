#!/usr/bin/env bash
# we can't use safe mode, shitty scripts being called
# set -o errexit -o nounset -o pipefail
shopt -s expand_aliases # script uses aliases, does't work in non-interactive shells unless explicitly enabled
# shitty scripts
# IFS=$'\n\t\v'
cd `dirname "${BASH_SOURCE[0]:-$0}"`

>&2 echo "starting setup"

tmp_dir=$(mktemp -d -t ci-XXXXXXXXXX)
cd $tmp_dir
>&2 echo "setting up daq application in $tmp_dir"

>&2 echo "sourcing daq built tools"
source /opt/dune/daq-buildtools/dbt-setup-env.sh

>&2 echo "setting up project"
dbt-create.sh dunedaq-v2.2.0
dbt-setup-build-environment

>&2 echo "checking out latest code"
cd sourcecode
>&2 echo "-> restcmd"
git clone https://github.com/DUNE-DAQ/restcmd.git
>&2 echo "-> cmdlib"
git clone https://github.com/DUNE-DAQ/cmdlib.git
>&2 echo "-> appfwk"
git clone https://github.com/DUNE-DAQ/appfwk.git

>&2 echo "building code"
dbt-build.sh --clean --install
dbt-setup-runtime-environment

>&2 echo "starting"
# exec daq_application --commandFacility rest://localhost:12345
exec daq_application "$@"
