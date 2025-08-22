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
	prommodel "github.com/prometheus/client_model/go"
	"github.com/thediveo/pyrotest/to"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	testHistoSamples = []float64{
		0.5, 0.11, 0.18, 0.3, 0.33, 0.38, 0.41, 0.42, 0.66, 0.7, 0.9, 1.0, 1.1, 1.2, 2.0,
	}
	testHistoBoundaries = []float64{0.1, 0.2, 0.4, 0.8, 1.6}

	testHistoSamplesSmall     = []float64{0.1}
	testHistoSamplesSmallInf  = []float64{6.66}
	testHistoBoundariesSingle = []float64{0.6}
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
	buckets, count, sum := to.SampledBuckets(testHistoSamples, testHistoBoundaries)
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
	buckets, count, sum = to.SampledBuckets(testHistoSamplesSmall, testHistoBoundariesSingle)
	ch <- prometheus.MustNewConstHistogram(
		prometheus.NewDesc(
			"foos",
			"histogram of foos",
			nil,
			prometheus.Labels{
				"label": "single",
			}),
		count,
		sum,
		buckets)
	ch <- prometheus.MustNewConstHistogram(
		prometheus.NewDesc(
			"foos",
			"histogram of foos",
			nil,
			prometheus.Labels{
				"label": "broken",
			}),
		0, // yes, that's incorrect on purpose.
		sum,
		buckets)
	buckets, count, sum = to.SampledBuckets(testHistoSamplesSmallInf, testHistoBoundariesSingle)
	ch <- prometheus.MustNewConstHistogram(
		prometheus.NewDesc(
			"foos",
			"histogram of foos",
			nil,
			prometheus.Labels{
				"label": "singleinf",
			}),
		count,
		sum,
		buckets)

	ch <- prometheus.MustNewConstHistogram(
		prometheus.NewDesc(
			"foos",
			"histogram of foos",
			nil,
			prometheus.Labels{
				"label": "empty",
			}),
		0, // yes, that's incorrect on purpose.
		0,
		map[float64]uint64{
			66.6: 0,
		})
}

var _ = Describe("single-metric properties", func() {

	It("works with histograms", func() {
		metfams := CollectAndLint(&someCollector{})
		Expect(metfams).To(ContainMetrics(
			Histogram(HaveBucketBoundaries(testHistoBoundaries))))
		Expect(metfams).NotTo(ContainMetrics(
			Histogram(HaveBucketBoundaries(append(testHistoBoundaries, 6.66)))))
	})

	It("checks bucket fillings", func() {
		metfams := CollectAndLint(&someCollector{})
		Expect(metfams).To(ContainMetrics(
			Histogram(HaveLabel("label=single"), HaveSomeFilledBuckets())))
		Expect(metfams).To(ContainMetrics(
			Histogram(HaveLabel("label=singleinf"), HaveSomeFilledBuckets())))
		Expect(metfams).To(ContainMetrics(
			Histogram(HaveLabel("label=broken"), HaveSomeFilledBuckets())))
		Expect(metfams).NotTo(ContainMetrics(
			Histogram(HaveLabel("label=empty"), HaveSomeFilledBuckets())))

		Expect(metfams).NotTo(ContainMetrics(
			Histogram(HaveLabel("label=single"), HaveAllEmptyBuckets())))
		Expect(metfams).NotTo(ContainMetrics(
			Histogram(HaveLabel("label=singleinf"), HaveAllEmptyBuckets())))
		Expect(metfams).To(ContainMetrics(
			Histogram(HaveLabel("label=empty"), HaveAllEmptyBuckets())))
	})

	Context("HistoryEmptyBucketsMatcher", func() {

		It("doesn't match non-histogram metrics", func() {
			Expect(HaveAllEmptyBuckets().(metricPropertyMatcher).matchProperty(&prommodel.Metric{})).
				To(BeFalse())
		})

	})

	Context("HistorySomeFilledBucketsMatcher", func() {

		It("doesn't match non-histogram metrics", func() {
			Expect(HaveSomeFilledBuckets().(metricPropertyMatcher).matchProperty(&prommodel.Metric{})).
				To(BeFalse())
		})

	})

	Context("HistoryBuckedBoundariesMatcher", func() {

		It("doesn't match non-histogram metrics", func() {
			Expect(HaveBucketBoundaries(nil).(metricPropertyMatcher).matchProperty(&prommodel.Metric{})).
				To(BeFalse())
		})

	})

})
