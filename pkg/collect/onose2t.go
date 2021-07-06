// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package collect

import (
	"context"
	"fmt"
	"io"
	"strings"

	adminapi "github.com/onosproject/onos-api/go/onos/e2t/admin"
	"github.com/onosproject/onos-exporter/pkg/kpis"
	"google.golang.org/grpc"
)

// onose2tCollector is the onos e2t collector.
// It extracts all the e2t related kpis using the Collect method.
type onose2tCollector struct {
	collector
}

// Collect implements the collector of the onos e2t service kpis.
// It uses the function(s) defined in onose2t.go to extract the kpis and return
// a list of them.
// This function can create go routines if needed in order to extract multiple
// onos e2t kpis using the same connection and multiple calls to functions
// defined in the file onose2t.go.
func (col *onose2tCollector) Collect() ([]kpis.KPI, error) {
	kpis := []kpis.KPI{}

	if len(col.config.getAddress()) == 0 {
		return kpis, fmt.Errorf("Onose2tCollector Collect missing service address")
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

	e2tconnectionsKPI, err := onose2tListConnections(conn)
	if err != nil {
		return kpis, err
	}

	kpis = append(kpis, e2tconnectionsKPI)

	return kpis, nil
}

// onose2tListConnections implements the extraction of the kpi OnosE2tConnections
// from the component onose2t. It connects to onos e2t service list the e2NodeConnections
// and fill the proper fields of the OnosE2tConnectionsKPI.
// Other functions must be implemented similar to this one in order to extract other
// kpis from onos e2t service.
func onose2tListConnections(conn *grpc.ClientConn) (kpis.KPI, error) {
	OnosE2tConnectionsKPI := kpis.OnosE2tConnections()
	OnosE2tConnectionsKPI.NumberConnections = make(map[string]kpis.E2tConnection)

	request := adminapi.ListE2NodeConnectionsRequest{}
	client := adminapi.NewE2TAdminServiceClient(conn)
	stream, err := client.ListE2NodeConnections(context.Background(), &request)

	if err != nil {
		return OnosE2tConnectionsKPI, err
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return OnosE2tConnectionsKPI, err
		}

		OnosE2tConnectionsKPI.NumberConnections[response.Id] = kpis.E2tConnection{
			Id:             response.Id,
			PlmnId:         response.PlmnId,
			NodeId:         response.NodeId,
			RemoteIp:       strings.Join(response.RemoteIp, ","),
			RemotePort:     fmt.Sprintf("%v", response.RemotePort),
			ConnectionType: response.ConnectionType.String(),
		}
	}

	return OnosE2tConnectionsKPI, nil
}
