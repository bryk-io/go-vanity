.DEFAULT_GOAL := help
.PHONY: all
VERSION_TAG=0.1.0
BINARY_NAME=govanity
DOCKER_IMAGE=govanity

# Linker tags
# https://golang.org/cmd/link/
LD_FLAGS="\
-X 'main.coreVersion=$(VERSION_TAG)' \
-X 'main.buildTimestamp=`date +'%s'`' \
-X 'main.buildCode=`git log --pretty=format:'%H' -n1`' \
"

help: ## Display available make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-16s\033[0m %s\n", $$1, $$2}'

ca-roots: ## Generate the list of valid CA certificates
	@docker run -dit --rm --name ca-roots debian:stable-slim
	@docker exec --privileged ca-roots sh -c "apt update"
	@docker exec --privileged ca-roots sh -c "apt install -y ca-certificates"
	@docker exec --privileged ca-roots sh -c "cat /etc/ssl/certs/* > /ca-roots.crt"
	@docker cp ca-roots:/ca-roots.crt ca-roots.crt
	@docker stop ca-roots

test: ## Run all tests excluding the vendor dependencies
	# Static analysis
	golangci-lint run -v ./...
	go-consistent -v ./...

	# Unit tests
	go test -race -cover -v ./...

clean: ## Download and compile all dependencies and intermediary products
	@-rm -rf vendor
	go mod tidy
	go mod verify
	go mod vendor

updates: ## List available updates for direct dependencies
	# https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies
	go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null

install: ## Install the binary to GOPATH and keep cached all compiled artifacts
	@go build -v -ldflags $(LD_FLAGS) -i -o ${GOPATH}/bin/$(BINARY_NAME)

build: ## Build for the current architecture in use, intended for devevelopment
	go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)

build-for: ## Build the availabe binaries for the specified 'os' and 'arch'
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) \
	go build -v -ldflags $(LD_FLAGS) \
	-o $(dest)$(BINARY_NAME)_$(VERSION_TAG)_$(os)_$(arch)$(suffix)

docker: ## Build docker image
	@-rm $(BINARY_NAME)_$(VERSION_TAG)_linux_amd64 ca-roots.crt
	@make ca-roots
	@make build-for os=linux arch=amd64
	@-docker rmi $(DOCKER_IMAGE):$(VERSION_TAG)
	@docker build --build-arg VERSION_TAG="$(VERSION_TAG)" --rm -t $(DOCKER_IMAGE):$(VERSION_TAG) .
	@-rm $(BINARY_NAME)_$(VERSION_TAG)_linux_amd64 ca-roots.crt
