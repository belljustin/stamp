## Configuration

### Geth

To configure private test network with geth:

1. `./scripts/initGeth.sh` to initialize with the the genesis block defined at `genesis.json`
    - If you have a keyfile for an existing address, place it in `truffle/data/keystore/`
2. Run `./truffle/scripts/startGeth.sh` to start the geth node
3. In a new terminal, run `make migrate` to compile and deploy the contract to the geth node
    - You'll have to provide the password to the keyfile
