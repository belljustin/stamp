# Stamp

The [StampStorage contract](truffle/contracts/StampStorage.sol) is a simple ethereum contract whose purpose is to record proof of document existence at a given timestamp.

The SHA256 hash of a document is submitted to the contract which will be "stamped" with a datetime represented as seconds since the epoch.
Anyone can verify the existence and check the timestamp of these stamps on the public blockchain.

## Stamp API

### Submitting stamps
To avoid broadcasting frequent and numerous document stamps - which may be prohibitively expensive - document hashes can be combined into a [Merkle Tree](https://en.wikipedia.org/wiki/Merkle_tree) and the tree root submitted for stamping.
Document hashes submitted to the Stamp API are periodically batched and submitted in this way.

### Document proofs
A subset of the tree nodes, which are sufficient to prove a valid document was included in the stamp, may be retrieved after the stamp has been broadcasted.
This proof size is proportianal to log(n), where n is the number of documents sumbitted in the stamp.
Others can independently and trustlessly:
- verify this proof forms a valid merkle tree,
- its root is stored in the StampStorage contract,
- and check the stamp timestamp, t
proving the document existed _at least since_ time t.

## Install

### Run `make go` to:
- generate the go contract code from the solidity contract
- build the stamp server binary

### Run `make truffle` to:
- compile all the solidity contracts
- perform ethereum migrations on geth network
_Note_: this will prompt for the password of the ethereum account which will become the contract

### Run `make db` to:
- start the postgres docker container
- apply all migrations against the database

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
