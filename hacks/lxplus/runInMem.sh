#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
IFS=$'\n\t\v'
cd `dirname "${BASH_SOURCE[0]:-$0}"`

WORKDIR=/tmp/dune-rc-lxplus

>&2 echo "WARNING"
>&2 echo "This is NOT a production setup"
sleep 2s

mkdir -p $WORKDIR

>&2 echo "checking if consul is installed"
if [[ ! -f $WORKDIR/bin/consul ]]; then
  >&2 echo "downloading consul"
  curl -Lo $WORKDIR/consul.zip https://releases.hashicorp.com/consul/1.9.3/consul_1.9.3_linux_amd64.zip
  >&2 echo "installing consul"
  unzip $WORKDIR/consul.zip -d $WORKDIR/bin
fi
>&2 echo "consul is installed"

>&2 echo "checking if nomad is installed"
if [[ ! -f $WORKDIR/bin/nomad ]]; then
  >&2 echo "downloading nomad"
  curl -Lo $WORKDIR/nomad.zip https://releases.hashicorp.com/nomad/1.0.3/nomad_1.0.3_linux_amd64.zip
  >&2 echo "installing nomad"
  unzip $WORKDIR/nomad.zip -d $WORKDIR/bin
fi
>&2 echo "nomad is installed"

>&2 echo "checking if consul is running"
if ! pgrep consul > /dev/null; then
  >&2 echo "starting consul"
  nohup $WORKDIR/bin/consul agent -dev -client=0.0.0.0 -datacenter=dune-rc </dev/null >/dev/null 2>&1 &
fi
>&2 echo "consul is running"

cp ../../ansible/single-host/nomad/config.hcl $WORKDIR/nomad-config.hcl

>&2 echo "checking if nomad is running"
if ! pgrep nomad > /dev/null; then
  >&2 echo "starting nomad"
  nohup $WORKDIR/bin/nomad agent -dev -bind=0.0.0.0 -config=$WORKDIR/nomad-config.hcl </dev/null >/dev/null 2>&1 &
  >&2 echo "waiting 10s"
  sleep 10s
fi
>&2 echo "nomad is running"

>&2 echo "installing daq-application manager"
rsync -a --exclude=venv/ ../../controllers $WORKDIR

>&2 echo "(re)starting daq-application manager"
! $WORKDIR/bin/nomad stop daq-application-manager 2>/dev/null
$WORKDIR/bin/nomad job run nomad-lxplus.hcl

>&2 echo "installing hack scripts"
mkdir -p /tmp/dune-rc-hacks
rsync -a ../ /tmp/dune-rc-hacks/

cd ..
>&2 echo ""
>&2 echo "the CLI is available at $(pwd)/../cli/build/run-control"
>&2 echo "for convenience, you could run:"
>&2 echo "export PATH=\"$(pwd)/../cli/build:\$PATH\""