# Stamp

## Install

Run `make` to generate the go contract code from the solidity contract and build the binary.

## Configuration

### Geth

To configure private test network with geth:

1. `./truffle/scripts/initGeth.sh` to initialize with the the genesis block defined at `./truffle/genesis.json`
    - If you have a keyfile for an existing address, place it in `truffle/data/keystore/`
2. Run `./truffle/scripts/startGeth.sh` to start the geth node
3. In a new terminal, run `make migrate` to compile and deploy the contract to the geth node
    - You'll have to provide the password to the keyfile

### Environment Variables

- STAMP\_PRIVATE\_KEY: [hex value of the private key used for stamping]

### Files

- config.json
- truffle/truffle-config.js
- (dev|prod).secret.env

## Requirements
- geth
- solc
- truffle
- [rambler](https://github.com/elwinar/rambler/releases/download/4.2.0/rambler-darwin-10.6-386)
