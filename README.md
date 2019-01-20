# Stamp

The [StampStorage contract](truffle/contracts/StampStorage.sol) is a simple ethereum contract whose purpose is to record proof of document existence at a given timestamp.

The SHA256 hash of a document is submitted to the contract which will be "stamped" with a datetime represented as seconds since the epoch.
Anyone can verify the existence and check the timestamp of these stamps on the public blockchain.

## Stamp API

### Submitting stamps
To avoid broadcasting frequent and numerous document stamps - which may be prohibitively expensive - document hashes can be combined into a [Merkle Tree](https://en.wikipedia.org/wiki/Merkle_tree) and the tree root submitted for stamping.
Document hashes submitted to the Stamp API are periodically batched and submitted in this way.

### Document proofs
A subset of the tree nodes, which are sufficient to prove a valid document is included in a stamp, may be retrieved after the stamp has been mined.
The proof's size is proportianal to log(n), where n is the number of documents submitted in the stamp.
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

### Database
Database migrations are run with rambler.
To apply the next database migration, run `rambler -c sql/rambler.json apply`.
To revert to the previous migration, run `rambler -c sql/rambler.json reverse`

### Run
```sh
docker-compose up
./build/stamp -config config.json
```

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

# API

- [Ping](api.md#ping) `GET /ping`
- [Get Document](api.md#get-document) `GET /document/{id}`
- [Create Document](api.md#create-document) `POST /document`
- [Get Stamp](api.md#get-stamp) `GET /stamp/{id}`
