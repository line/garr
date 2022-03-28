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
	"testing"
)

func TestExponentialBackoff(t *testing.T) {
	backoff, _ := NewExponentialBackoff(100, 2000, 1.3)
	for i := 1; i < 100; i++ {
		backoff.NextDelayMillis(i)
	}

	if _, err := NewExponentialBackoff(-1, 12, 3); err == nil {
		t.FailNow()
	}

	if _, err := NewExponentialBackoff(3, 2, 3); err == nil {
		t.FailNow()
	}

	if _, err := NewExponentialBackoff(3, 12, 0.3); err == nil {
		t.FailNow()
	}

	// fake
	if saturatedMultiply(3, float64(math.MaxInt64)) != math.MaxInt64 {
		t.FailNow()
	}
}
