contracts = truffle/contracts

all:
	abigen --sol=$(contracts)/StampStorage.sol --pkg=stamper --out=pkg/stamper/stampStorageContract.go
	go build -o server cmd/server.go

migrate:
	$(MAKE) -C truffle all

db:
	docker-compose -f build/docker-compose.yml up -d db
	sleep 20s # TODO: healthcheck
	rambler -c sql/rambler.json apply
