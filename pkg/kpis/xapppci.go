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

type CellConflict struct {
	CellID            string
	ResolvedConflicts float64
	OriginalPci       string
	ResolvedPci       string
}

type CellInfo struct {
	CellID        string
	NodeID        string
	CellType      string
	CellPci       string
	CellDlearfcn  float64
	CellNeighbors string
}

// xapppciNumConflicts defines the common data that can be used
// to output the format of a KPI (e.g., PrometheusFormat).
// CellInfo stores the cell info.
type xappPciNumConflicts struct {
	name        string
	description string
	Labels      []string
	LabelValues []string
	Cells       map[string]CellInfo
}

// xappPciResolvedConflicts defines the common data that can be used
// to output the format of a KPI (e.g., PrometheusFormat).
// CellConflict stores the number of conflicts per cell id.
type xappPciResolvedConflicts struct {
	name        string
	description string
	Labels      []string
	LabelValues []string
	Cells       map[string]CellConflict
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for xapppciNumConflicts.
func (c *xappPciNumConflicts) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	c.Labels = []string{"cellid", "celltype", "e2id", "pci", "neighbors"}
	metricDesc := xappPciBuilder.NewMetricDesc(c.name, c.description, c.Labels, staticLabelsXappPci)

	for _, cell := range c.Cells {
		metric := xappPciBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			cell.CellDlearfcn,
			cell.CellID,
			cell.CellType,
			cell.NodeID,
			cell.CellPci,
			cell.CellNeighbors,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// PrometheusFormat implements the contract behavior of the kpis.KPI
// interface for xappPciResolvedConflicts.
func (c *xappPciResolvedConflicts) PrometheusFormat() ([]prometheus.Metric, error) {
	metrics := []prometheus.Metric{}

	c.Labels = []string{"cellid", "original_pci", "resolved_pci"}
	metricDesc := xappPciBuilder.NewMetricDesc(c.name, c.description, c.Labels, staticLabelsXappPci)

	for _, cell := range c.Cells {
		metric := xappPciBuilder.MustNewConstMetric(
			metricDesc,
			prometheus.GaugeValue,
			cell.ResolvedConflicts,
			cell.CellID,
			cell.OriginalPci,
			cell.ResolvedPci,
		)
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
