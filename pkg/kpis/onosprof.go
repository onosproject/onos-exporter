// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package kpis

import (
	"github.com/onosproject/onos-lib-go/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
)

// Var definitions of e2t metrics onose2tBuilder and static labels.
// builder is used to create metrics in the PrometheusFormat.
var (
	staticLabelsProf   = map[string]string{"sdran": "profile"}
	onosProfileBuilder = prom.NewBuilder("onos", "profile", staticLabelsProf)
)

type HeapObject struct {
	Value  int64
	Name   string
	Source string
	Format string
}

type onosProfileHeap struct {
	name        string
	description string
	Labels      []string
	LabelValues []string
	Objects     map[string]HeapObject
}

func (c *onosProfileHeap) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	c.Labels = []string{"name", "source", "format"}
	metricDesc := onosProfileBuilder.NewMetricDesc(c.name, c.description, c.Labels, staticLabelsProf)

	for _, obj := range c.Objects {
		metric := onosProfileBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			float64(obj.Value),
			obj.Name,
			obj.Source,
			obj.Format,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
