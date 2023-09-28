# Copyright (c) 2023 Julian Klaiber

GOCMD=go
BINARY_NAME=hawkwing
CLANG ?= clang
CFLAGS :=  -O2 -g -Wall $(CFLAGS) -DDEBUG

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all build clean

all: go-gen build ## Build the entire project

install-deps: ## Install development dependencies
	go install honnef.co/go/tools/cmd/staticcheck@latest
	sudo apt install clang clang-format llvm gcc libbpf-dev libelf-dev make linux-headers-$(uname -r)
	sudo ln -s /usr/include/x86_64-linux-gnu/asm /usr/include/asm
# https://github.com/xdp-project/xdp-tools
# https://github.com/libbpf/bpftool/blob/master/README.md

build: ## Compile the Go binary
	mkdir -p out/bin
	$(GOCMD) build -o out/bin/$(BINARY_NAME) ./$(BINARY_NAME)/main.go

clean: ## Clean build artifacts
	rm -fr out

go-gen: export BPF_CLANG := $(CLANG)
go-gen: export BPF_CFLAGS := $(CFLAGS)
go-gen: ## Generate BPF code and Go bindings
	go generate ./...

test: ## Run go tests
	go clean -testcache
	go test ./...

test-coverage: ## Run go tests with coverage
	go clean -testcache
	go test ./... -coverprofile=coverage.out

grpc-gen: ## Generate gRPC code
	protoc -I ./api/proto --go_out=plugins=grpc:./api/proto ./api/proto/*.proto

setup-network: ## Setup the development network environment
	cd tools && sudo ./network.sh -s

clean-network: ## Clean the development network environment
	cd tools && sudo ./network.sh -c

start-client: ## Start the client
	@echo "Starting client..."
	cd tools && sudo ip netns exec ns-host-a ./network.sh -p

start-server: ## Start the server
	@echo "Starting server..."
	cd tools && sudo ip netns exec ns-host-b ./network.sh -q

start-dns-server: ## Start the dns server in namespace ns-dns
	@echo "Starting dns server..."
	cd tools/dns && sudo ip netns exec ns-dns ./dns

start-webserver_%: ## Start the webserver (usage: make start-webserver_<namespace>-<port>)
	NAMESPACE=`echo $* | cut -d- -f1` && \
	PORT=`echo $* | cut -d- -f2` && \
	echo "Starting webserver in namespace ns-host-$$NAMESPACE on port $$PORT..." && \
	cd tools/web && sudo ip netns exec ns-host-$$NAMESPACE ./webserver host-$$NAMESPACE $$PORT
	# cd tools/web && sudo ip netns exec ns-host-$* ./webserver host-$* 80

start-controller: ## Start the controller
	@echo "Starting controller..."

fix-clang-style: ## Fix the clang style
	find . -iname *.h -o -iname *.c | xargs clang-format -i
	
help: ## Show this help message
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_0-9-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)