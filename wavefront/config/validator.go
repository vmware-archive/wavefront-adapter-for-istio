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

package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
)

var (
	// InvalidDirectCredsError is returned when invalid direct credentials are found in the configuration.
	InvalidDirectCredsError = errors.New("invalid server or token found in configuration")
	// InvalidProxyCredsError is returned when invalid proxy credentials are found in the configuration.
	InvalidProxyCredsError = errors.New("invalid proxy address found in configuration")
	// NoCredentialsError is returned when no valid credentials are found in the configuration.
	NoCredentialsError = errors.New("no credentials were found in the configuration")
)

// ValidateCredentials validates the credentials given a Params instance.
func ValidateCredentials(cfg *Params) error {
	if cfg.GetCredentials() == nil {
		return NoCredentialsError
	} else if cfg.GetDirect() == nil && cfg.GetProxy() == nil {
		return NoCredentialsError
	} else if direct := cfg.GetDirect(); direct != nil {
		if direct.Server == "" || direct.Token == "" {
			return InvalidDirectCredsError
		} else if _, err := url.ParseRequestURI(direct.Server); err != nil {
			return err
		}
	} else if proxy := cfg.GetProxy(); proxy != nil {
		if proxy.Address == "" {
			return InvalidProxyCredsError
		} else if _, err := net.ResolveTCPAddr("tcp", proxy.Address); err != nil {
			return err
		}
	}
	return nil
}

// validateMetric validates a given metric instance.
func validateMetric(m *Params_MetricInfo) error {
	switch m.Type {
	case HISTOGRAM:
		if m.Sample == nil || m.Sample.GetDefinition() == nil || (m.Sample.GetExpDecay() == nil && m.Sample.GetUniform() == nil) {
			return fmt.Errorf("no sample definition was found for histogram metric %s", MetricName(m))
		}
	}
	return nil
}

// ValidateMetrics validates the metrics configuration given a Params instance.
func ValidateMetrics(cfg *Params) error {
	names := make(map[string]struct{})
	instanceNames := make(map[string]struct{})

	for _, m := range cfg.Metrics {
		if err := validateMetric(m); err != nil {
			return err
		}

		// check for duplicate metric names
		name := MetricName(m)
		if _, exists := names[name]; exists {
			return fmt.Errorf("duplicate metric %s found, please supply or change the metric name", name)
		}
		names[name] = struct{}{}

		// check for duplicate instance names
		if _, exists := instanceNames[m.InstanceName]; exists {
			return fmt.Errorf("duplicate metrics found for instance %s", m.InstanceName)
		}
		instanceNames[m.InstanceName] = struct{}{}
	}

	return nil
}
