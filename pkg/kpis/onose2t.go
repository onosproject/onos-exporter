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
	staticLabelsE2t = map[string]string{"sdran": "e2t"}
	onose2tBuilder  = prom.NewBuilder("onos", "e2t", staticLabelsE2t)
)

type E2tSubscription struct {
	Id                  string
	Revision            string
	ServiceModelName    string
	ServiceModelVersion string
	E2NodeID            string
	Encoding            string
	StatusPhase         string
	StatusState         string
}

// onosE2tSubscriptions defines the common data that can be used
// to output the format of a KPI (e.g., PrometheusFormat).
// Subs stores each data structure for a subsctiption
// which contains the annotations as defined by E2tSubscription struct.
type onosE2tSubscriptions struct {
	name        string
	description string
	Labels      []string
	LabelValues []string
	Subs        map[string]E2tSubscription
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for onosE2tSubscriptions.
func (c *onosE2tSubscriptions) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	c.Labels = []string{"id", "revision", "service_model_name", "service_model_version", "node_id", "encoding", "status_phase", "status_state"}
	metricDesc := onose2tBuilder.NewMetricDesc(c.name, c.description, c.Labels, staticLabelsE2t)

	for _, e2tSub := range c.Subs {
		metric := onose2tBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			1,
			e2tSub.Id,
			e2tSub.Revision,
			e2tSub.ServiceModelName,
			e2tSub.ServiceModelVersion,
			e2tSub.E2NodeID,
			e2tSub.Encoding,
			e2tSub.StatusPhase,
			e2tSub.StatusState,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
