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
	prommodel "github.com/prometheus/client_model/go"
)

type BeAMetricMatcher struct {
	Expected MetricMatcher
}

var _ types.GomegaMatcher = (*BeAMetricMatcher)(nil)

func (m *BeAMetricMatcher) Match(actual any) (bool, error) {
	mf, ok := actual.(*prommodel.MetricFamily)
	if !ok {
		return false, fmt.Errorf("BeAMetricMatcher expects a Prometheus MetricFamily.  Got:\n%s",
			format.Object(actual, 1))
	}
	return m.Expected.match(mf)
}

func (m *BeAMetricMatcher) FailureMessage(actual any) string {
	return format.Message(actual, "to match", m.Expected)
}

func (m *BeAMetricMatcher) NegatedFailureMessage(actual any) string {
	return format.Message(actual, "not to match", m.Expected)
}
