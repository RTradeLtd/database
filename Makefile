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
	(cd testenv/roach_clip ; make start-testenv)

.PHONY: clean
clean:
	(cd testenv/roach_clip ; make stop-test-node)

.PHONY: test
test: vendor
	go test -race -cover ./...

.PHONY: lint
lint: vendor
	golint $(go list ./... | grep -v /vendor/)
