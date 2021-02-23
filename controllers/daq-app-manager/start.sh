#!/usr/bin/env bash
cd `dirname "${BASH_SOURCE[0]:-$0}"`
shopt -s expand_aliases # script uses aliases, does't work in non-interactive shells unless explicitly enabled

# venv's aren't portable across machines, rebuild it
rm -rf venv
make build

# activate venv
source venv/bin/activate

exec python3 main.py