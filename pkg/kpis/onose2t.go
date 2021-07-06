// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

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

type E2tConnection struct {
	Id             string
	NodeId         string
	PlmnId         string
	RemoteIp       string
	RemotePort     string
	ConnectionType string
}

// onosE2tConnections defines the common data that can be used
// to output the format of a KPI (e.g., PrometheusFormat).
// NumberConnections stores each data structure for a connection
// which contains the annotations as defined by E2tConnection struct.
type onosE2tConnections struct {
	name              string
	description       string
	Labels            []string
	LabelValues       []string
	NumberConnections map[string]E2tConnection
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for onosE2tConnections.
func (c *onosE2tConnections) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	c.Labels = []string{"id", "e2id", "plmnid", "remote_ip", "remote_port", "connection_type"}
	metricDesc := onose2tBuilder.NewMetricDesc(c.name, c.description, c.Labels, staticLabelsE2t)

	for _, e2tCon := range c.NumberConnections {
		metric := onose2tBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			1,
			e2tCon.Id,
			e2tCon.NodeId,
			e2tCon.PlmnId,
			e2tCon.RemoteIp,
			e2tCon.RemotePort,
			e2tCon.ConnectionType,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
