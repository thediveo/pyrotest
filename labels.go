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

	prommodel "github.com/prometheus/client_model/go"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// HaveLabelMatcher succeeds if it matches an actual metric label by name and
// optionally by value. The HaveLabelMatcher supports describing the name and
// value matches as either verbatim string matches or alternatively using
// [types.GomegaMatcher].
type HaveLabelMatcher struct {
	name         any
	value        any
	matcherName  string
	nameMatcher  types.GomegaMatcher
	valueMatcher types.GomegaMatcher
}

var (
	_ MetricPropertyMatcher = (*HaveLabelMatcher)(nil)
	_ metricLabelMatcher    = (*HaveLabelMatcher)(nil)
	_ format.GomegaStringer = (*HaveLabelMatcher)(nil)
)

func newHaveLabelMatcher(name, value any, matchername string) MetricPropertyMatcher {
	if s, ok := name.(string); ok && value == nil {
		// single plain string argument, so let's see if it is in "NAME=VALUE"
		// format...?
		if nam, val, found := strings.Cut(s, "="); found {
			return &HaveLabelMatcher{
				name:         nam,
				value:        val,
				matcherName:  matchername,
				nameMatcher:  asStringMatcher(nam),
				valueMatcher: asStringMatcher(val),
			}
		}
	}
	return &HaveLabelMatcher{
		name:         name,
		value:        value,
		matcherName:  matchername,
		nameMatcher:  asStringMatcher(name),
		valueMatcher: asStringMatcher(value),
	}
}

func (m *HaveLabelMatcher) yesimametricpropertymatcher() {}

// GomegaString returns an optimized string representation for failure
// reporting, reducing visual clutter as much as possible. In case both the
// expected name and the value are plain string values. Otherwise, it falls back
// to reporting both name and value using Gomega's [format.Object]. In any case,
// it never reports useless private state, such as the name of the matcher
// constructor used, or the derived name and value matcher objects.
func (m *HaveLabelMatcher) GomegaString() string {
	nameStr, nameOk := m.name.(string)
	valueStr, valueOk := m.value.(string)
	if nameOk && valueOk {
		return fmt.Sprintf("label {%s=%s}", nameStr, valueStr)
	}
	return fmt.Sprintf("label\n%sname: %s\n%svalue: %s",
		format.Indent, format.Object(m.name, 1),
		format.Indent, format.Object(m.value, 1))
}

// matchLabel returns true if the passed-in label matches the configured name
// and optionally the value. Otherwise, it returns itself any errors returned by
// and GomegaMatcher specified for matching a label name and/or value.
func (m *HaveLabelMatcher) matchLabel(label *prommodel.LabelPair) (bool, error) {
	if m.nameMatcher == nil {
		return false, fmt.Errorf("name matcher must not be <nil>")
	}
	if m.value != nil && m.valueMatcher == nil {
		return false, fmt.Errorf("expected value to match to be either string or types.GomegaMatcher.  Got:\n%T",
			m.value)
	}
	success, err := m.nameMatcher.Match(label.GetName())
	if err != nil {
		return false, err
	}
	if !success {
		return false, nil
	}
	if m.valueMatcher == nil { // no value to match, so we've found a matching label
		return true, nil
	}
	success, err = m.valueMatcher.Match(label.GetValue())
	if err != nil {
		return false, err
	}
	return success, nil
}

// matchAllLabels succeeds if the all expected labels match (a subset of) the
// actual labels. It returns an error as soon as any underlying label matcher
// returns an error.
func matchAllLabels(actual []*prommodel.LabelPair, expected []metricLabelMatcher) (bool, error) {
nextMatcher:
	for _, matcher := range expected {
		for _, lblPair := range actual {
			success, err := matcher.matchLabel(lblPair)
			if err != nil {
				return false, err
			}
			if success {
				continue nextMatcher
			}
		}
		return false, nil
	}
	return true, nil
}
