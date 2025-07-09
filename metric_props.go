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

type MetricBucketMatcher struct {
	matcher  types.GomegaMatcher
	expected any
}

var (
	_ (MetricPropertyMatcher) = (*MetricBucketMatcher)(nil)
	_ (metricPropertyMatcher) = (*MetricBucketMatcher)(nil)
	_ (format.GomegaStringer) = (*MetricBucketMatcher)(nil)
)

func (m *MetricBucketMatcher) GomegaString() string {
	if _, ok := m.expected.(types.GomegaMatcher); ok {
		return fmt.Sprintf("buckets: %s", format.Object(m.expected, 1))
	}
	return fmt.Sprintf("buckets: %v", m.expected)
}

func (m *MetricBucketMatcher) yesimametricpropertymatcher() {}

func (m *MetricBucketMatcher) matchProperty(metric *prommodel.Metric) (bool, error) {
	switch {
	case metric.Histogram != nil:
		buckets := metric.Histogram.GetBucket()
		b := make([]uint64, 0, len(buckets))
		for _, bucket := range buckets {
			b = append(b, bucket.GetCumulativeCount())
		}
		return m.matcher.Match(b)
	default:
		return false, nil
	}
}
