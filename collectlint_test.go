// Copyright 2025 Harald Albrecht.
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

package pyrotest

import (
	"maps"

	"github.com/prometheus/client_golang/prometheus"
	prommodel "github.com/prometheus/client_model/go"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type brokenCollector struct{}

var _ prometheus.Collector = (*brokenCollector)(nil)

func (c *brokenCollector) Describe(ch chan<- *prometheus.Desc) { prometheus.DescribeByCollect(c, ch) }

func (c *brokenCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"foo_total",
			"there's no help",
			nil,
			prometheus.Labels{
				"label": "scam",
			}),
		prometheus.CounterValue, 42.0)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"foo_total",
			"42",
			nil,
			prometheus.Labels{
				"foobar": "baz",
			}),
		prometheus.GaugeValue, 66.6)
}

type fishyCollector struct{}

var _ prometheus.Collector = (*brokenCollector)(nil)

func (c *fishyCollector) Describe(ch chan<- *prometheus.Desc) { prometheus.DescribeByCollect(c, ch) }

func (c *fishyCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"foo",
			"there's no help",
			nil,
			prometheus.Labels{
				"label": "scam",
			}),
		prometheus.CounterValue, 42.0)
}

var _ = Describe("collecting/gathering, and linting metric families", func() {

	metfams := []*prommodel.MetricFamily{
		{Name: pstr("foo")},
		{Name: pstr("bar")},
		{Name: pstr("baz")},
	}

	It("filters metric families", func() {
		Expect(filterMetrics(metfams, []string{"foo", "baz", "zoo"})).To(ConsistOf(
			HaveField("GetName()", "foo"),
			HaveField("GetName()", "baz")))
	})

	It("iterates over metric families", func() {
		m := maps.Collect(allFamilies(metfams))
		Expect(m).To(HaveLen(len(metfams)))
		Expect(m).To(Equal(maps.Collect(
			func(yield func(string, *prommodel.MetricFamily) bool) {
				for _, metfam := range metfams {
					if !yield(metfam.GetName(), metfam) {
						return
					}
				}
			})))
	})

	When("things fail", Serial, func() {

		var g Gomega
		var msg string

		BeforeEach(func() {
			g = NewGomega(func(message string, callerSkip ...int) { msg = message })
		})

		It("fails given an invalid collector", func() {
			collectAndLint(g, &brokenCollector{})
			Expect(msg).To(ContainSubstring("have inconsistent label names or help strings"))
		})

		It("fails given linting problems", func() {
			collectAndLint(g, &fishyCollector{})
			Expect(msg).To(ContainSubstring("counter metrics should have \\\"_total\\\" suffix"))
		})

	})

})
