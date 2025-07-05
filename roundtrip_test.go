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

package pyrotest_test

import (
	"github.com/prometheus/client_golang/prometheus"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/thediveo/pyrotest" // this ensures that dot importing doesn't cause any conflicts
)

type testCollector struct{}

var _ prometheus.Collector = (*testCollector)(nil)

func (c *testCollector) Describe(ch chan<- *prometheus.Desc) { prometheus.DescribeByCollect(c, ch) }

func (c *testCollector) Collect(ch chan<- prometheus.Metric) {
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
			"bar_baz",
			"42",
			nil,
			prometheus.Labels{
				"foobar": "baz",
			}),
		prometheus.GaugeValue, 66.6)
}

var _ = Describe("roundtrip", func() {

	It("collects and reasons about metrics", func() {
		c := &testCollector{}
		families := CollectAndLint(c)
		Expect(families).To(ContainMetrics(
			Gauge(HaveName("bar_baz"),
				HaveHelp(Not(BeEmpty())),
				HaveLabel("foobar=baz")),
			Counter(HaveName(ContainSubstring("_total")),
				HaveHelp(ContainSubstring("no help")),
				HaveLabelWithValue("label", "scam")),
		))

		Expect(CollectAndLint(c, "foo_total")).To(HaveLen(1))
	})

	It("gathers and reasons about metrics", func() {
		c := &testCollector{}
		reg := prometheus.NewPedanticRegistry()
		Expect(reg.Register(c)).To(Succeed())
		Expect(GatherAndLint(reg, "foo_total")).To(HaveLen(1))
	})

})
