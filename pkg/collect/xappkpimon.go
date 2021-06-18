// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package collect

import (
	"context"
	"fmt"
	"strings"

	prototypes "github.com/gogo/protobuf/types"

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

	request := kpimonapi.GetRequest{}
	client := kpimonapi.NewKpimonClient(conn)

	respGetMeasurement, err := client.ListMeasurements(context.Background(), &request)
	if err != nil {
		return xappKpiMonKPI, err
	}

	for key, measItems := range respGetMeasurement.GetMeasurements() {
		for _, measItem := range measItems.MeasurementItems {
			for _, measRecord := range measItem.MeasurementRecords {
				// timeStamp := measRecord.Timestamp
				measName := measRecord.MeasurementName
				measValue := measRecord.MeasurementValue

				ids := strings.Split(key, ":")
				tmpE2ID := ids[0]
				tmpCellID := ids[1]

				var value interface{}

				switch {
				case prototypes.Is(measValue, &kpimonapi.IntegerValue{}):
					v := kpimonapi.IntegerValue{}
					err := prototypes.UnmarshalAny(measValue, &v)
					if err != nil {
						log.Warn(err)
					}
					value = v.GetValue()

				case prototypes.Is(measValue, &kpimonapi.RealValue{}):
					v := kpimonapi.RealValue{}
					err := prototypes.UnmarshalAny(measValue, &v)
					if err != nil {
						log.Warn(err)
					}
					value = v.GetValue()

				case prototypes.Is(measValue, &kpimonapi.NoValue{}):
					v := kpimonapi.NoValue{}
					err := prototypes.UnmarshalAny(measValue, &v)
					if err != nil {
						log.Warn(err)
					}
					value = v.GetValue()

				}

				uKey := fmt.Sprintf("%s:%s", key, measName)
				xappKpiMonKPI.Data[uKey] = kpis.KpimonData{
					CellID:     tmpCellID,
					E2ID:       tmpE2ID,
					MetricType: measName,
					Value:      fmt.Sprintf("%v", value),
				}
			}
		}
	}

	return xappKpiMonKPI, nil
}
