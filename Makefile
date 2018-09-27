TESTCONFIG=https://raw.githubusercontent.com/RTradeLtd/Temporal/V2/test/config.json
TESTCOMPOSE=https://raw.githubusercontent.com/RTradeLtd/Temporal/V2/test/docker-compose.yml

COMPOSECOMMAND=env ADDR_NODE1=1 ADDR_NODE2=2 docker-compose -f test/docker-compose.yml

all: build

.PHONY: build
build:
	go build ./...

.PHONY: testenv
testenv:
	mkdir -p test
	curl $(TESTCONFIG) --output test/config.json
	curl $(TESTCOMPOSE) --output test/docker-compose.yml
	$(COMPOSECOMMAND) up -d postgres

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: lint
lint:
	golint $(go list ./... | grep -v /vendor/)
