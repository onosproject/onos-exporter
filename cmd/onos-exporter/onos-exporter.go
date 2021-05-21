// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-exporter/pkg/export"
)

const (
	endpoint_address   = ":9861"
	endpoint_path      = "/metrics"
	exporter_mode      = "prometheus"
	e2tEndpoint        = "onos-e2t:5150"
	e2subEndpoint      = "onos-e2sub:5150"
	xappPciEndpoint    = "onos-pci:5150"
	xappKpimonEndpoint = "onos-kpimon-v2:5150"
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
	e2tEndpoint := flag.String("e2tEndpoint", e2tEndpoint, "E2T service endpoint")
	e2subEndpoint := flag.String("e2subEndpoint", e2subEndpoint, "E2Sub service endpoint")
	xappPciEndpoint := flag.String("xappPciEndpoint", xappPciEndpoint, "XApp PCI service endpoint")
	xappKpimonEndpoint := flag.String("xappKpimonEndpoint", xappKpimonEndpoint, "XApp Kpimon service endpoint")

	flag.Parse()

	log.Info("Starting onos-exporter")

	cfgs := map[string]export.CollectorConfig{
		export.ONOSE2T: {
			ServiceAddress: *e2tEndpoint,
		},
		export.ONOSE2SUB: {
			ServiceAddress: *e2subEndpoint,
		},
		export.ONOSXAPPPCI: {
			ServiceAddress: *xappPciEndpoint,
		},
		export.ONOSXAPPKPIMON: {
			ServiceAddress: *xappKpimonEndpoint,
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
