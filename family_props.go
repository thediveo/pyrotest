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
	"errors"
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	prommodel "github.com/prometheus/client_model/go"
)

// MetricFamilyPropertyMatcher matches properties of a metric family (as opposed
// to properties of individual metrics).
type MetricFamilyPropertyMatcher interface {
	match(*prommodel.MetricFamily) (bool, error)
}

// ----

// MetricFamilyNameMatcher matches the name property of a metric family and
// optionally indicates a plain name for direct family lookup.
type MetricFamilyNameMatcher struct {
	plainname string              // either a plain string name to match against...
	matcher   types.GomegaMatcher // ...or a "complex" GomegaMatcher for advanced use cases.
	expected  any                 // original expected value for error reporting.
}

var (
	_ (MetricPropertyMatcher)       = (*MetricFamilyNameMatcher)(nil)
	_ (metricNamer)                 = (*MetricFamilyNameMatcher)(nil)
	_ (metricFamilyPropertyMatcher) = (*MetricFamilyNameMatcher)(nil)
	_ (format.GomegaStringer)       = (*MetricFamilyNameMatcher)(nil)
)

// GomegaString returns an optimized string representation of this metric family
// name matcher as to make failure reporting more useful. In particular it hides
// implementation detail state information from the representation.
func (m *MetricFamilyNameMatcher) GomegaString() string {
	if s, ok := m.expected.(string); ok {
		return fmt.Sprintf("name: %s", s)
	}
	return fmt.Sprintf("name: %s", format.Object(m.expected, 1))
}

func (m *MetricFamilyNameMatcher) yesimametricpropertymatcher() {}

func (m *MetricFamilyNameMatcher) indexname() string {
	return m.plainname
}

func (m *MetricFamilyNameMatcher) matchFamilyProperty(mf *prommodel.MetricFamily) (bool, error) {
	if m.matcher == nil {
		return false, errors.New(format.Message(
			m.expected, "to be either a string or GomegaMatcher"))
	}
	return m.matcher.Match(mf.GetName())
}

// ----

// MetricFamilyHelpMatcher matches the help property of a metric family.
type MetricFamilyHelpMatcher struct {
	matcher  types.GomegaMatcher
	expected any // original expected value for error reporting.
}

var (
	_ (MetricPropertyMatcher)       = (*MetricFamilyHelpMatcher)(nil)
	_ (metricFamilyPropertyMatcher) = (*MetricFamilyHelpMatcher)(nil)
	_ (format.GomegaStringer)       = (*MetricFamilyHelpMatcher)(nil)
)

func (m *MetricFamilyHelpMatcher) GomegaString() string {
	if s, ok := m.expected.(string); ok {
		return fmt.Sprintf("help: %s", s)
	}
	return fmt.Sprintf("help: %s", format.Object(m.expected, 1))
}

func (m *MetricFamilyHelpMatcher) yesimametricpropertymatcher() {}

func (m *MetricFamilyHelpMatcher) matchFamilyProperty(mf *prommodel.MetricFamily) (bool, error) {
	if m.matcher == nil {
		return false, errors.New(format.Message(
			m.expected, "to be either a string or GomegaMatcher"))
	}
	return m.matcher.Match(mf.GetHelp())
}

// ----

// MetricFamilyUnitMatcher matches the unit property of a metric family.
type MetricFamilyUnitMatcher struct {
	matcher  types.GomegaMatcher
	expected any // original expected value for error reporting.
}

var (
	_ (MetricPropertyMatcher)       = (*MetricFamilyUnitMatcher)(nil)
	_ (metricFamilyPropertyMatcher) = (*MetricFamilyUnitMatcher)(nil)
	_ (format.GomegaStringer)       = (*MetricFamilyUnitMatcher)(nil)
)

func (m *MetricFamilyUnitMatcher) GomegaString() string {
	if s, ok := m.expected.(string); ok {
		return fmt.Sprintf("unit: %s", s)
	}
	return fmt.Sprintf("unit: %s", format.Object(m.expected, 1))
}

func (m *MetricFamilyUnitMatcher) yesimametricpropertymatcher() {}

func (m *MetricFamilyUnitMatcher) matchFamilyProperty(mf *prommodel.MetricFamily) (bool, error) {
	if m.matcher == nil {
		return false, errors.New(format.Message(
			m.expected, "to be either a string or GomegaMatcher"))
	}
	return m.matcher.Match(mf.GetUnit())
}
