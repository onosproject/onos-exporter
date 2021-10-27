// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package collect

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-exporter/pkg/kpis"
	"google.golang.org/grpc"
)

// onosTopoCollector is the onos topo collector.
// It extracts all the topo related kpis using the Collect method.
type onosTopoCollector struct {
	collector
}

// Collect implements the Collector interface behavior for
// onosTopoCollector, returning a list of kpis.KPI.
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

	topoEntityObjs, err := getTopoObjects(conn, topoapi.Object_ENTITY)
	if err != nil {
		return kpis, err
	}

	entitiesKPI := listEntities(topoEntityObjs)
	slicesKPI := listSlices(topoEntityObjs)

	kpis = append(kpis, entitiesKPI)
	kpis = append(kpis, slicesKPI)

	relationsKPI, err := listRelations(conn)
	if err != nil {
		return kpis, err
	}

	kpis = append(kpis, relationsKPI)

	return kpis, err
}

// getTopoObjects gets topo objects based on type, which
// can be topoapi.Object_ENTITY or topoapi.Object_RELATION.
func getTopoObjects(conn *grpc.ClientConn, objType topoapi.Object_Type) ([]topoapi.Object, error) {
	entitiesKPI := kpis.OnosTopoEntities()
	entitiesKPI.Entities = make(map[string]kpis.TopoEntity)

	filters := &topoapi.Filters{}
	filters.ObjectTypes = []topoapi.Object_Type{objType}
	objects, err := listObjects(conn, filters)

	return objects, err
}

// listEntities receives a list of topo Objects and store them according to the
// data structure of the kpis.OnosTopoEntities KPI.
func listEntities(objects []topoapi.Object) kpis.KPI {
	entitiesKPI := kpis.OnosTopoEntities()
	entitiesKPI.Entities = make(map[string]kpis.TopoEntity)

	for _, object := range objects {
		entity := parseObjectEntity(object)
		entitiesKPI.Entities[entity.ID] = entity
	}

	return entitiesKPI
}

// listSlices receives a list of topo Objects and store them according to the
// data structure of the kpis.OnosTopoSlices KPI.
func listSlices(objects []topoapi.Object) kpis.KPI {
	slicesKPI := kpis.OnosTopoSlices()
	slicesKPI.Slices = make(map[string]kpis.TopoEntitySlice)

	for _, object := range objects {
		entitySlices := parseSlicesEntity(object)

		for _, entitySlice := range entitySlices {
			sliceKey := fmt.Sprintf("%s-%s", entitySlice.NodeID, entitySlice.SliceID)
			slicesKPI.Slices[sliceKey] = entitySlice
		}

	}

	return slicesKPI
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
// data structure of the kpis.OnosTopoRelations KPI.
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

func parseSliceUeIdList(UeIdList []*topoapi.UeIdentity) []string {
	sliceUEs := []string{}

	for _, UeIdentity := range UeIdList {
		sliceUE := fmt.Sprintf(
			"PreferredIDType=%s,AMFUeNgapID=%s,CuUeF1apID=%s,DuUeF1apID=%s,EnbUeS1apID=%s,RANUeNgapID=%s",
			UeIdentity.PreferredIDType.String(), UeIdentity.AMFUeNgapID.String(), UeIdentity.CuUeF1apID.String(),
			UeIdentity.DuUeF1apID.String(), UeIdentity.EnbUeS1apID.String(), UeIdentity.RANUeNgapID.String())
		sliceUEs = append(sliceUEs, sliceUE)
	}
	return sliceUEs
}

func parseSlicesEntity(object topoapi.Object) []kpis.TopoEntitySlice {

	var kindID topoapi.ID
	if e := object.GetEntity(); e != nil {
		kindID = e.KindID
	}

	NodeID := string(object.ID)
	Kind := string(kindID)

	entitySlices := []kpis.TopoEntitySlice{}

	for aspectType, aspect := range object.Aspects {

		if strings.Contains(aspectType, "RSMSliceItemList") {
			topoSliceItemList := topoapi.RSMSliceItemList{}
			jm := jsonpb.Unmarshaler{}
			avb := bytes.NewBuffer(aspect.Value)
			err := jm.Unmarshal(avb, &topoSliceItemList)
			if err != nil {
				log.Warn(err)
			}

			for _, topoSliceItem := range topoSliceItemList.RsmSliceList {
				log.Info(topoSliceItem)

				sliceUEs := parseSliceUeIdList(topoSliceItem.UeIdList)

				sliceItem := kpis.TopoEntitySlice{
					NodeID:        NodeID,
					Kind:          Kind,
					SliceID:       topoSliceItem.ID,
					SliceDesc:     topoSliceItem.SliceDesc,
					SchedulerType: topoSliceItem.SliceParameters.GetSchedulerType().String(),
					Weight:        fmt.Sprintf("%d", topoSliceItem.SliceParameters.GetWeight()),
					QosLevel:      fmt.Sprintf("%d", topoSliceItem.SliceParameters.GetQosLevel()),
					SliceType:     topoSliceItem.SliceType.String(),
					UeIdList:      strings.Join(sliceUEs, ","),
				}

				entitySlices = append(entitySlices, sliceItem)
			}

		}
	}
	return entitySlices
}
