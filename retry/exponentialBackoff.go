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
	"fmt"
	"math"
)

// ExponentialBackoff waits for an exponentially-increasing amount of time between attempts.
type ExponentialBackoff struct {
	initialDelayMillis int64
	maxDelayMillis     int64
	multiplier         float64
}

// NewExponentialBackoff creates new ExponentialBackoff.
func NewExponentialBackoff(initialDelayMillis, maxDelayMillis int64, multiplier float64) (b *ExponentialBackoff, err error) {
	if multiplier <= 1 {
		err = fmt.Errorf("multiplier: %.3f (expected: > 1.0)", multiplier)
	} else if initialDelayMillis < 0 {
		err = fmt.Errorf("initialDelayMillis: %d (expected: >= 0)", initialDelayMillis)
	} else if initialDelayMillis > maxDelayMillis {
		err = fmt.Errorf("maxDelayMillis: %d (expected: >= %d)", maxDelayMillis, initialDelayMillis)
	} else {
		b = &ExponentialBackoff{
			initialDelayMillis: initialDelayMillis,
			maxDelayMillis:     maxDelayMillis,
			multiplier:         multiplier,
		}
	}
	return
}

// NextDelayMillis returns the number of milliseconds to wait for before attempting a retry.
func (f *ExponentialBackoff) NextDelayMillis(numAttemptsSoFar int) (nextDelay int64) {
	if numAttemptsSoFar == 1 {
		return f.initialDelayMillis
	}

	nextDelay = saturatedMultiply(f.initialDelayMillis, math.Pow(f.multiplier, float64(numAttemptsSoFar-1)))
	if nextDelay > f.maxDelayMillis {
		nextDelay = f.maxDelayMillis
	}
	return
}
