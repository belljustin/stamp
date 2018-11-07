#!/bin/bash
set -e
parent_path=$(dirname ${BASH_SOURCE[0]})

geth --identity "MyTestNetNode" \
    --nodiscover \
    --networkid 1999 \
    --datadir $parent_path/../data \
    --ws \
    --rpc
