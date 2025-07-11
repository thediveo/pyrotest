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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("conversions", func() {

	Context("histogram buckets", func() {

		It("converts samples into cumulative histogram buckets ", func() {
			ubs := []float64{1.0, 2.0, 4.0}
			samples := []float64{0.0, 0.5, 1.0, 1.5, 3.2, 6.66}
			buckets, count, sum := SampledBuckets(samples, ubs)

			Expect(count).To(Equal(uint64(len(samples))))

			expectedsum := 0.0
			for _, sample := range samples {
				expectedsum += sample
			}
			Expect(sum).To(Equal(expectedsum))

			Expect(buckets).To(HaveLen(len(ubs)))
			Expect(buckets).To(Equal(MappedBuckets{
				1.0: 3,
				2.0: 3 + 1,
				4.0: 3 + 1 + 1,
			}))
		})

		It("converts a histogram bucket map into a sorted bucket+boundary list", func() {
			ubs := []float64{1.0, 2.0, 4.0}
			samples := []float64{0.0, 0.5, 1.0, 1.5, 3.2, 6.66}
			buckets, _, _ := SampledBuckets(samples, ubs)
			bucketlist := OrderedBuckets(buckets)
			Expect(bucketlist).To(HaveLen(len(ubs)))
			for idx := range bucketlist {
				if idx == 0 {
					continue
				}
				Expect(bucketlist[idx-1].GetUpperBound()).To(
					BeNumerically("<", bucketlist[idx].GetUpperBound()))
			}
			for _, bucket := range bucketlist {
				Expect(buckets).To(HaveKeyWithValue(
					bucket.GetUpperBound(), bucket.GetCumulativeCount()))
			}
		})

	})

})
