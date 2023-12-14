SHELL := /bin/bash
.PHONY: help deps build

LD_FLAGS=-ldflags="-s -w"

# Ref: https://gist.github.com/prwhite/8168133
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

deps:  ## Update vendor.
	go mod verify
	go mod tidy -v

build: ## build tsdb tester binary
	go build -o 'bin/tsdb-tester' $(LD_FLAGS) ./cmd

run: ## run demo
	go run github.com/lindb/tsdb-tester/cmd
