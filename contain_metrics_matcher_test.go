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
	prommodel "github.com/prometheus/client_model/go"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContainMetricsMatcher", func() {

	famsmap := map[string]*prommodel.MetricFamily{
		"bottled_boris": {
			Type: prommodel.MetricType_COUNTER.Enum(),
			Name: pstr("bottled_boris"),
			Unit: pstr("booze"),
			Metric: []*prommodel.Metric{
				{Label: []*prommodel.LabelPair{{Name: pstr("type"), Value: pstr("champagne")}}},
				{Label: []*prommodel.LabelPair{{Name: pstr("type"), Value: pstr("schaumwein")}}},
			},
		},
		"angry_angie": {
			Type: prommodel.MetricType_GAUGE.Enum(),
			Name: pstr("angry_angie"),
			Unit: pstr("snooze"),
			Metric: []*prommodel.Metric{
				{Label: []*prommodel.LabelPair{{Name: pstr("realm"), Value: pstr("east")}}},
				{Label: []*prommodel.LabelPair{{Name: pstr("realm"), Value: pstr("west")}}},
			},
		},
	}

	It("rejects actual values if they're not families maps", func() {
		Expect(ContainMetrics().Match(nil)).Error().To(MatchError(
			ContainSubstring("ContainMetrics matcher expects a non-nil map of metric families")))
		Expect(ContainMetrics().Match(42)).Error().To(MatchError(
			ContainSubstring("ContainMetrics matcher expects a non-nil map of metric families")))
	})

	It("uses the index when possible", func() {
		Expect(famsmap).To(ContainMetrics(
			Counter(HaveName("bottled_boris"))))
		Expect(famsmap).NotTo(ContainMetrics(
			Counter(HaveName("pritti_prattl"))))
	})

	It("returns the error of a sub matcher", func() {
		Expect(ContainMetrics(Gauge(HaveName(42))).Match(famsmap)).Error().To(MatchError(
			ContainSubstring("to be either a string or GomegaMatcher")))
		Expect(ContainMetrics(Counter(HaveName("bottled_boris"), HaveHelp(42))).Match(famsmap)).Error().To(MatchError(
			ContainSubstring("to be either a string or GomegaMatcher")))
	})

	It("succeeds when expected all metrics are matched", func() {
		Expect(famsmap).To(ContainMetrics(
			Counter(HaveLabel("type=schaumwein")),
			Gauge(HaveLabel("realm=east"))))
		Expect(famsmap).To(ContainMetrics(
			Counter(HaveLabel("type=schaumwein")),
			Gauge(HaveName("angry_angie"), HaveLabel("realm=east"))))
	})

	It("fails when some expected metrics are missing", func() {
		Expect(famsmap).NotTo(ContainMetrics(
			Counter(HaveLabel("type=schaumwein")),
			Gauge(HaveLabel("realm=elsewhere"))))
		Expect(famsmap).NotTo(ContainMetrics(
			Counter(HaveLabel("type=schaumwein")),
			Gauge(HaveName("angry_angie"), HaveLabel("foo=bar"))))
	})

})
