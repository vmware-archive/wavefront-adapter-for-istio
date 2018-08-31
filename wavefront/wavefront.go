// Copyright 2018 VMware, Inc.
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
//go:generate $GOPATH/src/istio.io/istio/bin/mixer_codegen.sh -a mixer/adapter/wavefront/config/config.proto -x "-s=false -n wavefront -t metric"

package wavefront

import (
	"context"
	"fmt"
	"net"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/vmware/wavefront-istio-mixer-adapter/wavefront/config"
	wf "github.com/wavefrontHQ/go-metrics-wavefront"

	"google.golang.org/grpc"

	"istio.io/api/mixer/adapter/model/v1beta1"
	policy "istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/template/metric"
	"istio.io/istio/pkg/log"
)

type (
	// Server is basic server interface
	Server interface {
		Addr() string
		Close() error
		Run(shutdown chan error)
	}

	// WavefrontAdapter supports metric template.
	WavefrontAdapter struct {
		listener            net.Listener
		server              *grpc.Server
		reporterInitialized bool
	}
)

// ensure that WavefrontAdapter implements the HandleMetricServiceServer interface.
var _ metric.HandleMetricServiceServer = &WavefrontAdapter{}

// hostTags holds the source tag information for Wavefront metrics.
var hostTags = map[string]string{"source": "istio"}

// createWavefrontReporter creates a reporter that periodically flushes metrics to Wavefront.
func createWavefrontReporter(cfg *config.Params) {
	if direct := cfg.GetDirect(); direct != nil {
		go wf.WavefrontDirect(metrics.DefaultRegistry, cfg.FlushInterval, hostTags, cfg.Prefix, direct.Server, direct.Token)
	} else if proxy := cfg.GetProxy(); proxy != nil {
		addr, _ := net.ResolveTCPAddr("tcp", proxy.Address)
		go wf.WavefrontProxy(metrics.DefaultRegistry, cfg.FlushInterval, hostTags, cfg.Prefix, addr)
	}
}

// verifyAndInitReporter checks if the Wavefront reporter is initialized, and if
// not, initializes it.
func (wa *WavefrontAdapter) verifyAndInitReporter(cfg *config.Params) {
	if !wa.reporterInitialized {
		log.Infof("trying to init wavefront reporter, config: %s", cfg.String())
		if err := config.ValidateCredentials(cfg); err != nil {
			log.Errorf("failed to create wavefront reporter, err: %s, config: %s", err.Error(), cfg.String())
		} else {
			createWavefrontReporter(cfg)
			wa.reporterInitialized = true
			log.Infof("wavefront reporter successfully initialized, config: %s", cfg.String())
		}
	}
}

// identifyMetric locates a metric in an array given the metric name.
func identifyMetric(ms []*config.Params_MetricInfo, name string) *config.Params_MetricInfo {
	for _, m := range ms {
		if m.Name == name {
			return m
		}
	}
	return nil
}

// translateToFloat64 converts a given number to float64 or returns an error.
func translateToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("couldn't convert %s to float64", value)
	}
}

// translateToInt64 converts a given number to int64 or returns an error.
func translateToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("couldn't convert %s to int64", value)
	}
}

// translateSample translates a config.Sample instance to a metrics.Sample instance.
func translateSample(s *config.Params_MetricInfo_Sample) metrics.Sample {
	if def := s.GetExpDecay(); def != nil {
		return metrics.NewExpDecaySample(int(def.ReservoirSize), def.Alpha)
	} else if def := s.GetUniform(); def != nil {
		return metrics.NewUniformSample(int(def.ReservoirSize))
	}
	return nil
}

// writeMetrics extracts metric information from metric.InstanceMsgs and writes
// it to the Wavefront metric registry.
func writeMetrics(cfg *config.Params, insts []*metric.InstanceMsg) {
	for _, inst := range insts {
		metricName := inst.Name
		metric := identifyMetric(cfg.Metrics, metricName)
		if metric == nil {
			log.Warnf("couldn't identify metric %s in configuration %s, ignoring", metricName, cfg.String())
			continue
		}

		value := decodeValue(inst.Value.GetValue())
		tags := decodeTags(inst.Dimensions)

		switch metric.Type {
		case config.GAUGE, config.COUNTER:
			if float64Val, err := translateToFloat64(value); err != nil {
				log.Warnf("couldn't translate metric value: %s %v, err: %v", metricName, value, err)
			} else {
				gauge := wf.GetOrRegisterMetric(metricName, metrics.NewGaugeFloat64(), tags).(metrics.GaugeFloat64)
				gauge.Update(float64Val)
				log.Debugf("updated gauge metric %s with %v, tags: %v", metricName, float64Val, tags)
			}
		case config.DELTA_COUNTER:
			if int64Val, err := translateToInt64(value); err != nil {
				log.Warnf("couldn't translate metric value: %s %v, err: %v", metricName, value, err)
			} else {
				deltaMetricName := wf.DeltaCounterName(metricName)
				counter := wf.GetOrRegisterMetric(deltaMetricName, metrics.NewCounter(), tags).(metrics.Counter)
				counter.Inc(int64Val)
				log.Debugf("updated delta counter metric %s with %v, tags: %v", deltaMetricName, int64Val, tags)
			}
		case config.HISTOGRAM:
			if int64Val, err := translateToInt64(value); err != nil {
				log.Warnf("couldn't translate metric value: %s %v, err: %v", metricName, value, err)
			} else {
				histogram := wf.GetMetric(metricName, tags)
				if histogram == nil {
					sample := translateSample(metric.Sample)
					histogram = metrics.NewHistogram(sample)
					wf.RegisterMetric(metricName, histogram, tags)
				}
				histogram.(metrics.Histogram).Update(int64Val)
				log.Debugf("updated histogram metric %s with %v, tags: %v", metricName, int64Val, tags)
			}
		default:
			log.Warnf("couldn't handle metric %s with value %s, tags: %v", metricName, value, tags)
		}
	}
}

// HandleMetric records metric entries.
func (wa *WavefrontAdapter) HandleMetric(ctx context.Context, r *metric.HandleMetricRequest) (*v1beta1.ReportResult, error) {
	log.Infof("received request %v\n", *r)

	// unmarshal configuration
	cfg := &config.Params{}
	if r.AdapterConfig != nil {
		if err := cfg.Unmarshal(r.AdapterConfig.Value); err != nil {
			log.Errorf("error unmarshalling adapter config: %v", err)
			return nil, err
		}
	}

	// init the Wavefront reporter if not initialized already
	wa.verifyAndInitReporter(cfg)

	// validate the metrics configuration
	if err := config.ValidateMetrics(cfg); err != nil {
		log.Errorf("error validating metrics config: %v %s", err, cfg.String())
		return nil, err
	}

	// write metrics
	writeMetrics(cfg, r.Instances)

	log.Infof("metrics were processed successfully!")
	return &v1beta1.ReportResult{}, nil
}

// decodeTags converts dimensions to a map of tags.
func decodeTags(dimensions map[string]*policy.Value) map[string]string {
	tags := make(map[string]string, len(dimensions))
	for i, d := range dimensions {
		tags[i] = fmt.Sprintf("%v", decodeValue(d.GetValue()))
	}
	return tags
}

// decodeValue decodes a policy.Value instance.
func decodeValue(in interface{}) interface{} {
	switch t := in.(type) {
	case *policy.Value_StringValue:
		return t.StringValue
	case *policy.Value_Int64Value:
		return t.Int64Value
	case *policy.Value_DoubleValue:
		return t.DoubleValue
	case *policy.Value_BoolValue:
		return t.BoolValue
	case *policy.Value_IpAddressValue:
		return t.IpAddressValue
	case *policy.Value_TimestampValue:
		return t.TimestampValue
	case *policy.Value_DurationValue:
		return t.DurationValue
	case *policy.Value_EmailAddressValue:
		return t.EmailAddressValue
	case *policy.Value_DnsNameValue:
		return t.DnsNameValue
	case *policy.Value_UriValue:
		return t.UriValue
	default:
		return fmt.Sprintf("%v", in)
	}
}

// Addr returns the listening address of the server.
func (wa *WavefrontAdapter) Addr() string {
	return wa.listener.Addr().String()
}

// Run starts the server run.
func (wa *WavefrontAdapter) Run(shutdown chan error) {
	shutdown <- wa.server.Serve(wa.listener)
}

// Close gracefully shuts down the server; used for testing.
func (wa *WavefrontAdapter) Close() error {
	if wa.server != nil {
		wa.server.GracefulStop()
	}
	if wa.listener != nil {
		_ = wa.listener.Close()
	}
	return nil
}

// NewWavefrontAdapter creates a new Wavefront adapter that listens at provided port.
func NewWavefrontAdapter(addr string) (Server, error) {
	if addr == "" {
		addr = "0"
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
	if err != nil {
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}

	adapter := &WavefrontAdapter{
		listener:            listener,
		server:              grpc.NewServer(),
		reporterInitialized: false,
	}
	metric.RegisterHandleMetricServiceServer(adapter.server, adapter)
	fmt.Printf("listening on \"%v\"\n", adapter.Addr())
	return adapter, nil
}
