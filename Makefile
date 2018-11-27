contracts = truffle/contracts

go:
	abigen --sol=$(contracts)/StampStorage.sol --pkg=stamper --out=pkg/stamper/stampStorageContract.go
	go build -o build/stamp cmd/stamp.go

db:
	rambler -c sql/rambler.json apply --all

truffle:
	$(MAKE) -C truffle all
