// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-exporter/pkg/config"
	"github.com/onosproject/onos-exporter/pkg/export"
)

const (
	endpoint_address          = ":9861"
	endpoint_path             = "/metrics"
	exporter_mode             = "prometheus"
	e2tEndpointDefault        = "onos-e2t:5150"
	xappPciEndpointDefault    = "onos-pci:5150"
	xappKpimonEndpointDefault = "onos-kpimon:5150"
	topoEndpointDefault       = "onos-topo:5150"
)

var log = logging.GetLogger("main")

var fatalErr error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

func main() {
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	address := flag.String("address", endpoint_address, "Exporter endpoint address:port or just :port")
	path := flag.String("path", endpoint_path, "Exporter endpoint path be used to export kpis")
	mode := flag.String("mode", exporter_mode, "Exporter mode (e.g., prometheus, ...)")
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	e2tEndpoint := flag.String("e2tEndpoint", e2tEndpointDefault, "E2T service endpoint")
	xappPciEndpoint := flag.String("xappPciEndpoint", xappPciEndpointDefault, "XApp PCI service endpoint")
	xappKpimonEndpoint := flag.String("xappKpimonEndpoint", xappKpimonEndpointDefault, "XApp Kpimon service endpoint")
	topoEndpoint := flag.String("topoEndpoint", topoEndpointDefault, "Onos topo service endpoint")

	flag.Parse()

	log.Info("Starting onos-exporter")

	cfgs := map[string]export.CollectorConfig{
		config.ONOSE2T: {
			ServiceAddress: *e2tEndpoint,
		},
		config.ONOSXAPPPCI: {
			ServiceAddress: *xappPciEndpoint,
		},
		config.ONOSXAPPKPIMON: {
			ServiceAddress: *xappKpimonEndpoint,
		},
		config.ONOSTOPO: {
			ServiceAddress: *topoEndpoint,
		},
	}

	cfg := export.Config{
		Address:           *address,
		Path:              *path,
		Mode:              *mode,
		CAPath:            *caPath,
		KeyPath:           *keyPath,
		CertPath:          *certPath,
		CollectorsConfigs: cfgs,
	}

	exporter := export.NewExporter(cfg)

	if err := exporter.Run(); err != nil {
		log.Errorf("onos exporter error")
		fatal(err)
	}
}
