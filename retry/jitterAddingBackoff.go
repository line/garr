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

// JitterAddingBackoff returns a Backoff that adds a random jitter value to the original delay using
// https://www.awsarchitectureblog.com/2015/03/backoff.html full jitter strategy.
type JitterAddingBackoff struct {
	minJitterRate float64
	maxJitterRate float64
	delegate      Backoff
}

// NewJitterAddingBackoff creates new JitterAddingBackoff.
func NewJitterAddingBackoff(delegate Backoff, minJitterRate, maxJitterRate float64) (b *JitterAddingBackoff, err error) {
	if delegate == nil {
		err = fmt.Errorf("Delegate must be not nil")
	} else if !(-1.0 <= minJitterRate && minJitterRate <= 1.0) {
		err = fmt.Errorf("minJitterRate: %.3f (expected: >= -1.0 and <= 1.0)", minJitterRate)
	} else if !(-1.0 <= maxJitterRate && maxJitterRate <= 1.0) {
		err = fmt.Errorf("maxJitterRate: %.3f (expected: >= -1.0 and <= 1.0)", maxJitterRate)
	} else if minJitterRate > maxJitterRate {
		err = fmt.Errorf("maxJitterRate: %.3f needs to be greater than or equal to minJitterRate: %.3f", maxJitterRate, minJitterRate)
	} else {
		b = &JitterAddingBackoff{minJitterRate: minJitterRate, maxJitterRate: maxJitterRate, delegate: delegate}
	}
	return
}

// NextDelayMillis returns the number of milliseconds to wait for before attempting a retry.
func (f *JitterAddingBackoff) NextDelayMillis(numAttemptsSoFar int) (nextDelay int64) {
	tmp := f.delegate.NextDelayMillis(numAttemptsSoFar)
	if tmp <= 0 {
		return tmp
	}

	minJitter := int64(float64(tmp) * (1 + f.minJitterRate))
	maxJitter := int64(float64(tmp) * (1 + f.maxJitterRate))
	if nextDelay = minJitter + nextRandomInt64IncludingZero(maxJitter-minJitter+1); nextDelay < 0 {
		nextDelay = 0
	}
	return
}
