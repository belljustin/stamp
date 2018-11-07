#!/bin/bash
set -e
parent_path=$(dirname ${BASH_SOURCE[0]})

geth --datadir=$parent_path/../data/ \
    init $parent_path/../genesis.json
