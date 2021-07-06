// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpis

import (
	"github.com/onosproject/onos-lib-go/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
)

// Var definitions of onos uenib metrics builder and static labels.
// builder is used to create metrics in the PrometheusFormat.
var (
	staticLabelsOnosUenib = map[string]string{"sdran": "uenib"}
	onosUenibBuilder      = prom.NewBuilder("onos", "uenib", staticLabelsOnosUenib)
)

type UE struct {
	ID            string
	Aspects       []string
	AspectsValues []string
	Relations     map[string]TopoRelation
}

type onosUenibUEs struct {
	name        string
	description string
	Labels      []string
	LabelValues []string
	UEs         map[string]UE
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for onosUenibUEs.
func (t *onosUenibUEs) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	for _, ue := range t.UEs {
		t.Labels = []string{"ueid"}
		t.Labels = append(t.Labels, ue.Aspects...)

		t.LabelValues = []string{ue.ID}
		t.LabelValues = append(t.LabelValues, ue.AspectsValues...)

		metricDesc := onosUenibBuilder.NewMetricDesc(t.name, t.description, t.Labels, staticLabelsOnosUenib)

		metric := onosUenibBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			1.0,
			t.LabelValues...,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
