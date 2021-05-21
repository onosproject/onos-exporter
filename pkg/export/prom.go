// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package export

import (
	"github.com/onosproject/onos-exporter/pkg/collect"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	log = logging.GetLogger("export", "prom")

	// All the names of collectors that prometheus exporter instantiates.
	collectorNames = []string{
		ONOSE2T,
		ONOSE2SUB,
		ONOSXAPPKPIMON,
		ONOSXAPPPCI,
	}
)

// CollectorsPrometheus defines a prometheus collector
// for all collectors.
type CollectorsPrometheus struct {
	collectors []collect.Collector
}

// Retrieve implements the method needed for a Collector interface
// in a prometheus exporter. It retrieves all the kpis from
// CollectorsPrometheus and pass them to the ch channel using the
// prometheus.Metric format.
// The function collect.KPIs performs the collection of each collector
// list of KPIs, and aggregates them in onosKPIs var.
func (c *CollectorsPrometheus) Retrieve(ch chan<- prometheus.Metric) error {
	onosKPIs := collect.KPIs(c.collectors)

	for _, kpi := range onosKPIs {
		promMetrics, err := kpi.PrometheusFormat()

		if err != nil {
			log.Errorf("onos kpi prometheus format error %s", err)
		} else {
			for _, m := range promMetrics {
				ch <- m
			}
		}
	}

	return nil
}

// Defines the set of collector used to extract KPIs for
// the prometheus exporter. Each collector implements the
// prom.Collector interface behavior via the method Collect.
func initCollectorsPrometheus(config Config) prom.Collector {
	collectors := []collect.Collector{}

	for _, collectorName := range collectorNames {
		collectorConfig, ok := config.CollectorsConfigs[collectorName]

		if ok {
			collector, err := collect.CreateCollector(collectorName, collectorConfig.ServiceAddress)

			if err != nil {
				log.Errorf("%s not added to collectors %s", collectorName, err)
			} else {
				collectors = append(collectors, collector)
			}

		} else {
			log.Errorf("%s not added to collectors no configuration provided", collectorName)
		}
	}

	return &CollectorsPrometheus{
		collectors: collectors,
	}
}

// PrometheusExporter uses Config to create an instance of a
// Prometheus exporter, registering all its collectors, which must
// implement the interface method Retrieve.
func PrometheusExporter(config Config) prom.Exporter {
	exporter := prom.NewExporter(config.Path, config.Address)

	log.Info("Registering collector sdran")
	err := exporter.RegisterCollector("sdran", initCollectorsPrometheus(config))
	if err != nil {
		log.Errorf("error registering collector sdran %s", err)
	}

	return exporter
}
