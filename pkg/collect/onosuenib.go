// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package collect

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/onosproject/onos-api/go/onos/uenib"
	"github.com/onosproject/onos-exporter/pkg/kpis"
	"google.golang.org/grpc"
)

// onosUenibCollector is the onos uenib collector.
// It extracts all the uenib related kpis using the Collect method.
type onosUenibCollector struct {
	collector
}

// Collect implements the Collector interface behavior for
// onosUenibCollector, returning a list of kpis.KPI.
func (col *onosUenibCollector) Collect() ([]kpis.KPI, error) {
	kpis := []kpis.KPI{}

	if len(col.config.getAddress()) == 0 {
		return kpis, fmt.Errorf("onosUenibCollector Collect missing service address")
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

	uenibKPI, err := listUEs(conn)
	if err != nil {
		return kpis, err
	}

	kpis = append(kpis, uenibKPI)

	return kpis, err
}

// listUEs receives a connection to a onos uenib service
// to retrieve the uenib UEs Aspects and store them according to the
// data structure of the kpis.OnosUenibUEs KPI.
func listUEs(conn *grpc.ClientConn) (kpis.KPI, error) {
	uenibKPI := kpis.OnosUenibUEs()
	uenibKPI.UEs = make(map[string]kpis.UE)

	aspectTypes := []string{"neighbors", "RRC.Conn.Avg"}

	client := uenib.CreateUEServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := client.ListUEs(ctx, &uenib.ListUERequest{AspectTypes: aspectTypes})
	if err != nil {

		return uenibKPI, err
	}

	if err != nil {
		return uenibKPI, err
	}

	for {
		resp, err := response.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return uenibKPI, err
		} else {
			ue := parseObjectUE(resp.UE)
			uenibKPI.UEs[ue.ID] = ue

		}
	}

	return uenibKPI, nil
}

func parseObjectUE(ue uenib.UE) kpis.UE {
	aspects := []string{}
	aspectsValues := []string{}

	for aspectType, any := range ue.Aspects {
		aspectType = strings.ToLower(strings.ReplaceAll(aspectType, ".", "_"))
		aspects = append(aspects, aspectType)
		aspectsValues = append(aspectsValues, string(any.Value))
	}

	return kpis.UE{
		ID:            string(ue.ID),
		Aspects:       aspects,
		AspectsValues: aspectsValues,
	}
}
