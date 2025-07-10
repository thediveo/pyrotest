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

package to

import "sort"

// Buckets maps upper boundaries to corresponding bucket sample counts.
type Buckets = map[float64]uint64

// BucketsFromSamples returns the (cumulative) buckets, count, and sum for the
// given lists of samples and bucket (inclusive) upper boundaries. The bucket
// boundaries list must not include “+Inf”. In the Prometheus histogram bucket
// model, the +Inf bucket is implicit in the sample sum and sample count.
func BucketsFromSamples(samples []float64, upperboundaries []float64) (b Buckets, count uint64, sum float64) {
	b = Buckets{}
	l := len(upperboundaries)
	for _, upperboundary := range upperboundaries {
		b[upperboundary] = 0
	}
	for _, sample := range samples {
		sum += sample
		idx := sort.Search(l, func(i int) bool { return upperboundaries[i] >= sample })
		if idx >= l {
			continue
		}
		b[upperboundaries[idx]] += 1
	}
	cumulative := uint64(0)
	for _, upperboundary := range upperboundaries {
		bucketcount := b[upperboundary]
		b[upperboundary] += cumulative
		cumulative += bucketcount
	}
	return b, uint64(len(samples)), sum
}
