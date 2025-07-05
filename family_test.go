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
	"github.com/onsi/gomega/types"
	prommodel "github.com/prometheus/client_model/go"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type invalidM struct{}

func (m *invalidM) yesimametricpropertymatcher() {}

var _ = Describe("matching metrics in families", func() {

	labels := []*prommodel.LabelPair{
		{Name: nil, Value: nil},
		{Name: pstr("foo"), Value: pstr("bar")},
		{Name: pstr("bar"), Value: nil},
		{Name: pstr("baz"), Value: pstr("zoo")},
	}

	counterFamily := &prommodel.MetricFamily{
		Name: pstr("foo_bar_total"),
		Type: prommodel.MetricType_COUNTER.Enum(),
		Unit: pstr("gotchas"),
		Help: pstr("help!"),
		Metric: []*prommodel.Metric{
			{
				Label: labels,
			},
			{
				Label: []*prommodel.LabelPair{
					{Name: pstr("foobar"), Value: pstr("barbaz")},
				},
			},
		},
	}

	It("rejects a misconfigured metric name matcher", func() {
		Expect((&MetricFamilyNameMatcher{}).matchProperty(counterFamily)).Error().To(HaveOccurred())
	})

	It("rejects an invalid matcher", func() {
		Expect(func() { Gauge(&invalidM{}) }).To(PanicWith(ContainSubstring("internal error")))
	})

	DescribeTable("incorrectly configured property matchers",
		func(m metricPropertyMatcher) {
			Expect(m.matchProperty(counterFamily)).Error().To(MatchError(
				ContainSubstring("to be either a string or GomegaMatcher")))
		},
		Entry("help property", HaveHelp(666)),
		Entry("help property", HaveUnit(42)),
	)

	DescribeTable("matching the metric name",
		func(m metricPropertyMatcher, matchExpectations types.GomegaMatcher) {
			Expect(m.matchProperty(counterFamily)).To(matchExpectations)
		},
		Entry("correct name", HaveName("foo_bar_total"), BeTrue()),
		Entry("correct name", HaveName(HavePrefix("foo_")), BeTrue()),
		Entry("wrong name", HaveName("rumpelpumpel"), BeFalse()),
	)

	DescribeTable("returning a families index",
		func(m metricPropertyMatcher, expected string) {
			Expect(m.(metricNamer).indexname()).To(Equal(expected))
		},
		Entry("plain name", HaveName("foobar"), "foobar"),
		Entry("name matcher", HaveName(Equal("foobar")), ""),
	)

	DescribeTable("matching the metric type",
		func(m MetricMatcher, matchExpectations types.GomegaMatcher) {
			Expect(m.match(counterFamily)).To(matchExpectations)
		},
		Entry("correct counter", Counter(), BeTrue()),
		Entry("wrong gauge", Gauge(), BeFalse()),
		Entry("wrong history", Histogram(), BeFalse()),
	)

	DescribeTable("matching the metric unit",
		func(m metricPropertyMatcher, matchExpectations types.GomegaMatcher) {
			Expect(m.matchProperty(counterFamily)).To(matchExpectations)
		},
		Entry("correct unit", HaveUnit("gotchas"), BeTrue()),
		Entry("correct unit regexp", HaveUnit(MatchRegexp(`go.*as`)), BeTrue()),
		Entry("wrong unit", HaveUnit("gorillas"), BeFalse()),
	)

	DescribeTable("matching the help",
		func(m metricPropertyMatcher, matchExpectations types.GomegaMatcher) {
			Expect(m.matchProperty(counterFamily)).To(matchExpectations)
		},
		Entry("correct help", HaveHelp("help!"), BeTrue()),
		Entry("correct help regexp", HaveHelp(MatchRegexp(`he..!`)), BeTrue()),
		Entry("wrong help", HaveHelp("hah!"), BeFalse()),
	)

	DescribeTable("matching a family",
		func(m MetricMatcher, matchExpectations types.GomegaMatcher) {
			ok, err := m.match(counterFamily)
			if err != nil {
				Expect(err).To(matchExpectations)
				return
			}
			Expect(ok).To(matchExpectations)
		},
		Entry("matching Counter",
			Counter(HaveName(*counterFamily.Name), HaveUnit(*counterFamily.Unit)),
			BeTrue()),
		Entry("matching Counter",
			Counter(HaveName(Equal(*counterFamily.Name)), HaveUnit(*counterFamily.Unit)),
			BeTrue()),
		Entry("wrong Counter name",
			Counter(HaveName("ballistic_boris"), HaveUnit(*counterFamily.Unit)),
			BeFalse()),
		Entry("wrong Counter unit",
			Counter(HaveUnit("kiloseconds")),
			BeFalse()),
		Entry("Gomega matcher error",
			Counter(HaveName(BeTrue()), HaveUnit(*counterFamily.Unit)),
			MatchError(ContainSubstring("Expected a boolean"))),
	)

	Context("labels", func() {

		It("matches metric's labels (or not)", func() {
			Expect(Counter( /* sic! */ ).match(counterFamily)).To(BeTrue())
			Expect(Counter(HaveLabel("baz=zoo"), HaveLabel("foo=bar")).match(counterFamily)).To(BeTrue())
			Expect(Counter(HaveLabel("foobar=barbaz")).match(counterFamily)).To(BeTrue())
			Expect(Counter(HaveLabel("baz=zoo"), HaveLabel("fool=bar")).match(counterFamily)).To(BeFalse())
		})

		It("reports metric label matcher errors", func() {
			Expect(Counter(HaveLabel(BeTrue())).match(counterFamily)).Error().To(HaveOccurred())
		})

	})

})
