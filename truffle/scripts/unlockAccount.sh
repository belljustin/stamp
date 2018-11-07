#!/bin/bash
set -e
parent_path=$(dirname ${BASH_SOURCE[0]})

echo -n Password: 
read -s password

geth attach $parent_path/../data/geth.ipc << EOF
personal.unlockAccount(personal.listAccounts[0])
$password
miner.start()
