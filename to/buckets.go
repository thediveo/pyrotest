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

import (
	"slices"
	"sort"

	prommodel "github.com/prometheus/client_model/go"
)

// MappedBuckets maps (inclusive) upper boundaries to their corresponding bucket
// sample counts.
type MappedBuckets = map[float64]uint64

// SampledBuckets returns the (cumulative) buckets, count, and sum for the given
// lists of samples and bucket (inclusive) upper boundaries. The bucket
// boundaries list must not include “+Inf”. In the Prometheus histogram bucket
// model, the +Inf bucket is implicit in the sample sum and sample count.
func SampledBuckets(samples []float64, upperboundaries []float64) (buckets MappedBuckets, count uint64, sum float64) {
	buckets = MappedBuckets{}
	l := len(upperboundaries)
	for _, upperboundary := range upperboundaries {
		buckets[upperboundary] = 0
	}
	for _, sample := range samples {
		sum += sample
		idx := bucketIndex(sample, upperboundaries)
		if idx >= l {
			continue
		}
		buckets[upperboundaries[idx]] += 1
	}
	cumulative := uint64(0)
	for _, upperboundary := range upperboundaries {
		bucketcount := buckets[upperboundary]
		buckets[upperboundary] += cumulative
		cumulative += bucketcount
	}
	return buckets, uint64(len(samples)), sum
}

// bucketIndex returns the index of the bucket the sample falls into according
// the the list of bucket upper boundaries. If the sample falls beyond the last
// bucket boundary, bucketIndex returns len(upperboundaries).
//
// bucketIndex leverages the sort's binary search.
func bucketIndex(sample float64, upperboundaries []float64) int {
	return sort.Search(len(upperboundaries), func(i int) bool { return upperboundaries[i] >= sample })
}

// OrderedBuckets returns the list of Prometheus protobuf buckets for the
// passed-in buckets map. The returned bucket list is sorted in ascending order
// by bucket inclusive upper limits.
func OrderedBuckets(buckets MappedBuckets) []*prommodel.Bucket {
	prombuckets := make([]*prommodel.Bucket, 0, len(buckets))
	for upperbound, cumcount := range buckets {
		prombuckets = append(prombuckets, &prommodel.Bucket{
			UpperBound:      &upperbound,
			CumulativeCount: &cumcount,
		})
	}
	slices.SortFunc(prombuckets, func(a, b *prommodel.Bucket) int {
		return int(a.GetUpperBound()) - int(b.GetUpperBound())
	})
	return prombuckets
}
