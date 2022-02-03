// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"context"
	"fmt"

	subapi "github.com/onosproject/onos-api/go/onos/e2t/e2/v1beta1"
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

	e2tsubscriptionKPI, err := onose2tListSubscriptions(conn)
	if err != nil {
		return kpis, err
	}

	kpis = append(kpis, e2tsubscriptionKPI)

	return kpis, nil
}

// onose2tListSubscriptions implements the extraction of the kpi OnosE2tSubscriptions
// from the component onose2t. It connects to onos e2t service list the e2NodeSubs
// and fill the proper fields of the OnosE2tSubscriptionsKPI.
// Other functions must be implemented similar to this one in order to extract other
// kpis from onos e2t service.
func onose2tListSubscriptions(conn *grpc.ClientConn) (kpis.KPI, error) {
	OnosE2tSubsKPI := kpis.OnosE2tSubscriptions()
	OnosE2tSubsKPI.Subs = make(map[string]kpis.E2tSubscription)

	client := subapi.NewSubscriptionAdminServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.ListSubscriptions(ctx, &subapi.ListSubscriptionsRequest{})
	if err != nil {
		return OnosE2tSubsKPI, err
	}

	if err != nil {
		return OnosE2tSubsKPI, err
	}

	for _, sub := range response.Subscriptions {

		OnosE2tSubsKPI.Subs[string(sub.ID)] = kpis.E2tSubscription{
			Id:                  string(sub.ID),
			Revision:            string(rune(sub.Revision)),
			ServiceModelName:    string(sub.SubscriptionMeta.ServiceModel.Name),
			ServiceModelVersion: string(sub.SubscriptionMeta.ServiceModel.Version),
			E2NodeID:            string(sub.SubscriptionMeta.E2NodeID),
			Encoding:            sub.SubscriptionMeta.Encoding.String(),
			StatusPhase:         sub.Status.Phase.String(),
			StatusState:         sub.Status.State.String(),
		}
	}

	return OnosE2tSubsKPI, nil
}
