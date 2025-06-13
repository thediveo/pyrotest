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
	"iter"
	"maps"
	"slices"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil/promlint"
	prommodel "github.com/prometheus/client_model/go"

	gi "github.com/onsi/ginkgo/v2"
	gom "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

// CollectAndLint collects metrics from the passed-in [prometheus.Collector],
// linting them, and finally returns them if there are neither errors nor
// linting issues. Otherwise, CollectAndLint will fail the current test with
// details about collecting or linting problems. CollectAndLint uses a newly
// created pedantic [prometheus.Registry].
//
// If any metric names are passed in, only metrics with those names are checked.
func CollectAndLint(coll prometheus.Collector, metricNames ...string) MetricsFamilies {
	gi.GinkgoHelper()
	return collectAndLint(gom.Default, coll, metricNames...)
}

func collectAndLint(gomega types.Gomega, coll prometheus.Collector, metricNames ...string) MetricsFamilies {
	gi.GinkgoHelper()
	reg := prometheus.NewPedanticRegistry()
	gomega.Expect(reg.Register(coll)).To(gom.Succeed(), "registering collector failed")
	return gatherAndLint(gomega, reg, metricNames...)
}

// GatherAndLint gathers all metrics from the passed-in [prometheus.Gatherer],
// linting them, and finally returns them if there are neither errors nor
// linting erros. Otherwise, GatherAndLint will fail the current test with
// details about gathering or linting problems.
//
// If any metric names are passed in, only metrics with those names are checked.
func GatherAndLint(g prometheus.Gatherer, metricNames ...string) MetricsFamilies {
	gi.GinkgoHelper()
	return gatherAndLint(gom.Default, g, metricNames...)
}

func gatherAndLint(gomega types.Gomega, g prometheus.Gatherer, metricNames ...string) MetricsFamilies {
	gi.GinkgoHelper()

	metfams, err := g.Gather()
	gomega.Expect(err).NotTo(gom.HaveOccurred(), "gathering metrics failed")

	if len(metricNames) != 0 {
		metfams = filterMetrics(metfams, metricNames)
	}

	problems, err := promlint.NewWithMetricFamilies(metfams).Lint()
	gomega.Expect(err).NotTo(gom.HaveOccurred(), "linting error")
	gomega.Expect(problems).To(gom.BeEmpty(), "linting problems")

	return maps.Collect(allFamilies(metfams))
}

// allFamilies returns an iterator over all metricfamilies elements, producing
// key-value pairs with the metric family names as keys and the metric families
// as their values.
func allFamilies(metricfamilies []*prommodel.MetricFamily) iter.Seq2[string, *prommodel.MetricFamily] {
	return func(yield func(string, *prommodel.MetricFamily) bool) {
		for _, metfam := range metricfamilies {
			if !yield(metfam.GetName(), metfam) {
				return
			}
		}
	}
}

// filterMetrics returns a new slices with only those MetricFamily elements
// whose names match.
func filterMetrics(metricfamilies []*prommodel.MetricFamily, names []string) []*prommodel.MetricFamily {
	return slices.Collect(func(yield func(metfam *prommodel.MetricFamily) bool) {
		for _, metfam := range metricfamilies {
			if !slices.Contains(names, metfam.GetName()) {
				continue
			}
			if !yield(metfam) {
				return
			}
		}
	})
}
