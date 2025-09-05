// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package probe

import (
	"fmt"
	"log"

	"github.com/prometheus-community/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeSystemStatus(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mSystemStatusInfo = prometheus.NewDesc(
			"fortigate_system_status_info",
			"Info about the system status.",
			[]string{"hostname", "status", "serial"}, nil,
		)
		mSystemVersionInfo = prometheus.NewDesc(
			"fortigate_system_version_info",
			"Info about the fortios version",
			[]string{"version", "build"}, nil,
		)
		mSystemHardwareInfo = prometheus.NewDesc(
			"fortigate_system_hardware_info",
			"Info about the make and model.",
			[]string{"model_name", "model_number", "model"}, nil,
		)
	)

	type Result struct {
		ModelName   string `json:"model_name"`
		ModelNumber string `json:"model_number"`
		Model		string `json:"model"`
		Hostname	string `json:"hostname"`
	}

	type SystemStatus struct {
		Result  Result `json:"results"`
		Status   string `json:"status"`
		Serial   string `json:"serial"`
		Version  string `json:"version"`
		Build    int64 `json:"build"`
	}
	var st SystemStatus

	if err := c.Get("api/v2/monitor/system/status", "", &st); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	m = append(m, prometheus.MustNewConstMetric(mSystemStatusInfo, prometheus.GaugeValue, 1.0, 
			st.Result.Hostname, st.Status, st.Serial))
	m = append(m, prometheus.MustNewConstMetric(mSystemVersionInfo, prometheus.GaugeValue, 1.0, 
			st.Version, fmt.Sprintf("%d", st.Build)))
	m = append(m, prometheus.MustNewConstMetric(mSystemHardwareInfo, prometheus.GaugeValue, 1.0, 
			st.Result.ModelName, st.Result.ModelNumber, st.Result.Model))

	return m, true
}
