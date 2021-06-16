// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package collect

import (
	"fmt"
	"sync"

	"github.com/onosproject/onos-exporter/pkg/kpis"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

// List of available collectors implemented.
const (
	ONOSE2T        = "onos-e2t"
	ONOSXAPPKPIMON = "onos-xappkpimon"
	ONOSXAPPPCI    = "onos-xapppci"
)

var log = logging.GetLogger("collect")

// Collector defines an interface for Collectors to retrieve
// a list of kpis.KPI via the Collect method.
type Collector interface {
	Collect() ([]kpis.KPI, error)
}

type collector struct {
	name   string
	config Configuration
}

func (col *collector) Collect() ([]kpis.KPI, error) {
	return []kpis.KPI{}, nil
}

// CreateCollector instantiates a new collector based on the const
// name of the collector specified. Available collectors must be defined
// in the cost set of strings.
func CreateCollector(name, serviceAddress string) (Collector, error) {
	colConfig := InitConfig(name)
	err := colConfig.set(map[string]string{addressKey: serviceAddress})

	if err != nil {
		return &collector{}, fmt.Errorf("could not configure collector %s error %s", name, err)

	}

	switch name {
	case ONOSE2T:
		return &onose2tCollector{
			collector: collector{
				name:   name,
				config: colConfig,
			},
		}, nil
	case ONOSXAPPKPIMON:
		return &xappKpimonCollector{
			collector: collector{
				name:   name,
				config: colConfig,
			},
		}, nil
	case ONOSXAPPPCI:
		return &xappPciCollector{
			collector: collector{
				name:   name,
				config: colConfig,
			},
		}, nil
	default:
		return &collector{}, fmt.Errorf("no collector found with name %s", name)

	}
}

// KPIs retrieves the list of kpis.KPI from each Collector.
// It handles each collector error locally, logging the error.
// In any case, kpis.KPI list is returned, e.g., if one collector
// presents error or even if all collectors present errors.
func KPIs(collectors []Collector) []kpis.KPI {
	kpis := []kpis.KPI{}
	mu := sync.RWMutex{}

	wg := sync.WaitGroup{}
	wg.Add(len(collectors))

	for _, col := range collectors {
		go func(c Collector) {
			colKPIs, err := c.Collect()

			if err != nil {
				log.Errorf("collector KPIs Collect error: %s", err)
			} else {
				mu.Lock()
				kpis = append(kpis, colKPIs...)
				mu.Unlock()
			}
			wg.Done()
		}(col)
	}
	wg.Wait()

	return kpis
}
