// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package collect

import (
	"fmt"
	"strings"

	"github.com/onosproject/onos-exporter/pkg/kpis"

	"github.com/google/pprof/profile"
	"github.com/onosproject/onos-exporter/pkg/internal/driver"
	"github.com/onosproject/onos-exporter/pkg/internal/plugin"
	"github.com/onosproject/onos-exporter/pkg/internal/report"
)

type onosProfileCollector struct {
	collector
}

func formatAddress(address, format string) (string, error) {
	switch format {
	case "heap":
		newFormat := "http://" + address + ":6060/debug/pprof/heap"
		return newFormat, nil
	case "goroutine":
		newFormat := "http://" + address + ":6060/debug/pprof/goroutine"
		return newFormat, nil
	case "cpu":
		newFormat := "http://" + address + ":6060/debug/pprof/profile?seconds=2"
		return newFormat, nil
	default:
		return "", fmt.Errorf("no address profile target format for %s", format)
	}
}

func (col *onosProfileCollector) Collect() ([]kpis.KPI, error) {
	kpis := []kpis.KPI{}

	if len(col.config.getAddress()) == 0 {
		return kpis, fmt.Errorf("OnosProfileCollector Collect missing service address(es)")
	}

	heapKPIs, err := onosProfiles(col.config.getAddress())
	if err != nil {
		return kpis, err
	}

	kpis = append(kpis, heapKPIs)

	return kpis, nil

}

func onosProfiles(addresses string) (kpis.KPI, error) {
	onosProfileHeapKPI := kpis.OnosProfileHeap()
	onosProfileHeapKPI.Objects = make(map[string]kpis.HeapObject)

	// Remove any temporary files created during pprof processing.
	defer func() {
		err := driver.CleanupTempFiles()
		if err != nil {
			log.Warn("onosProfiles could not cleam temp files")
		}
	}()

	profileTypes := []string{"heap", "cpu", "goroutine"}

	var addressesSplit []string
	if strings.Contains(addresses, ",") {
		addressesSplit = strings.Split(addresses, ",")
	} else {
		addressesSplit = []string{addresses}
	}

	for _, address := range addressesSplit {
		for _, profileType := range profileTypes {

			profileAddress, err := getProfile(address, profileType)
			if err != nil {
				return onosProfileHeapKPI, err
			}

			for _, prof := range profileAddress.objects {
				obj := kpis.HeapObject{
					Name:   prof.name,
					Value:  prof.value,
					Source: address,
					Format: profileType,
				}
				objID := strings.Join([]string{address, profileType, prof.name}, "-")
				onosProfileHeapKPI.Objects[objID] = obj
			}
		}
	}
	return onosProfileHeapKPI, nil
}

func getProfile(address, profileType string) (profiles, error) {
	eo := &plugin.Options{}

	profs := profiles{
		format: profileType,
	}
	profs.objects = make(map[string]profileObject)

	fmtAddress, err := formatAddress(address, profileType)
	if err != nil {
		return profs, err
	}

	if profileType == "heap" {
		cfg := driver.CurrentConfig()
		cfg.SampleIndex = "inuse_space"
	}

	o := driver.SetDefaults(eo)
	cmd := []string{"text"}
	src := &driver.Source{
		Sources:            []string{fmtAddress},
		ExecName:           "",
		BuildID:            "",
		Seconds:            -1,
		Timeout:            -1,
		Symbolize:          "",
		HTTPHostport:       "",
		HTTPDisableBrowser: true,
		Comment:            "",
	}

	p, err := driver.FetchProfiles(src, o)
	if err != nil {
		return profs, err
	}

	if cmd != nil {
		err = reportMetrics(p, cmd, o, profs)
		return profs, err
	}

	return profs, nil
}

func reportMetrics(p *profile.Profile, cmd []string, o *plugin.Options, profs profiles) error {
	cfg := driver.CurrentConfig()

	_, rpt, err := driver.GenerateRawReport(p, cmd, cfg, o)
	if err != nil {
		return err
	}

	items, _ := report.TextItems(rpt)
	for _, item := range items {
		obj := profileObject{
			name:  item.Name,
			value: item.Flat,
		}
		profs.objects[item.Name] = obj
	}

	return nil
}

type profileObject struct {
	name  string
	value int64
}

type profiles struct {
	format  string
	objects map[string]profileObject
}
