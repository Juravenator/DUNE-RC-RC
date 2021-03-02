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
git clone https://github.com/DUNE-DAQ/daq-buildtools.git -b v2.1.1
source daq-buildtools/dbt-setup-env.sh

>&2 echo "setting up project"
mkdir listrev-app
cd listrev-app
dbt-create.sh dunedaq-v2.2.0

>&2 echo "checking out newer code"
cd sourcecode
git clone https://github.com/DUNE-DAQ/daq-cmake.git -b v1.3.1
git clone https://github.com/DUNE-DAQ/ers.git -b v1.1.0
git clone https://github.com/DUNE-DAQ/logging.git -b v1.0.1
git clone https://github.com/DUNE-DAQ/cmdlib.git -b v1.1.1
git clone https://github.com/DUNE-DAQ/restcmd.git -b v1.1.0 
git clone https://github.com/DUNE-DAQ/rcif.git -b v1.0.1
git clone https://github.com/DUNE-DAQ/opmonlib.git -b v1.0.0
git clone https://github.com/DUNE-DAQ/appfwk.git -b v2.2.0
git clone https://github.com/DUNE-DAQ/listrev.git -b v2.1.1
echo 'set(build_order "daq-cmake" "ers" "logging" "cmdlib" "rcif" "restcmd" "opmonlib" "appfwk" "listrev")' > dbt-build-order.cmake
cd ..

>&2 echo "building code"
dbt-setup-build-environment
pip install -U https://github.com/brettviren/moo/archive/0.5.5.tar.gz
dbt-build.sh --clean --install
dbt-setup-runtime-environment

>&2 echo "starting daq_application $@"
set -o errexit
daq_application "$@"