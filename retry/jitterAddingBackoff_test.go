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
	"testing"
)

func TestJitterAddingBackoff(t *testing.T) {
	if _, err := NewJitterAddingBackoff(nil, 0, 0); err == nil {
		t.FailNow()
	}

	if _, err := NewJitterAddingBackoff(NoDelayBackoff, -1.1, 1); err == nil {
		t.FailNow()
	}

	if _, err := NewJitterAddingBackoff(NoDelayBackoff, 1.1, 1); err == nil {
		t.FailNow()
	}

	if _, err := NewJitterAddingBackoff(NoDelayBackoff, 0.9, -1.1); err == nil {
		t.FailNow()
	}

	if _, err := NewJitterAddingBackoff(NoDelayBackoff, 0.9, 1.1); err == nil {
		t.FailNow()
	}

	if _, err := NewJitterAddingBackoff(NoDelayBackoff, 0.9, 0.5); err == nil {
		t.FailNow()
	}

	// fake backoff
	if b, err := NewJitterAddingBackoff(&FixedBackoff{delayMillis: -1}, 0.5, 0.9); err != nil || b == nil {
		t.FailNow()
	} else if b.NextDelayMillis(2) >= 0 {
		t.FailNow()
	}

	// real backoff
	if b, err := NewJitterAddingBackoff(&ExponentialBackoff{initialDelayMillis: 100, maxDelayMillis: 1200, multiplier: 1.2},
		0.7, 0.97); err != nil || b == nil {
		t.FailNow()
	} else {
		// fake call
		for i := 0; i < 10000; i++ {
			b.NextDelayMillis(i)
		}
	}
}

func TestJitterAddingBackoff_Stochastic(t *testing.T) {
	initialDelay := int64(100)
	minJitter := -0.02
	maxJitter := 0.03
	if b, err := NewJitterAddingBackoff(&FixedBackoff{delayMillis: initialDelay}, minJitter, maxJitter); err != nil || b == nil {
		t.FailNow()
	} else {
		histogram := make(map[int64]int)
		tryouts := 500
		expectedLowerBound := int64(float64(initialDelay) * (1 + minJitter))
		expectedUpperBound := int64(float64(initialDelay) * (1 + maxJitter))
		for tryout := 0; tryout < tryouts; tryout += 1 {
			delay := b.NextDelayMillis(1)
			if delay < expectedLowerBound || delay > expectedUpperBound {
				t.FailNow()
			}
			histogram[delay] += 1
		}
		for d := expectedLowerBound; d <= expectedUpperBound; d += 1 {
			if histogram[d] == 0 {
				t.FailNow()
			}
		}
	}
}
