// Copyright 2018 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wavefront

import (
	"runtime"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	metrics "github.com/rcrowley/go-metrics"
	wf "github.com/wavefronthq/go-metrics-wavefront/reporting"
	"istio.io/istio/pkg/log"
)

// delay between memory and cpu metrics sample
const delay = 5 * time.Second

// createSystemStatsReporter creates a reporter that periodically flushes adapter system metrics to Wavefront.
func createSystemStatsReporter(hostTags map[string]string) {
	log.Info("Preparing adapter metrics")
	ticker := time.NewTicker(delay)
	go func() {
		previous, err := cpu.Get()
		startTime := time.Now()
		if err != nil {
			log.Errorf("Error getting CPU stats - %s", err.Error())
			return
		}

		for t := range ticker.C {
			log.Infof("reporting memory stats - %s", t)
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			gauge := wf.GetOrRegisterMetric("adapter.memory.alloc", metrics.NewGauge(), hostTags).(metrics.Gauge)
			gauge.Update(int64(m.Alloc))

			gauge = wf.GetOrRegisterMetric("adapter.memory.totalalloc", metrics.NewGauge(), hostTags).(metrics.Gauge)
			gauge.Update(int64(m.TotalAlloc))

			gauge = wf.GetOrRegisterMetric("adapter.memory.sys", metrics.NewGauge(), hostTags).(metrics.Gauge)
			gauge.Update(int64(m.Sys))

			gauge = wf.GetOrRegisterMetric("adapter.memory.numgc", metrics.NewGauge(), hostTags).(metrics.Gauge)
			gauge.Update(int64(m.NumGC))

			current, err := cpu.Get()
			if err != nil {
				log.Errorf("Error getting CPU stats - %s", err.Error())
				return
			}
			total := float64(current.Total - previous.Total)

			gaugeCPU := wf.GetOrRegisterMetric("adapter.cpu.user", metrics.NewGaugeFloat64(), hostTags).(metrics.GaugeFloat64)
			gaugeCPU.Update(float64(current.User-previous.User) / total)

			gaugeCPU = wf.GetOrRegisterMetric("adapter.cpu.system", metrics.NewGaugeFloat64(), hostTags).(metrics.GaugeFloat64)
			gaugeCPU.Update(float64(current.System-previous.System) / total)

			gaugeCPU = wf.GetOrRegisterMetric("adapter.cpu.nice", metrics.NewGaugeFloat64(), hostTags).(metrics.GaugeFloat64)
			gaugeCPU.Update(float64(current.Nice-previous.Nice) / total)

			gaugeCPU = wf.GetOrRegisterMetric("adapter.cpu.idle", metrics.NewGaugeFloat64(), hostTags).(metrics.GaugeFloat64)
			gaugeCPU.Update(float64(current.Idle-previous.Idle) / total)

			gaugeCPU = wf.GetOrRegisterMetric("adapter.uptime", metrics.NewGaugeFloat64(), hostTags).(metrics.GaugeFloat64)
			gaugeCPU.Update(time.Since(startTime).Seconds())

			previous = current
		}
	}()
}
