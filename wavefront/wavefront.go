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
	"bytes"
	"context"
	"fmt"
	"net"
	"os"

	"github.com/vmware/wavefront-istio-mixer-adapter/wavefront/config"

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
		listener net.Listener
		server   *grpc.Server
	}
)

var _ metric.HandleMetricServiceServer = &WavefrontAdapter{}

// HandleMetric records metric entries
func (wa *WavefrontAdapter) HandleMetric(ctx context.Context, r *metric.HandleMetricRequest) (*v1beta1.ReportResult, error) {

	log.Infof("received request %v\n", *r)
	var b bytes.Buffer
	cfg := &config.Params{}

	if r.AdapterConfig != nil {
		if err := cfg.Unmarshal(r.AdapterConfig.Value); err != nil {
			log.Errorf("error unmarshalling adapter config: %v", err)
			return nil, err
		}
	}

	b.WriteString(fmt.Sprintf("HandleMetric invoked with:\n  Adapter config: %s\n  Instances: %s\n",
		cfg.String(), instances(r.Instances)))

	_, err := os.OpenFile("out.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Errorf("error creating file: %v", err)
	}
	f, err := os.OpenFile("out.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Errorf("error opening file for append: %v", err)
	}
	defer f.Close()

	log.Infof("writing instances to file %s", f.Name())
	if _, err = f.Write(b.Bytes()); err != nil {
		log.Errorf("error writing to file: %v", err)
	}

	for _, instance := range r.Instances {
		metricName := instance.Name
		metric := identifyMetric(metricName, cfg.Metrics)
		if metric != nil {
			switch metric.Type {
			case config.GAUGE:
				log.Infof("Gauge %s: %v -- Dimensions: %v", metricName, decodeValue(instance.Value.GetValue()), decodeDimensions(instance.Dimensions))
			case config.COUNTER:
				log.Infof("Counter %s: %v -- Dimensions: %v", metricName, decodeValue(instance.Value.GetValue()), decodeDimensions(instance.Dimensions))
			case config.DELTA_COUNTER:
				log.Infof("Delta Counter %s: %v -- Dimensions: %v", metricName, decodeValue(instance.Value.GetValue()), decodeDimensions(instance.Dimensions))
			case config.HISTOGRAM:
				log.Infof("Histogram %s: %v -- Dimensions: %v", metricName, decodeValue(instance.Value.GetValue()), decodeDimensions(instance.Dimensions))
			default:
				log.Warnf("Couldn't handle metric type %s, data: %v", metric.Type, instance)
			}
		} else {
			log.Warnf("Couldn't identify metric %s", metricName)
		}
	}

	log.Infof("success!!")
	return &v1beta1.ReportResult{}, nil
}

func identifyMetric(name string, metrics []*config.Params_MetricInfo) *config.Params_MetricInfo {
	for _, metric := range metrics {
		if metric.Name == name {
			return metric
		}
	}
	return nil
}

func decodeDimensions(in map[string]*policy.Value) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = decodeValue(v.GetValue())
	}
	return out
}

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

func instances(in []*metric.InstanceMsg) string {
	var b bytes.Buffer
	for _, inst := range in {
		b.WriteString(fmt.Sprintf("'%s':\n"+
			"  {\n"+
			"		Value = %v\n"+
			"		Dimensions = %v\n"+
			"  }", inst.Name, decodeValue(inst.Value.GetValue()), decodeDimensions(inst.Dimensions)))
	}
	return b.String()
}

// Addr returns the listening address of the server
func (wa *WavefrontAdapter) Addr() string {
	return wa.listener.Addr().String()
}

// Run starts the server run
func (wa *WavefrontAdapter) Run(shutdown chan error) {
	shutdown <- wa.server.Serve(wa.listener)
}

// Close gracefully shuts down the server; used for testing
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
		listener: listener,
		server:   grpc.NewServer(),
	}
	metric.RegisterHandleMetricServiceServer(adapter.server, adapter)
	fmt.Printf("listening on \"%v\"\n", adapter.Addr())
	return adapter, nil
}
