// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package collect

import (
	"context"
	"fmt"

	"github.com/onosproject/onos-exporter/pkg/kpis"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/subscription"
	"google.golang.org/grpc"
)

var (
	onose2subConfig = InitConfig("onos-e2sub")
)

// Onose2subCollector is the onos e2sub collector.
// It extracts all the e2sub related kpis using the Collect method.
type Onose2subCollector struct {
	E2subServiceAddress string
}

// Collect implements the Collector interface behavior for
// Onose2subCollector, returning a list of kpis.KPI.
func (col Onose2subCollector) Collect() ([]kpis.KPI, error) {
	kpis := []kpis.KPI{}

	err := onose2subConfig.set(map[string]string{addressKey: col.E2subServiceAddress})
	if err != nil {
		return kpis, err
	}

	conn, err := GetConnection(
		onose2subConfig.getAddress(),
		onose2subConfig.getCertPath(),
		onose2subConfig.getKeyPath(),
		onose2subConfig.noTLS(),
	)
	if err != nil {
		return kpis, err
	}
	defer conn.Close()

	numSubsKPI, err := listSubscriptions(conn)
	if err != nil {
		return kpis, err
	}

	kpis = append(kpis, numSubsKPI)

	return kpis, err
}

// listSubscriptions receives a connection to a onos e2sub service
// to retrieve the e2sub metrics and store them according to the
// data structure of the kpis.OnosE2subs KPI.
func listSubscriptions(conn *grpc.ClientConn) (kpis.KPI, error) {
	onose2subsKPI := kpis.OnosE2subs()
	onose2subsKPI.Subscriptions = make(map[string]kpis.E2Subscription)

	client := subscription.NewClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.List(ctx)
	if err != nil {
		return onose2subsKPI, err
	}

	for _, sub := range response {
		onose2subsKPI.Subscriptions[string(sub.ID)] = kpis.E2Subscription{
			ID:                  string(sub.ID),
			Revision:            fmt.Sprintf("%v", sub.Revision),
			AppID:               string(sub.AppID),
			ServiceModelName:    string(sub.Details.ServiceModel.Name),
			ServiceModelVersion: string(sub.Details.ServiceModel.Version),
			E2NodeID:            string(sub.Details.E2NodeID),
			LifecycleStatus:     fmt.Sprintf("%v", sub.Lifecycle.Status),
		}
	}

	return onose2subsKPI, nil
}
