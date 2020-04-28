.PHONY: install
install:
	./scripts/install.sh

# used to start a single node for easy testing
# to start in insecure mode:
# env FLAGS=--insecure make single-node
.PHONY: single-node
single-node:
	cockroach start-single-node $(FLAGS)

# used to generate the required certificates for our test node
.PHONY: gen-test-certs
gen-test-certs:
	rm -rf certs safe-dir cockroach-data
	mkdir certs && mkdir safe-dir
	./scripts/gen_certs.sh ca ./certs ./safe-dir/ca.key
	./scripts/gen_certs.sh node ./certs ./safe-dir/ca.key localhost
	./scripts/gen_certs.sh client ./certs ./safe-dir/ca.key temporal
	./scripts/gen_certs.sh client ./certs ./safe-dir/ca.key root

# used to start our secure single-node test server
.PHONY: start-test-node
start-test-node:
	cockroach \
		start-single-node \
		--certs-dir=certs \
		--listen-addr=localhost:26257 \
		--http-addr=localhost:8080 \
		--background

# setups the database after it starts
.PHONY: setup-test-node
setup-test-node:
	cockroach \
		sql \
		--certs-dir=certs \
		--user=root \
		--host=localhost:26257 \
		--execute="CREATE USER temporaladmin WITH PASSWORD 'temporaladmin'"