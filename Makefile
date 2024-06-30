# SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build

ONOS_EXPORTER_VERSION ?= latest

GOLANG_CI_VERSION := v1.52.2

all: build docker-build

build: # @HELP build the Go binaries and run all validations (default)
	GOPRIVATE="github.com/onosproject/*" go build -o build/_output/onos-exporter ./cmd/onos-exporter

test: # @HELP run the unit tests and source code validation
test: build lint license
	go test -race github.com/onosproject/onos-exporter/pkg/...
	go test -race github.com/onosproject/onos-exporter/cmd/...


docker-build-onos-exporter: # @HELP build onos-exporter Docker image
	@go mod vendor
	docker build . -f build/onos-exporter/Dockerfile \
		-t onosproject/onos-exporter:${ONOS_EXPORTER_VERSION}
	@rm -rf vendor

images: # @HELP build all Docker images
docker-build: build docker-build-onos-exporter

docker-push-onos-exporter: # @HELP push onos-exporter Docker image
	docker push onosproject/onos-exporter:${ONOS_EXPORTER_VERSION}

docker-push: # @HELP push docker images
docker-push: docker-push-onos-exporter

lint: # @HELP examines Go source code and reports coding problems
	golangci-lint --version | grep $(GOLANG_CI_VERSION) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b `go env GOPATH`/bin $(GOLANG_CI_VERSION)
	golangci-lint run --timeout 15m

license: # @HELP run license checks
	rm -rf venv
	python3 -m venv venv
	. ./venv/bin/activate;\
	python3 -m pip install --upgrade pip;\
	python3 -m pip install reuse;\
	reuse lint

check-version: # @HELP check version is duplicated
	./build/bin/version_check.sh all

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos-pci/onos-exporter ./cmd/onos/onos venv
	go clean github.com/onosproject/onos-exporter/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '