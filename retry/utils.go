// Copyright 2022 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package retry

import (
	"math"

	"github.com/valyala/fastrand"
)

const (
	limit64 = (1 << 63) - 1
)

func randomInt64() (result int64) {
	result |= (int64(fastrand.Uint32()) << 32) & limit64
	result |= int64(fastrand.Uint32())
	return
}

func saturatedMultiply(left int64, right float64) int64 {
	if tmp := float64(left) * right; tmp < math.MaxInt64 {
		return int64(tmp)
	}
	return math.MaxInt64
}

// generates a random number in range [1, bound].
// If the given bound is not positive, fast return the bound.
func nextRandomInt64(bound int64) int64 {
	if bound <= 0 {
		return bound
	}
	return nextRandomInt64IncludingZero(bound-1) + 1
}

// generates a random number in range [0, bound].
// If the given bound is not positive, fast return the bound.
func nextRandomInt64IncludingZero(bound int64) (result int64) {
	if bound <= 0 {
		return bound
	}

	mask := bound - 1
	result = randomInt64()

	if bound&mask == 0 {
		result &= mask
	} else {
		u := result >> 1
		for {
			if result = u % bound; u < result-mask {
				u = randomInt64() >> 1
			} else {
				break
			}
		}
	}

	return
}
