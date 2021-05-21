// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package collect

import (
	"context"
	"fmt"
	"strings"

	kpimonapi "github.com/onosproject/onos-api/go/onos/kpimon"
	"github.com/onosproject/onos-exporter/pkg/kpis"
	"google.golang.org/grpc"
)

// xappKpimonCollector is the onos xapp kpm collector.
// It extracts all the kpm related kpis using the Collect method.
type xappKpimonCollector struct {
	collector
}

// Collect implements the Collector interface behavior for
// XappKpimonCollector, returning a list of kpis.KPI.
func (col *xappKpimonCollector) Collect() ([]kpis.KPI, error) {
	kpis := []kpis.KPI{}

	if len(col.config.getAddress()) == 0 {
		return kpis, fmt.Errorf("XappKpimonCollector Collect missing service address")
	}

	conn, err := GetConnection(
		col.config.getAddress(),
		col.config.getCertPath(),
		col.config.getKeyPath(),
		col.config.noTLS(),
	)
	if err != nil {
		return kpis, err
	}
	defer conn.Close()

	kpmKPI, err := listKpmMetrics(conn)
	if err != nil {
		return kpis, err
	}

	kpis = append(kpis, kpmKPI)

	return kpis, err

}

// listKpmMetrics receives a connection to a kpm xapp service
// to retrieve the kpm metrics and store them according to the
// data structure of the kpis.XappKpiMon KPI.
func listKpmMetrics(conn *grpc.ClientConn) (kpis.KPI, error) {
	xappKpiMonKPI := kpis.XappKpiMon()
	xappKpiMonKPI.Data = make(map[string]kpis.KpimonData)

	request := kpimonapi.GetRequest{
		Id: "kpimon",
	}
	client := kpimonapi.NewKpimonClient(conn)

	respGetMetrics, err := client.GetMetrics(context.Background(), &request)
	if err != nil {
		return xappKpiMonKPI, err
	}

	for k, v := range respGetMetrics.GetObject().GetAttributes() {
		ids := strings.Split(k, ":")
		tmpCid := ids[0]
		tmpPlmnID := ids[1]
		tmpEgnbID := ids[2]
		tmpMetricType := ids[3]
		// tmpTimestamp := ids[4]

		xappKpiMonKPI.Data[k] = kpis.KpimonData{
			CellID:     tmpCid,
			PlmnID:     tmpPlmnID,
			EgnbID:     tmpEgnbID,
			MetricType: tmpMetricType,
			Value:      v,
		}
	}

	return xappKpiMonKPI, nil
}
