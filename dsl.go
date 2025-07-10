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
)

// MetricFamilies maps from a metric family name to its metric family.
type MetricsFamilies = map[string]*prommodel.MetricFamily

// BeAMetric succeeds if actual is a Prometheus [*prommodel.MetricFamily] that
// matches the passed-in metric properties in form of a MetricMatcher.
func BeAMetric(m MetricMatcher) types.GomegaMatcher {
	return &BeAMetricMatcher{
		Expected: m,
	}
}

// MetricMatcher succeeds if it matches any metric within a metric family, based
// on family and metric properties. All properties not explicitly specified are
// taken as “don't care”.
//
// Note: when a MetricMatcher matches the name of a metric to a plain string, it
// can report so using the unexported indexname method in order to allow outer
// matchers to optimize finding a matching metric family using a direct
// dictionary lookup.
//
// A MetricMatcher is a composite matcher that succeeds only if all its
// [MetricPropertyMatcher] succeed, such as matching metric family properties
// like name, unit, help, and labels.
type MetricMatcher interface {
	match(*prommodel.MetricFamily) (bool, error)
	indexname() string
}

// MetricPropertyMatcher identifies metric family property or metric property
// matchers.
type MetricPropertyMatcher interface {
	yesimametricpropertymatcher()
}

// Gauge succeeds if a metric (metric family) is a Prometheus Gauge and
// additionally satisfies all optionally specified name, help, and labels
// matchers.
func Gauge(props ...MetricPropertyMatcher) MetricMatcher {
	return metricOfType(prommodel.MetricType_GAUGE, props...)
}

// Counter succeeds if a metric (metric family) is a Prometheus Counter and
// additionally satisfies all optionally specified name, help, and labels
// matchers.
func Counter(props ...MetricPropertyMatcher) MetricMatcher {
	return metricOfType(prommodel.MetricType_COUNTER, props...)
}

// Histogram succeeds if a metric (metric family) is a Prometheus Histogram and
// additionally satisfies all optionally specified name, help, and labels
// matchers.
func Histogram(props ...MetricPropertyMatcher) MetricMatcher {
	return metricOfType(prommodel.MetricType_HISTOGRAM, props...)
}

// HaveLabel succeeds if a metric has a label with the specified name (and
// optional value).
//
// The value passed into the label parameter can be either a string or
// GomegaMatcher:
//   - a string in the form of “name” where it must match a label name, or
//     in the “name=value” form where it must match both the label name and
//     value.
//   - a GomegaMatcher that matches the name only.
//   - any other type of value is an error.
//
// See also [HaveLabelWithValue].
func HaveLabel(label any) MetricPropertyMatcher {
	return newHaveLabelMatcher(label, nil, "HaveLabel")
}

// HaveLabelWithValue succeeds if a metric has a label with the specified name
// and value.
//
// The value passed into the name parameter can be either a string or a
// GomegaMatcher. Similarly, the value passed into the value parameter can also
// be either a string or a GomegaMatcher. Passing any other type of value to
// either the name or value parameter is an error.
//
// See also [HaveLabel].
func HaveLabelWithValue(name, value any) MetricPropertyMatcher {
	return newHaveLabelMatcher(name, value, "HaveLabelWithValue")
}

// HaveName succeeds if a metric family has a name that either equals the passed
// string or matches the passed GomegaMatcher.
func HaveName(name any) MetricPropertyMatcher {
	var plainname string
	if str, ok := name.(string); ok {
		plainname = str
	}
	return &MetricFamilyNameMatcher{
		plainname: plainname,
		matcher:   asStringMatcher(name),
		expected:  name,
	}
}

// HaveUnit succeeds if a metric family has a unit that either equals the passed
// string or matches the passed [types.GomegaMatcher].
func HaveUnit(unit any) MetricPropertyMatcher {
	return &MetricFamilyUnitMatcher{
		matcher:  asStringMatcher(unit),
		expected: unit,
	}
}

// HaveHelp succeeds if a metric family has a help text that either equals the
// passed string or matches the passed [types.GomegaMatcher].
func HaveHelp(help any) MetricPropertyMatcher {
	return &MetricFamilyHelpMatcher{
		matcher:  asStringMatcher(help),
		expected: help,
	}
}

// HaveMetricValue succeeds if a metric is a Counter or a Gauge and has a
// matching float64 value.
func HaveMetricValue(value any) MetricPropertyMatcher {
	return &MetricValueMatcher{
		matcher:  asMatcher(value),
		expected: value,
	}
}

// HaveBucketBoundaries succeeds if a metric is a Histogram and has the matching
// (float64) bucket upper boundaries. Please note the “+Inf” boundary is
// implicit and thus must not be specified in the bucket boundaries slice passed
// to this matcher. If passed a [types.GomegaMatcher], this matcher gets passed
// a slice of []uint64 explicit boundaries as the actual value.
func HaveBucketBoundaries(values any) MetricPropertyMatcher {
	return &HistoryBucketBoundariesMatcher{
		matcher:  asMatcher(values),
		expected: values,
	}
}
