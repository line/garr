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

// AttemptLimitingBackoff is a backoff which limits the number of attempts up to the specified value.
type AttemptLimitingBackoff struct {
	delegate Backoff
	limit    int
}

// NewAttemptLimitingBackoff creates new AttemptLimitingBackoff.
func NewAttemptLimitingBackoff(delegate Backoff, limit int) (b *AttemptLimitingBackoff, err error) {
	if delegate == nil {
		err = fmt.Errorf("Delegate must be not nil")
	} else if limit <= 0 {
		err = fmt.Errorf("maxAttempts: %d (expected: > 0)", limit)
	} else {
		b = &AttemptLimitingBackoff{delegate: delegate, limit: limit}
	}
	return
}

// NextDelayMillis returns the number of milliseconds to wait for before attempting a retry.
func (f *AttemptLimitingBackoff) NextDelayMillis(numAttemptsSoFar int) int64 {
	if numAttemptsSoFar >= f.limit {
		return -1
	}
	return f.delegate.NextDelayMillis(numAttemptsSoFar)
}
