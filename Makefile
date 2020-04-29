COMPOSECOMMAND=env ADDR_NODE1=1 ADDR_NODE2=2 docker-compose -f testenv/docker-compose.yml

all: build

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

.PHONY: build
build: vendor
	go build ./...

.PHONY: testenv
testenv:
	(cd roach_clip ; make start-testenv)
#	$(COMPOSECOMMAND) up -d postgres

.PHONY: clean
clean:
	(cd roach_clip ; make stop-test-node)
#	$(COMPOSECOMMAND) down

.PHONY: test
test: vendor
	go test -race -cover ./...

.PHONY: lint
lint: vendor
	golint $(go list ./... | grep -v /vendor/)
