// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package kpis

import (
	"github.com/onosproject/onos-lib-go/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
)

// Var definitions of onos topo metrics builder and static labels.
// builder is used to create metrics in the PrometheusFormat.
var (
	staticLabelsOnosTopo = map[string]string{"sdran": "topo"}
	onosTopoBuilder      = prom.NewBuilder("onos", "topo", staticLabelsOnosTopo)
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

type TopoEntitySlice struct {
	NodeID        string
	Kind          string
	SliceID       string
	SliceDesc     string
	SchedulerType string
	Weight        string
	QosLevel      string
	SliceType     string
	UeIdList      string
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

type topoSlices struct {
	name        string
	description string
	Labels      []string
	LabelValues []string
	Slices      map[string]TopoEntitySlice
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for topoRelations.
func (t *topoRelations) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	t.Labels = []string{"relationid", "kind", "source", "target", "labels", "aspects"}
	metricDesc := onosTopoBuilder.NewMetricDesc(t.name, t.description, t.Labels, staticLabelsOnosTopo)

	for _, relation := range t.Relations {
		metric := onosTopoBuilder.MustNewConstMetric(
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
	metricDesc := onosTopoBuilder.NewMetricDesc(t.name, t.description, t.Labels, staticLabelsOnosTopo)

	for _, entity := range t.Entities {
		metric := onosTopoBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			1.0,
			entity.ID,
			entity.Kind,
			entity.Labels,
			entity.Aspects,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for topoSlices.
func (t *topoSlices) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	t.Labels = []string{"entityid", "kind", "slice_id", "slice_desc", "scheduler_type", "weight", "qoslevel", "slice_type", "ue_id_list"}
	metricDesc := onosTopoBuilder.NewMetricDesc(t.name, t.description, t.Labels, staticLabelsOnosTopo)

	for _, entitySlice := range t.Slices {
		metric := onosTopoBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			1.0,
			entitySlice.NodeID,
			entitySlice.Kind,
			entitySlice.SliceID,
			entitySlice.SliceDesc,
			entitySlice.SchedulerType,
			entitySlice.Weight,
			entitySlice.QosLevel,
			entitySlice.SliceType,
			entitySlice.UeIdList,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
