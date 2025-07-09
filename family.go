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
	"strings"

	"github.com/onsi/gomega/format"
	prommodel "github.com/prometheus/client_model/go"
)

// metricNamer returns the plain string metric (family) name a
// MetricPropertyMatcher matches on, if any; otherwise, it returns an empty name
// in case a GomegaMatcher is used for more complex name matching, so that there
// is no plain name for quick metric family lookup.
type metricNamer interface {
	indexname() string
}

// metricFamilyPropertyMatcher succeeds if a specific property of a metric
// family match, such as name, unit, help, and labels. Per the Prometheus
// metrics model only the name, unit and help properties apply at the metrics
// family level.
type metricFamilyPropertyMatcher interface {
	matchFamilyProperty(*prommodel.MetricFamily) (bool, error)
}

type metricPropertyMatcher interface {
	matchProperty(*prommodel.Metric) (bool, error)
}

// metricLabelMatcher succeeds if an actual LabelPair matches the specified name
// and optionally the specified value.
type metricLabelMatcher interface {
	matchLabel(*prommodel.LabelPair) (bool, error)
}

// TypedMetricFamilyMatcher implements MetricMatcher to match metrics within a
// metric family that satisfy a mandatory type, optional name, optional
// properties other than name and labels, and finally a set of labels.
type TypedMetricFamilyMatcher struct {
	plainName              string                        // non-zero if plain string to match, otherwise "".
	typ                    prommodel.MetricType          // type of metric, such as counter, gauge, ...
	familyPropertyMatchers []metricFamilyPropertyMatcher // the metric family properties to match.
	metricPropertyMatchers []metricPropertyMatcher       // the properties of the same metric to match.
	labelMatchers          []metricLabelMatcher          // metric labels that must be all matched on the same metric.
}

var (
	_ MetricMatcher         = (*TypedMetricFamilyMatcher)(nil)
	_ format.GomegaStringer = (*TypedMetricFamilyMatcher)(nil)
)

// GomegaString returns an optimized string representation for failure
// reporting, reducing visual clutter compared to simply dumbing the matcher
// using Gomega's format.Object.
func (m *TypedMetricFamilyMatcher) GomegaString() string {
	return fmt.Sprintf("\n%s:%s%s%s",
		m.typ.String(),
		m.expectedName(),
		m.expectedProperties(),
		m.expectedLabels())
}

func (m *TypedMetricFamilyMatcher) expectedName() string {
	if m.plainName == "" {
		// return nothing as name matching is done and documented using a metric
		// family property matcher later.
		return ""
	}
	return "\n" + format.IndentString(fmt.Sprintf("name: %s", m.plainName), 1)
}

func (m *TypedMetricFamilyMatcher) expectedProperties() string {
	if len(m.familyPropertyMatchers) == 0 {
		return ""
	}
	var s strings.Builder
	for _, propMatcher := range m.familyPropertyMatchers {
		s.WriteRune('\n')
		s.WriteString(format.Indent + propMatcher.(format.GomegaStringer).GomegaString())
	}
	for _, propMatcher := range m.metricPropertyMatchers {
		s.WriteRune('\n')
		s.WriteString(format.Indent + propMatcher.(format.GomegaStringer).GomegaString())
	}
	return s.String()
}

func (m *TypedMetricFamilyMatcher) expectedLabels() string {
	if len(m.labelMatchers) == 0 {
		return ""
	}
	var s strings.Builder
	for _, label := range m.labelMatchers {
		s.WriteRune('\n')
		s.WriteString(format.IndentString(label.(format.GomegaStringer).GomegaString(), 1))
	}
	return s.String()
}

// metricOfType returns a new MetricMatcher that matches the specified metrics
// (family) type, and optional metrics (family) properties.
func metricOfType(mettype prommodel.MetricType, props ...MetricPropertyMatcher) MetricMatcher {
	m := &TypedMetricFamilyMatcher{
		typ: mettype,
	}
	// Now separate the metric (family) property matcher into different buckets
	// or even drop them:
	//  - if it's a metricNamer that matches on a plain string name, remember that
	//    name and otherwise drop the matcher completely under the carpet; we'll
	//    later match directly to the plain string name. If it doesn't match on a
	//    plain name then instead keep it as a normal metric (family) property matcher.
	//  - if it's a labelMatcher then put it into its separate list of label matchers.
	//  - everything else is "just" a metric (family) property matcher.
	for _, propm := range props {
		switch matcher := propm.(type) {
		case metricNamer:
			m.plainName = matcher.indexname()
			if m.plainName != "" {
				// don't add name matcher as a normal property matcher as we are
				// returning the plain name as the index key into the families
				// map.
				continue
			}
			m.familyPropertyMatchers = append(m.familyPropertyMatchers, matcher.(metricFamilyPropertyMatcher))
		case metricFamilyPropertyMatcher:
			m.familyPropertyMatchers = append(m.familyPropertyMatchers, matcher)
		case metricPropertyMatcher:
			m.metricPropertyMatchers = append(m.metricPropertyMatchers, matcher)
		case metricLabelMatcher:
			m.labelMatchers = append(m.labelMatchers, matcher)
		default:
			panic(fmt.Sprintf("internal error: unsupported MetricProperyMatcher of type %T", propm))
		}
	}
	return m
}

// match succeeds if the passed MetricFamily...
//   - matches the expected metric type,
//   - matches the expected plain name, if specified,
//   - matches all expected metric family properties (including the name in case of
//     complex name matching).
//   - matches all expected labels within any, but same, metric of this family.
func (m *TypedMetricFamilyMatcher) match(metfam *prommodel.MetricFamily) (bool, error) {
	if metfam.GetType() != m.typ {
		return false, nil
	}
	if m.plainName != "" && m.plainName != metfam.GetName() {
		return false, nil
	}
	for _, propmatcher := range m.familyPropertyMatchers {
		success, err := propmatcher.matchFamilyProperty(metfam)
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}
	// nota bene: on a valid metric family we always have at least one metric;
	// if the test doesn't care about labels at all, we can shortcut things here.
	if len(m.labelMatchers) == 0 {
		return true, nil
	}
	for _, metric := range metfam.GetMetric() {
		success, err := matchAllProperties(metric, m.metricPropertyMatchers)
		if err != nil {
			return false, err
		}
		if !success {
			continue
		}
		success, err = matchAllLabels(metric.GetLabel(), m.labelMatchers)
		if err != nil {
			return false, err
		}
		if success {
			// finally, we passed all exams!!!
			return true, nil
		}
	}
	return false, nil
}

// indexname returns the plain string metrics name if explicitly specified,
// otherwise an empty string in case a Gomega matcher was specified to match
// metric (family) names.
func (m *TypedMetricFamilyMatcher) indexname() string {
	return m.plainName
}

func matchAllProperties(metric *prommodel.Metric, matchers []metricPropertyMatcher) (bool, error) {
	for _, matcher := range matchers {
		success, err := matcher.matchProperty(metric)
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}
	return true, nil
}
