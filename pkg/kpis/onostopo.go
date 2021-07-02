// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpis

import (
	"github.com/onosproject/onos-lib-go/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
)

// Var definitions of onos topo metrics builder and static labels.
// builder is used to create metrics in the PrometheusFormat.
var (
	staticLabelsOnosTopo = map[string]string{"sdran": "topo"}
	xappTopoBuilder      = prom.NewBuilder("onos", "topo", staticLabelsOnosTopo)
)

type TopoRelation struct {
	ID      string
	Kind    string
	Source  string
	Target  string
	Labels  string
	Aspects string
}

type TopoEntity struct {
	ID      string
	Kind    string
	Labels  string
	Aspects string
}

type topoRelations struct {
	name        string
	description string
	Labels      []string
	LabelValues []string
	Relations   map[string]TopoRelation
}

type topoEntities struct {
	name        string
	description string
	Labels      []string
	LabelValues []string
	Entities    map[string]TopoEntity
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for topoRelations.
func (t *topoRelations) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	t.Labels = []string{"relationid", "kind", "source", "target", "labels", "aspects"}
	metricDesc := xappTopoBuilder.NewMetricDesc(t.name, t.description, t.Labels, staticLabelsOnosTopo)

	for _, relation := range t.Relations {
		metric := xappTopoBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			1.0,
			relation.ID,
			relation.Kind,
			relation.Source,
			relation.Target,
			relation.Labels,
			relation.Aspects,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for topoEntities.
func (t *topoEntities) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	t.Labels = []string{"entityid", "kind", "labels", "aspects"}
	metricDesc := xappTopoBuilder.NewMetricDesc(t.name, t.description, t.Labels, staticLabelsOnosTopo)

	for _, relation := range t.Entities {
		metric := xappTopoBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			1.0,
			relation.ID,
			relation.Kind,
			relation.Labels,
			relation.Aspects,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
