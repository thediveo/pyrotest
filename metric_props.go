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
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	prommodel "github.com/prometheus/client_model/go"
)

// MetricValueMatcher matches the value property of a Counter or Gauge metric.
type MetricValueMatcher struct {
	matcher  types.GomegaMatcher
	expected any
}

var (
	_ (MetricPropertyMatcher) = (*MetricValueMatcher)(nil)
	_ (metricPropertyMatcher) = (*MetricValueMatcher)(nil)
	_ (format.GomegaStringer) = (*MetricValueMatcher)(nil)
)

func (m *MetricValueMatcher) GomegaString() string {
	if _, ok := m.expected.(types.GomegaMatcher); ok {
		return fmt.Sprintf("value: %s", format.Object(m.expected, 1))
	}
	return fmt.Sprintf("value: %v", m.expected)
}

func (m *MetricValueMatcher) yesimametricpropertymatcher() {}

func (m *MetricValueMatcher) matchProperty(metric *prommodel.Metric) (bool, error) {
	switch {
	case metric.Gauge != nil:
		return m.matcher.Match(metric.Gauge.GetValue())
	case metric.Counter != nil:
		return m.matcher.Match(metric.Counter.GetValue())
	default:
		return false, nil
	}
}

// ----

type HistoryBucketsMatcher struct {
	expectedBuckets any

	countMatcher  types.GomegaMatcher
	expectedCount any
}

var (
	_ (MetricPropertyMatcher) = (*HistoryBucketsMatcher)(nil)
	_ (metricPropertyMatcher) = (*HistoryBucketsMatcher)(nil)
	_ (format.GomegaStringer) = (*HistoryBucketsMatcher)(nil)
)

func (m *HistoryBucketsMatcher) GomegaString() string {
	return "buckets:" // FIXME:
}

func (m *HistoryBucketsMatcher) yesimametricpropertymatcher() {}

func (m *HistoryBucketsMatcher) matchProperty(metric *prommodel.Metric) (bool, error) {
	if metric.Histogram == nil {
		return false, nil
	}

	// Finally match the implicit "+Inf" bucket's count property.
	return m.countMatcher.Match(metric.Histogram.SampleCount)
}

// ----

type HistoryEmptyBucketsMatcher struct{}

var (
	_ (MetricPropertyMatcher) = (*HistoryEmptyBucketsMatcher)(nil)
	_ (metricPropertyMatcher) = (*HistoryEmptyBucketsMatcher)(nil)
	_ (format.GomegaStringer) = (*HistoryEmptyBucketsMatcher)(nil)
)

func (m *HistoryEmptyBucketsMatcher) GomegaString() string {
	return "buckets: <empty>"
}

func (m *HistoryEmptyBucketsMatcher) yesimametricpropertymatcher() {}

func (m *HistoryEmptyBucketsMatcher) matchProperty(metric *prommodel.Metric) (bool, error) {
	if metric.Histogram == nil {
		return false, nil
	}
	if metric.Histogram.GetSampleCount() != 0 {
		return false, nil
	}
	for _, bucket := range metric.Histogram.Bucket {
		if bucket.GetCumulativeCount() != 0 {
			return false, nil
		}
	}
	return true, nil
}

// ----

// HistoryBucketBoundariesMatcher matches the explicit bucket upper boundaries.
type HistoryBucketBoundariesMatcher struct {
	matcher  types.GomegaMatcher
	expected any
}

var (
	_ (MetricPropertyMatcher) = (*HistoryBucketBoundariesMatcher)(nil)
	_ (metricPropertyMatcher) = (*HistoryBucketBoundariesMatcher)(nil)
	_ (format.GomegaStringer) = (*HistoryBucketBoundariesMatcher)(nil)
)

func (m *HistoryBucketBoundariesMatcher) GomegaString() string {
	if _, ok := m.expected.(types.GomegaMatcher); ok {
		return fmt.Sprintf("bucket upper boundaries: %s", format.Object(m.expected, 1))
	}
	return fmt.Sprintf("bucket upper boundaries: %v", m.expected)
}

func (m *HistoryBucketBoundariesMatcher) yesimametricpropertymatcher() {}

func (m *HistoryBucketBoundariesMatcher) matchProperty(metric *prommodel.Metric) (bool, error) {
	switch {
	case metric.Histogram != nil:
		buckets := metric.Histogram.GetBucket()
		upperboundaries := make([]float64, 0, len(buckets))
		for _, bucket := range buckets {
			upperboundaries = append(upperboundaries, bucket.GetUpperBound())
		}
		return m.matcher.Match(upperboundaries)
	default:
		return false, nil
	}
}
