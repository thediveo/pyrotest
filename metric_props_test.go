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
	"github.com/thediveo/pyrotest/to"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	testHistoSamples = []float64{
		0.5, 0.11, 0.18, 0.3, 0.33, 0.38, 0.41, 0.42, 0.66, 0.7, 0.9, 1.0, 1.1, 1.2, 2.0,
	}
	testHistoBoundaries = []float64{0.1, 0.2, 0.4, 0.8, 1.6}
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
	buckets, count, sum := to.BucketsFromSamples(testHistoSamples, testHistoBoundaries)
	ch <- prometheus.MustNewConstHistogram(
		prometheus.NewDesc(
			"foos",
			"histogram of foos",
			nil,
			prometheus.Labels{
				"label": "fools",
			}),
		count,
		sum,
		buckets)
}

var _ = Describe("single-metric properties", func() {

	It("works with histograms", func() {
		metfams := CollectAndLint(&someCollector{})
		Expect(metfams).To(ContainMetrics(
			Histogram(HaveBucketBoundaries(testHistoBoundaries))))
		Expect(metfams).NotTo(ContainMetrics(
			Histogram(HaveBucketBoundaries(append(testHistoBoundaries, 6.66)))))
	})

})
