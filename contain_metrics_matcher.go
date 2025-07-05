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
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// ContainMetricsMatcher is a [types.GomegaMatcher] that succeeds if an actual
// value is assignable to a [MetricsFamilies] map and that metric families map
// contains all expected metrics.
type ContainMetricsMatcher struct {
	ExpectedMetrics []MetricMatcher
	missingMetrics  []MetricMatcher
}

var _ types.GomegaMatcher = (*ContainMetricsMatcher)(nil)

// ContainMetrics succeeds if actual represents a [MetricsFamilies] map and
// contains the passed-in metrics described by [MetricMatcher] elements.
func ContainMetrics(ms ...MetricMatcher) types.GomegaMatcher {
	return &ContainMetricsMatcher{
		ExpectedMetrics: ms,
	}
}

// asFamiliesMap returns the actual value as a MetricsFamilies typed value if
// possible, otherwise it returns false.
func asFamiliesMap(actual any) (MetricsFamilies, bool) {
	if actual == nil {
		return nil, false
	}
	family, ok := actual.(MetricsFamilies)
	return family, ok
}

func (m *ContainMetricsMatcher) Match(actual any) (bool, error) {
	familiesMap, ok := asFamiliesMap(actual)
	if !ok {
		return false, fmt.Errorf(
			"ContainMetrics matcher expects a non-nil map of metric families, indexed by their names.  Got:\n%s",
			format.Object(actual, 1))
	}

	// first, do any fast direct family-by-name lookups where they are
	// possible...
	var lastError error
	slowExpecteds := []MetricMatcher{}
	for _, matcher := range m.ExpectedMetrics {
		name := matcher.indexname()
		if name == "" {
			// we cannot directly look up this metric family by its plain string
			// name, so we put them into the slow list.
			slowExpecteds = append(slowExpecteds, matcher)
			continue
		}
		family, ok := familiesMap[name]
		if !ok || family == nil {
			m.missingMetrics = append(m.missingMetrics, matcher)
			continue
		}
		ok, err := matcher.match(family)
		if err != nil {
			lastError = err
			m.missingMetrics = append(m.missingMetrics, matcher) // kinda superfluous in this case
			continue
		}
		if !ok {
			m.missingMetrics = append(m.missingMetrics, matcher)
			continue
		}
	}
	if len(slowExpecteds) == 0 {
		if lastError != nil {
			return false, lastError
		}
		return len(m.missingMetrics) == 0, nil
	}

	// slow path...
nextFamily:
	for _, family := range familiesMap {
		if len(slowExpecteds) == 0 {
			break
		}
		for idx := 0; idx < len(slowExpecteds); {
			ok, err := slowExpecteds[idx].match(family)
			if err != nil {
				lastError = err
				m.missingMetrics = append(m.missingMetrics, slowExpecteds[idx])
				// move the last matcher element into our current matcher
				// element and then cut off the last element; we then rinse and
				// repeat at the same index, but now with a different matcher.
				slowExpecteds[idx] = slowExpecteds[len(slowExpecteds)-1]
				slowExpecteds = slowExpecteds[:len(slowExpecteds)-1]
			} else if ok {
				// move the last matcher element into our current matcher
				// element and then cut off the last element; we then rinse and
				// repeat at the same index, but now with a different matcher.
				slowExpecteds[idx] = slowExpecteds[len(slowExpecteds)-1]
				slowExpecteds = slowExpecteds[:len(slowExpecteds)-1]
				continue nextFamily
			}
			idx++
		}
	}
	// any error picked up finally takes priority...
	if lastError != nil {
		return false, lastError
	}
	// we were successful when we could match all remaining/slow expected
	// metrics; otherwise, if we still have expected metrics to match but
	// exhausted all our metric families, then we have failed.
	return len(slowExpecteds) == 0 && len(m.missingMetrics) == 0, nil
}

func (m *ContainMetricsMatcher) FailureMessage(actual any) string {
	s := format.Message(actual, "to contain", m.ExpectedMetrics)
	if len(m.missingMetrics) == 0 {
		return s
	}
	return fmt.Sprintf("%s\nthe missing elements were\n%s",
		s, format.Object(m.missingMetrics, 1))
}

func (m *ContainMetricsMatcher) NegatedFailureMessage(actual any) string {
	return format.Message(actual, "not to contain", m.ExpectedMetrics)
}
