// Copyright 2025 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package pyrotest

import (
	"github.com/prometheus/client_golang/prometheus"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type someCollector struct{}

var _ prometheus.Collector = (*someCollector)(nil)

func (c *someCollector) Describe(ch chan<- *prometheus.Desc) { prometheus.DescribeByCollect(c, ch) }

func (c *someCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"foo_total",
			"there's no help",
			nil,
			prometheus.Labels{
				"label": "scam",
			}),
		prometheus.CounterValue, 42.0)
	buckets := map[float64]uint64{
		0.1: 1,
		0.2: 2,
		0.4: 3,
		0.8: 4,
		1.6: 5,
	}
	sum := 0.5 + 0.11 + 0.18 + 0.3 + 0.33 + 0.38 + 0.41 + 0.42 + 0.66 + 0.7 + 0.9 + 1.0 + 1.1 + 1.2 + 1.3
	ch <- prometheus.MustNewConstHistogram(
		prometheus.NewDesc(
			"foos",
			"histogram of foos",
			nil,
			prometheus.Labels{
				"label": "fools",
			}),
		1+2+3+4+5,
		sum,
		buckets)
}

var _ = Describe("single-metric properties", func() {

	It("works with histograms", func() {
		metfams := CollectAndLint(&someCollector{})
		Expect(metfams).To(ContainMetrics(
			Histogram(HaveBucketValues([]uint64{1, 2, 3, 4, 5}))))
	})

})
