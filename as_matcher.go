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
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

// asStringMatcher expects a to be either a string or a types.GomegaMatcher and
// then always returns a suitable types.GomegaMatcher, otherwise nil in case of
// an unsupported value type of a.
func asStringMatcher(a any) types.GomegaMatcher {
	switch v := a.(type) {
	case string:
		return gomega.Equal(v)
	case types.GomegaMatcher:
		return v
	default:
		return nil
	}
}

// asMatcher returns a types.Matcher to match the passed-in value, or the
// passed-in Gomega matcher, or nil.
func asMatcher(a any) types.GomegaMatcher {
	if a == nil {
		return nil
	}
	if gm, ok := a.(types.GomegaMatcher); ok {
		return gm
	}
	return gomega.Equal(a)
}
