// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package collect

import (
	"bytes"
	"context"
	"fmt"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-exporter/pkg/kpis"
	"google.golang.org/grpc"
)

// onosTopo is the onos xapp pci collector.
// It extracts all the pci related kpis using the Collect method.
type onosTopoCollector struct {
	collector
}

// Collect implements the Collector interface behavior for
// XappPciCollector, returning a list of kpis.KPI.
func (col *onosTopoCollector) Collect() ([]kpis.KPI, error) {
	kpis := []kpis.KPI{}

	if len(col.config.getAddress()) == 0 {
		return kpis, fmt.Errorf("onosTopoCollector Collect missing service address")
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

	entitiesKPI, err := listEntities(conn)
	if err != nil {
		return kpis, err
	}

	relationsKPI, err := listRelations(conn)
	if err != nil {
		return kpis, err
	}

	kpis = append(kpis, entitiesKPI)
	kpis = append(kpis, relationsKPI)

	return kpis, err
}

// listEntities receives a connection to a onos topo service
// to retrieve the topo Entities and store them according to the
// data structure of the kpis.OnosTopoEntities KPI.
func listEntities(conn *grpc.ClientConn) (kpis.KPI, error) {
	entitiesKPI := kpis.OnosTopoEntities()
	entitiesKPI.Entities = make(map[string]kpis.TopoEntity)

	filters := &topoapi.Filters{}
	filters.ObjectTypes = []topoapi.Object_Type{topoapi.Object_ENTITY}
	objects, err := listObjects(conn, filters)

	if err != nil {
		return entitiesKPI, err
	}

	for _, object := range objects {
		entity := parseObjectEntity(object)
		entitiesKPI.Entities[entity.ID] = entity
	}

	return entitiesKPI, nil
}

func parseObjectEntity(obj topoapi.Object) kpis.TopoEntity {
	labels := labelsAsCSV(obj)
	aspects := aspectsAsCSV(obj, false)

	var kindID topoapi.ID
	if e := obj.GetEntity(); e != nil {
		kindID = e.KindID
	}

	return kpis.TopoEntity{
		ID:      string(obj.ID),
		Kind:    string(kindID),
		Labels:  labels,
		Aspects: aspects,
	}
}

// listRelations receives a connection to a onos topo service
// to retrieve the topo Relations and store them according to the
// data structure of the kpis.OnosTopoEntities KPI.
func listRelations(conn *grpc.ClientConn) (kpis.KPI, error) {
	relationsKPI := kpis.OnosTopoRelations()
	relationsKPI.Relations = make(map[string]kpis.TopoRelation)

	filters := &topoapi.Filters{}
	filters.ObjectTypes = []topoapi.Object_Type{topoapi.Object_RELATION}
	objects, err := listObjects(conn, filters)

	if err != nil {
		return relationsKPI, err
	}

	for _, object := range objects {
		relation := parseObjectRelation(object)
		relationsKPI.Relations[relation.ID] = relation
	}

	return relationsKPI, nil
}

func parseObjectRelation(obj topoapi.Object) kpis.TopoRelation {
	labels := labelsAsCSV(obj)
	aspects := aspectsAsCSV(obj, false)
	r := obj.GetRelation()

	return kpis.TopoRelation{
		ID:      string(obj.ID),
		Kind:    string(r.KindID),
		Labels:  labels,
		Source:  string(r.SrcEntityID),
		Target:  string(r.TgtEntityID),
		Aspects: aspects,
	}
}

func listObjects(conn *grpc.ClientConn, filters *topoapi.Filters) ([]topoapi.Object, error) {
	client := topoapi.CreateTopoClient(conn)

	resp, err := client.List(context.Background(), &topoapi.ListRequest{Filters: filters})
	if err != nil {
		return nil, err
	}
	return resp.Objects, nil
}

func labelsAsCSV(object topoapi.Object) string {
	var buffer bytes.Buffer
	first := true
	for k, v := range object.Labels {
		if !first {
			buffer.WriteString(",")
		}
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(v)
		first = false
	}
	return buffer.String()
}

func aspectsAsCSV(object topoapi.Object, verbose bool) string {
	var buffer bytes.Buffer
	first := true
	if object.Aspects != nil {
		for aspectType, aspect := range object.Aspects {

			if !first {
				buffer.WriteString(",")
			}
			buffer.WriteString(aspectType)
			if verbose {
				buffer.WriteString("=")
				buffer.WriteString(bytes.NewBuffer(aspect.Value).String())
			}
			first = false
		}
	}
	return buffer.String()
}
