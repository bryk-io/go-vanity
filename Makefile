.DEFAULT_GOAL := help
.PHONY: all
VERSION_TAG=0.1.2
BINARY_NAME=govanity
DOCKER_IMAGE_NAME=docker.pkg.github.com/bryk-io/go-vanity/govanity

# Linker tags
# https://golang.org/cmd/link/
LD_FLAGS += -s -w
LD_FLAGS += -X main.coreVersion=$(VERSION_TAG)
LD_FLAGS += -X main.buildTimestamp=$(shell date +'%s')
LD_FLAGS += -X main.buildCode=$(shell git log --pretty=format:'%H' -n1)

## help: Prints this help message
help:
	@echo "Commands available"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /' | sort

## clean: Download and compile all dependencies and intermediary products
clean:
	@-rm -rf vendor
	go mod tidy
	go mod verify
	go mod download
	go mod vendor

## updates: List available updates for direct dependencies
# https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies
updates:
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null

## lint: Static analysis
lint:
	helm lint helm/*
	golangci-lint run -v ./...

## scan: Look for known vulnerabilities in the project dependencies
# https://github.com/sonatype-nexus-community/nancy
scan:
	@nancy -quiet go.sum

## test: Run unit tests excluding the vendor dependencies
test:
	go test -v -race -failfast -coverprofile=coverage.report ./...
	go tool cover -html coverage.report -o coverage.html

## ca-roots: Generate the list of valid CA certificates
ca-roots:
	@docker run -dit --rm --name ca-roots debian:stable-slim
	@docker exec --privileged ca-roots sh -c "apt update"
	@docker exec --privileged ca-roots sh -c "apt install -y ca-certificates"
	@docker exec --privileged ca-roots sh -c "cat /etc/ssl/certs/* > /ca-roots.crt"
	@docker cp ca-roots:/ca-roots.crt ca-roots.crt
	@docker stop ca-roots

## build: Build for the current architecture in use, intended for development
build:
	go build -v -ldflags '$(LD_FLAGS)' -o $(BINARY_NAME)

## build-for: Build the available binaries for the specified 'os' and 'arch'
build-for:
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) \
	go build -v -ldflags '$(LD_FLAGS)' \
	-o $(dest)$(BINARY_NAME)_$(VERSION_TAG)_$(os)_$(arch)$(suffix)

## install: Install the binary to GOPATH and keep cached all compiled artifacts
install:
	@go build -v -ldflags '$(LD_FLAGS)' -i -o ${GOPATH}/bin/$(BINARY_NAME)

## docker: Build docker image
docker:
	make build-for os=linux arch=amd64
	@-docker rmi $(DOCKER_IMAGE_NAME):$(VERSION_TAG)
	@docker build --build-arg VERSION="$(VERSION_TAG)" --rm -t $(DOCKER_IMAGE_NAME):$(VERSION_TAG) .
	rm $(BINARY_NAME)_$(VERSION_TAG)_linux_amd64

## docs: Display package documentation on local server
docs:
	@echo "Docs available at: http://localhost:8080/pkg/github.com/bryk-io/go-vanity/"
	godoc -http=:8080

## release: Prepare artifacts for a new tagged release
release:
	@-rm -rf release-$(VERSION_TAG)
	mkdir release-$(VERSION_TAG)
	make build-for os=linux arch=amd64 dest=release-$(VERSION_TAG)/
	make build-for os=darwin arch=amd64 dest=release-$(VERSION_TAG)/
	make build-for os=windows arch=amd64 suffix=".exe" dest=release-$(VERSION_TAG)/
	make build-for os=windows arch=386 suffix=".exe" dest=release-$(VERSION_TAG)/
