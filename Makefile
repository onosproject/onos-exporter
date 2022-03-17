# SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build

ONOS_EXPORTER_VERSION := latest

build: # @HELP build the Go binaries and run all validations (default)
build:
	GOPRIVATE="github.com/onosproject/*" go build -o build/_output/onos-exporter ./cmd/onos-exporter

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

test: # @HELP run the unit tests and source code validation
test: build deps linters license
	go test -race github.com/onosproject/onos-exporter/pkg/...
	go test -race github.com/onosproject/onos-exporter/cmd/...

jenkins-test:  # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: deps license linters
	TEST_PACKAGES=github.com/onosproject/onos-exporter/... ./build/build-tools/build/jenkins/make-unit

buflint: #@HELP run the "buf check lint" command on the proto files in 'api'
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-exporter \
		-w /go/src/github.com/onosproject/onos-exporter/api \
		bufbuild/buf:${BUF_VERSION} check lint

onos-exporter-docker: # @HELP build onos-exporter Docker image
onos-exporter-docker:
	@go mod vendor
	docker build . -f build/onos-exporter/Dockerfile \
		-t onosproject/onos-exporter:${ONOS_EXPORTER_VERSION}
	@rm -rf vendor
	
images: # @HELP build all Docker images
images: build onos-exporter-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-exporter:${ONOS_EXPORTER_VERSION}

all: build images

publish: # @HELP publish version on github and dockerhub
	./build/build-tools/publish-version ${VERSION} onosproject/onos-exporter

jenkins-publish: # @HELP Jenkins calls this to publish artifacts
	./build/bin/push-images
	./build/build-tools/release-merge-commit

clean:: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos-exporter/onos-exporter ./cmd/onos/onos
	go clean -testcache github.com/onosproject/onos-exporter/...

