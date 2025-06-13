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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	prommodel "github.com/prometheus/client_model/go"
)

var _ = Describe("BeAMetric", func() {

	family := &prommodel.MetricFamily{
		Type: prommodel.MetricType_COUNTER.Enum(),
		Name: pstr("bottled_boris"),
		Unit: pstr("booze"),
		Help: pstr("beyond any"),
		Metric: []*prommodel.Metric{
			{Label: []*prommodel.LabelPair{{Name: pstr("type"), Value: pstr("champagne")}}},
			{Label: []*prommodel.LabelPair{{Name: pstr("type"), Value: pstr("schaumwein")}}},
		},
	}

	It("rejects an actual nil value", func() {
		Expect(BeAMetric(Gauge()).Match(nil)).Error().To(MatchError(
			MatchRegexp(`BeAMetricMatcher expects a Prometheus MetricFamily.  Got:\n.*<nil>: nil`)))
	})

	It("correctly matches or not", func() {
		Expect(family).NotTo(BeAMetric(Counter(HaveName("angry_angie"))))
		Expect(family).To(BeAMetric(Counter(HaveLabel("type=schaumwein"))))
	})

	It("produces a useful failure message (incl. plain string name matching)", func() {
		m := BeAMetric(Histogram(
			HaveName("abc"),
			HaveHelp("foobar"),
			HaveUnit(Equal("bottle")),
			HaveLabel("foo=bar"),
			HaveLabelWithValue("bar", Equal("baz"))))
		Expect(m.Match(family)).To(BeFalse())
		Expect(m.FailureMessage(family)).To(
			MatchRegexp(
				`HISTOGRAM:
.*name: abc
.*help: foobar
.*unit: .* \{
.*Expected: .*"bottle",
.*\}
.*label \{foo=bar\}
.*label
.*name: .*: bar
.*value: .* \{
.*Expected: .*"baz",
.*\}`))
	})

	It("produces a useful failure message (incl. name GomegaMatcher)", func() {
		m := BeAMetric(Histogram(
			HaveName(Equal("abc")),
			HaveHelp("foobar")))
		Expect(m.Match(family)).To(BeFalse())
		Expect(m.FailureMessage(family)).To(
			MatchRegexp(
				`HISTOGRAM:
.*name: .*: \{
.*Expected: .*"abc",
.*\}
.*help: foobar`))
	})

})
