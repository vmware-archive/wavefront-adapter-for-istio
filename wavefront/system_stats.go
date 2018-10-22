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

// nolint:lll
// Generates the wavefront adapter's resource yaml. It contains the adapter's
// configuration, name, supported template names (metric in this case), and
// whether it is session or no-session based.

package wavefront

import (
	"runtime"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	metrics "github.com/rcrowley/go-metrics"
	wf "github.com/wavefrontHQ/go-metrics-wavefront"
	"istio.io/istio/pkg/log"
)

const tickerDuration = 5 * time.Second

var before *cpu.Stats
var after *cpu.Stats
var err error

// CreateSystemStatsReporter creates a reporter that periodically flushes adpater system metrics to Wavefront.
func CreateSystemStatsReporter(hostTags map[string]string) {
	log.Info("Preparing adapter metrics")
	ticker := time.NewTicker(tickerDuration)
	go func() {
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

			if before == nil {
				before, err = cpu.Get()
				if err != nil {
					log.Errorf("Error getting CPU stats - %s", err.Error())
					return
				}
			} else {
				after, err = cpu.Get()
				if err != nil {
					log.Errorf("Error getting CPU stats - %s", err.Error())
					return
				}
				total := float64(after.Total - before.Total)

				gaugeCPU := wf.GetOrRegisterMetric("adapter.cpu.user", metrics.NewGaugeFloat64(), hostTags).(metrics.GaugeFloat64)
				gaugeCPU.Update(float64(after.User-before.User) / total)

				gaugeCPU = wf.GetOrRegisterMetric("adapter.cpu.system", metrics.NewGaugeFloat64(), hostTags).(metrics.GaugeFloat64)
				gaugeCPU.Update(float64(after.System-before.System) / total)

				gaugeCPU = wf.GetOrRegisterMetric("adapter.cpu.nice", metrics.NewGaugeFloat64(), hostTags).(metrics.GaugeFloat64)
				gaugeCPU.Update(float64(after.Nice-before.Nice) / total)

				gaugeCPU = wf.GetOrRegisterMetric("adapter.cpu.idle", metrics.NewGaugeFloat64(), hostTags).(metrics.GaugeFloat64)
				gaugeCPU.Update(float64(after.Idle-before.Idle) / total)

				before = after
			}
		}
	}()
}
