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
)

// RandomBackoff computes backoff delay which is a random value between
// minDelayMillis} and maxDelayMillis.
type RandomBackoff struct {
	minDelayMillis int64
	maxDelayMillis int64
	bound          int64
}

// NewRandomBackoff creates new RandomBackoff.
func NewRandomBackoff(minDelayMillis, maxDelayMillis int64) (b *RandomBackoff, err error) {
	if minDelayMillis < 0 {
		err = fmt.Errorf("minDelayMillis: %d (expected: >= 0)", minDelayMillis)
	} else if minDelayMillis > maxDelayMillis {
		err = fmt.Errorf("maxDelayMillis: %d (expected: >= %d)", maxDelayMillis, minDelayMillis)
	} else {
		b = &RandomBackoff{minDelayMillis: minDelayMillis, maxDelayMillis: maxDelayMillis, bound: maxDelayMillis - minDelayMillis}
	}
	return
}

// NextDelayMillis returns number of milliseconds to wait for before attempting a retry.
func (f *RandomBackoff) NextDelayMillis(numAttemptsSoFar int) int64 {
	if f.minDelayMillis != f.maxDelayMillis {
		return nextRandomInt64(f.bound) + f.minDelayMillis
	}
	return f.minDelayMillis
}
