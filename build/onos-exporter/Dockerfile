# SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

ARG ONOS_BUILD_VERSION=undefined

FROM onosproject/golang-build:v1.0 as build

ENV GO111MODULE=on
ARG ONOS_MAKE_TARGET=build

COPY Makefile go.mod go.sum /go/src/github.com/onosproject/onos-exporter/
COPY cmd/ /go/src/github.com/onosproject/onos-exporter/cmd/
COPY pkg/ /go/src/github.com/onosproject/onos-exporter/pkg/
COPY vendor/ /go/src/github.com/onosproject/onos-exporter/vendor/
COPY api/ /go/src/github.com/onosproject/onos-exporter/api
COPY build/ /go/src/github.com/onosproject/onos-exporter/build/


RUN cd /go/src/github.com/onosproject/onos-exporter && GOFLAGS=-mod=vendor make ${ONOS_MAKE_TARGET}

FROM alpine:3.20
RUN apk add bash bash-completion libc6-compat

RUN addgroup -S onos && adduser -S -G onos onos --home /home/onos
USER onos
WORKDIR /home/onos
RUN mkdir /home/onos/.onos

COPY --from=build /go/src/github.com/onosproject/onos-exporter/build/_output/onos-exporter /usr/local/bin/onos-exporter

ENTRYPOINT ["onos-exporter"]
