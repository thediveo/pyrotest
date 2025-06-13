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

func pstr(s string) *string {
	return &s
}

func matchLabels(actual []*prommodel.LabelPair, ms ...MetricPropertyMatcher) (bool, error) {
	GinkgoHelper()
	Expect(ms).NotTo(BeEmpty())
	var lm metricLabelMatcher = &HaveLabelMatcher{} // BeAssignableToTypeOf doesn't work with typed nils.
	lms := make([]metricLabelMatcher, 0, len(ms))
	for _, m := range ms {
		Expect(m).To(BeAssignableToTypeOf(lm), "expected metricLabelMatcher")
		lms = append(lms, m.(metricLabelMatcher))
	}
	return matchAllLabels(actual, lms)
}

var _ = Describe("matching labels", func() {

	labels := []*prommodel.LabelPair{
		{Name: nil, Value: nil},
		{Name: pstr("foo"), Value: pstr("bar")},
		{Name: pstr("bar"), Value: nil},
		{Name: pstr("baz"), Value: pstr("zoo")},
	}

	Context("matching a single label in a set in multiple ways", func() {

		It("doesn't match a non-existing label name", func() {
			Expect(matchLabels(labels, HaveLabel("zoo"))).To(BeFalse())
			Expect(matchLabels(labels, HaveLabel(ContainSubstring("zoo")))).To(BeFalse())
		})

		It("doesn't match a non-existing label value", func() {
			Expect(matchLabels(labels, HaveLabelWithValue("foo", "42"))).To(BeFalse())
			Expect(matchLabels(labels, HaveLabelWithValue("foo", BeEmpty()))).To(BeFalse())
		})

		DescribeTable("matches arguments",
			func(name, value any) {
				Expect(matchLabels(labels, HaveLabelWithValue(name, value))).To(BeTrue())
			},
			Entry(nil, "baz", "zoo"),
			Entry(nil, "bar", ""),
			Entry(nil, "bar", nil),
			Entry(nil, HavePrefix("ba"), "zoo"),
			Entry(nil, "baz", HaveSuffix("oo")),
			Entry(nil, "baz=zoo", nil),
		)

		DescribeTable("failing matching arguments",
			func(name, value any) {
				Expect(matchLabels(labels, HaveLabelWithValue(name, value))).Error().To(HaveOccurred())
			},
			Entry(nil, nil, nil),
			Entry(nil, 42, nil),
			Entry(nil, BeTrue(), nil),
			Entry(nil, "foo", 42),
			Entry(nil, "foo", BeTrue()),
		)

	})

	It("matches (or not) all expected labels", func() {
		Expect(matchLabels(labels, HaveLabel("foo=bar"), HaveLabel("baz=zoo"))).To(BeTrue())
		Expect(matchLabels(labels, HaveLabel("foo=bar"), HaveLabel("baz=booze"))).To(BeFalse())
	})

})
