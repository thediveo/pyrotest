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
	"github.com/onsi/gomega/format"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("family property matchers", func() {

	DescribeTable("GomegaString",
		func(m metricFamilyPropertyMatcher, expected string) {
			gs, ok := m.(format.GomegaStringer)
			Expect(ok).To(BeTrue(), "not a GomegaStringer: %T", m)
			Expect(gs.GomegaString()).To(MatchRegexp(expected))
		},
		Entry("name property", HaveName("foobar"),
			`name: foobar`),
		Entry("name property/matcher", HaveName(Equal("foobar")),
			`name: .*EqualMatcher.*\n.*Expected: <string>"foobar"`),

		Entry("help property", HaveHelp("foobar"),
			`help: foobar`),
		Entry("help property/matcher", HaveHelp(Equal("foobar")),
			`help: .*EqualMatcher.*\n.*Expected: <string>"foobar"`),

		Entry("unit property", HaveUnit("foobar"),
			`unit: foobar`),
		Entry("unit property/matcher", HaveUnit(Equal("foobar")),
			`unit: .*EqualMatcher.*\n.*Expected: <string>"foobar"`),
	)

})
