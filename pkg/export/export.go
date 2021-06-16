// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package export

import "github.com/onosproject/onos-exporter/pkg/collect"

// Consts define the names of the available collector names.
// It defines those names based on the collectors available
// at the collect package.
const (
	ONOSE2T        = collect.ONOSE2T
	ONOSXAPPKPIMON = collect.ONOSXAPPKPIMON
	ONOSXAPPPCI    = collect.ONOSXAPPPCI
)

// CollectorConfig states the parameters that enables a Collector.
type CollectorConfig struct {
	ServiceAddress string
	CAPath         string
	KeyPath        string
	CertPath       string
}

// Config establishes the fields needed for the instantiation of
// an exporter.
// Address and Path define the exporter endpoint from where KPIs can
// be pulled or pushed.
// Mode defines the exporter mode, i.e., the exporter implementation mode,
// for instance, prometheus.
// CAPath, KeyPath and CertPath are defined by the utilization of
// a northbound implementation of needed certificates for an exporter.
// The remaining fields define the needed data needed for the exporters,
// those fields can be defined in their own structs if needed.
type Config struct {
	Address           string
	Path              string
	Mode              string
	CAPath            string
	KeyPath           string
	CertPath          string
	CollectorsConfigs map[string]CollectorConfig
}

// exporter defines the behavior expected from an exporter.
type exporter interface {
	Run() error
}

// NewExporter defines a factory for an exporter interface.
// PrometheusExporter realizes that interface behavior.
// Other exporters can be added similarly. Turning the implementation
// of onos-exporter independent from a single exporter.
func NewExporter(cfg Config) exporter {
	switch cfg.Mode {
	case "prometheus":
		log.Info("Creating prometheus exporter")
		return PrometheusExporter(cfg)
	default:
		log.Info("Creating default exporter (prometheus)")
		return PrometheusExporter(cfg)
	}
}
