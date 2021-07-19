// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpis

import (
	"strconv"
	"strings"

	"github.com/onosproject/onos-lib-go/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
)

// Var definitions of xapp kpimon metrics builder and static labels.
// builder is used to create metrics in the PrometheusFormat.
var (
	staticLabelsXappKpimon = map[string]string{"sdran": "xappkpimon"}
	xappKpimonBuilder      = prom.NewBuilder("onos", "xappkpimon", staticLabelsXappKpimon)
)

type KpimonData struct {
	NodeID       string
	CellID       string
	CellGlobalID string
	MetricType   string
	Value        string
}

// xappkpimon defines the common data that can be used
// to output the format of a KPI (e.g., PrometheusFormat).
// Data stores the KpimonData structure defined for each kpimon
// metric.
type xappkpimon struct {
	name        string
	description string
	Labels      []string
	LabelValues []string
	Data        map[string]KpimonData
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for xappkpimon.
func (c *xappkpimon) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	for _, data := range c.Data {
		metricName := strings.ReplaceAll(strings.ToLower(data.MetricType), ".", "_")
		metricValue, err := strconv.ParseInt(data.Value, 0, 64)

		if err != nil {
			return metrics, err
		}

		c.Labels = []string{"nodeid", "cellid", "cell_global_id"}
		metricDesc := xappKpimonBuilder.NewMetricDesc(metricName, c.description, c.Labels, staticLabelsXappKpimon)

		metric := xappKpimonBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			float64(metricValue),
			data.NodeID,
			data.CellID,
			data.CellGlobalID,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
