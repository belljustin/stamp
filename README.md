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

An example configuration file can be found at `config.json`.

### contract
- address: the hex address of the contract
- host: geth node hostname
- port: geth node port
- interval: period, in seconds, to submit a new stamp
- privateKey: path to the Keystore File (UTC / JSON)
- password: used to decrypt the keystore file

## Requirements
- geth
- solc
- truffle
- [rambler](https://github.com/elwinar/rambler/releases/download/4.2.0/rambler-darwin-10.6-386)
