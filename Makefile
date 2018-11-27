contracts = truffle/contracts

go:
	abigen --sol=$(contracts)/StampStorage.sol --pkg=stamper --out=pkg/stamper/stampStorageContract.go
	go build -o build/stamp cmd/stamp.go

db:
	docker-compose -f build/docker-compose.yml up -d db
	sleep 20s # TODO: healthcheck
	rambler -c sql/rambler.json apply --all

truffle:
	$(MAKE) -C truffle all
