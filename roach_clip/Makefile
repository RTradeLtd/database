.PHONY: install
install:
	./scripts/install.sh

# used to start a single node for easy testing
# to start in insecure mode:
# env FLAGS=--insecure make single-node
.PHONY: single-node
single-node:
	cockroach start-single-node $(FLAGS)