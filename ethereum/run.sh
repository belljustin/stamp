#!/bin/sh
set -e

BASE=$(dirname "$0")

# initialize the node
geth --datadir=$BASE/data \
    init $BASE/genesis.json

geth --identity "MyTestNetNode" \
    --nodiscover \
    --networkid 1999 \
    --datadir "$BASE/data" \
    --rpc \
    --rpcaddr 0.0.0.0 \
    --ws \
    --wsaddr 0.0.0.0 \
    --wsorigins="*" \
    --mine \
    -minerthreads=1 \
    --etherbase=0x8b764be433cee171951328586e38ef2450a62def
