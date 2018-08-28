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

package wavefront

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	adapter_integration "istio.io/istio/mixer/pkg/adapter/test"
)

func TestReport(t *testing.T) {
	adptCrBytes, err := ioutil.ReadFile("config/wavefront.yaml")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}

	operatorCfgBytes, err := ioutil.ReadFile("operatorconfig/sample_operator_cfg.yaml")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	operatorCfg := string(operatorCfgBytes)
	shutdown := make(chan error, 1)

	var outFile *os.File
	outFile, err = os.OpenFile("out.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if removeErr := os.Remove(outFile.Name()); removeErr != nil {
			t.Logf("Could not remove temporary file %s: %v", outFile.Name(), removeErr)
		}
	}()

	adapter_integration.RunTest(
		t,
		nil,
		adapter_integration.Scenario{
			Setup: func() (ctx interface{}, err error) {
				pServer, err := NewWavefrontAdapter("")
				if err != nil {
					return nil, err
				}
				go func() {
					pServer.Run(shutdown)
					_ = <-shutdown
				}()
				return pServer, nil
			},
			Teardown: func(ctx interface{}) {
				s := ctx.(Server)
				s.Close()
			},
			ParallelCalls: []adapter_integration.Call{
				{
					CallKind: adapter_integration.REPORT,
					Attrs:    map[string]interface{}{"request.size": 555},
				},
			},
			GetState: func(ctx interface{}) (interface{}, error) {
				// validate if the content of "out.txt" is as expected
				bytes, err := ioutil.ReadFile("out.txt")
				if err != nil {
					return nil, err
				}
				s := string(bytes)
				wantStr := `HandleMetric invoked with:
       Adapter config: &Params{Credentials:&Params_Direct{Direct:&Params_WavefrontDirect{Domain:domain.wavefront.com,Token:api-token,},},FlushInterval:5s,Prefix:testprefix,Metrics:[&Params_MetricInfo{Name:requestsize.instance.istio-system,Type:COUNTER,}],}
       Instances: 'requestsize.instance.istio-system':
       {
           Value = 555
           Dimensions = map[response_code:200]
       }
`
				if normalize(s) != normalize(wantStr) {
					return nil, fmt.Errorf("got adapters state as : '%s'; want '%s'", s, wantStr)
				}
				return nil, nil
			},
			GetConfig: func(ctx interface{}) ([]string, error) {
				s := ctx.(Server)
				return []string{
					// CRs for built-in templates (metric is what we need for this test)
					// are automatically added by the integration test framework.
					string(adptCrBytes),
					strings.Replace(operatorCfg, "{ADDRESS}", s.Addr(), 1),
				}, nil
			},
			Want: `
     {
      "AdapterState": null,
      "Returns": [
       {
        "Check": {
         "Status": {},
         "ValidDuration": 0,
         "ValidUseCount": 0
        },
        "Quota": null,
        "Error": null
       }
      ]
     }`,
		},
	)
}

func normalize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}
