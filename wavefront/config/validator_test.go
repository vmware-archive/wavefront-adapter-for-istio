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

package config_test

import (
	"fmt"
	"testing"

	"github.com/vmware/wavefront-adapter-for-istio/wavefront/config"
)

func TestValidateCredentials(t *testing.T) {
	table := []struct {
		params config.Params
		err    error
	}{
		{config.Params{Credentials: nil}, config.NoCredentialsError},
		{config.Params{Credentials: &config.Params_Direct{
			Direct: nil,
		}}, config.NoCredentialsError},
		{config.Params{Credentials: &config.Params_Proxy{
			Proxy: nil,
		}}, config.NoCredentialsError},
		{config.Params{Credentials: &config.Params_Direct{
			Direct: &config.Params_WavefrontDirect{Server: ""},
		}}, config.InvalidDirectCredsError},
		{config.Params{Credentials: &config.Params_Direct{
			Direct: &config.Params_WavefrontDirect{
				Server: "not-a-valid-uri",
				Token:  "dummy-token",
			},
		}}, fmt.Errorf("parse not-a-valid-uri: invalid URI for request")},
		{config.Params{Credentials: &config.Params_Direct{
			Direct: &config.Params_WavefrontDirect{Token: ""},
		}}, config.InvalidDirectCredsError},
		{config.Params{Credentials: &config.Params_Direct{
			Direct: &config.Params_WavefrontDirect{
				Server: "https://server.wavefront.com",
				Token:  "dummy-token",
			},
		}}, nil},
		{config.Params{Credentials: &config.Params_Proxy{
			Proxy: &config.Params_WavefrontProxy{Address: ""},
		}}, config.InvalidProxyCredsError},
		{config.Params{Credentials: &config.Params_Proxy{
			Proxy: &config.Params_WavefrontProxy{Address: "not-a-valid-address"},
		}}, fmt.Errorf("address not-a-valid-address: missing port in address")},
		{config.Params{Credentials: &config.Params_Proxy{
			Proxy: &config.Params_WavefrontProxy{Address: "192.168.99.100:8080"},
		}}, nil},
	}

	for _, entry := range table {
		if err := config.ValidateCredentials(&entry.params); fmt.Sprint(err) != fmt.Sprint(entry.err) {
			t.Errorf("Validation failed for %v, got: %v, want: %v.", entry.params, err, entry.err)
		}
	}
}

func TestValidateMetrics(t *testing.T) {
	err := fmt.Errorf("no sample definition was found for histogram metric metric-name")

	table := []struct {
		metric config.Params_MetricInfo
		err    error
	}{
		{config.Params_MetricInfo{Type: config.GAUGE}, nil},
		{config.Params_MetricInfo{Type: config.COUNTER}, nil},
		{config.Params_MetricInfo{Type: config.DELTA_COUNTER}, nil},
		{config.Params_MetricInfo{
			Name: "metric-name",
			Type: config.HISTOGRAM,
		}, err},
		{config.Params_MetricInfo{
			Name:   "metric-name",
			Type:   config.HISTOGRAM,
			Sample: &config.Params_MetricInfo_Sample{},
		}, err},
		{config.Params_MetricInfo{
			Name: "metric-name",
			Type: config.HISTOGRAM,
			Sample: &config.Params_MetricInfo_Sample{
				Definition: &config.Params_MetricInfo_Sample_ExpDecay_{
					ExpDecay: &config.Params_MetricInfo_Sample_ExpDecay{
						ReservoirSize: 1024,
						Alpha:         0.015,
					},
				},
			},
		}, nil},
		{config.Params_MetricInfo{
			Name: "metric-name",
			Type: config.HISTOGRAM,
			Sample: &config.Params_MetricInfo_Sample{
				Definition: &config.Params_MetricInfo_Sample_Uniform_{
					Uniform: &config.Params_MetricInfo_Sample_Uniform{
						ReservoirSize: 1024,
					},
				},
			},
		}, nil},
	}

	for _, entry := range table {
		cfg := &config.Params{Metrics: []*config.Params_MetricInfo{&entry.metric}}
		if err := config.ValidateMetrics(cfg); fmt.Sprint(err) != fmt.Sprint(entry.err) {
			t.Errorf("Validation failed for %v, got: %v, want: %v.", entry.metric, err, entry.err)
		}
	}
}
