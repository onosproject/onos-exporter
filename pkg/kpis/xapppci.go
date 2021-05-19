// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpis

import (
	"github.com/onosproject/onos-lib-go/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
)

// Var definitions of xapp pci metrics builder and static labels.
// builder is used to create metrics in the PrometheusFormat.
var (
	staticLabelsXappPci = map[string]string{"sdran": "xapppci"}
	xappPciBuilder      = prom.NewBuilder("onos", "xapppci", staticLabelsXappPci)
)

// xapppciNumConflicts defines the common data that can be used
// to output the format of a KPI (e.g., PrometheusFormat).
// NumberConflicts stores the number of conflicts per cell id.
type xappPciNumConflicts struct {
	name            string
	description     string
	Labels          []string
	LabelValues     []string
	NumberConflicts map[string]float64
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for xapppciNumConflicts.
func (c *xappPciNumConflicts) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	c.Labels = []string{"cellid"}
	metricDesc := xappPciBuilder.NewMetricDesc(c.name, c.description, c.Labels, staticLabelsXappPci)

	for cellID, numConflicts := range c.NumberConflicts {
		metric := xappPciBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			numConflicts,
			cellID,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
