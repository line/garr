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

// FixedBackoff waits for a fixed delay between attempts.
type FixedBackoff struct {
	delayMillis int64
}

// NewFixedBackoff creates new fixed backoff.
func NewFixedBackoff(delayMillis int64) (b *FixedBackoff, err error) {
	if delayMillis >= 0 {
		b = &FixedBackoff{delayMillis: delayMillis}
	} else {
		err = fmt.Errorf("delayMillis: %d (expected: >= 0)", delayMillis)
	}
	return
}

// NextDelayMillis returns the number of milliseconds to wait for before attempting a retry.
func (f *FixedBackoff) NextDelayMillis(numAttemptsSoFar int) int64 {
	return f.delayMillis
}

// NoDelayBackoff returns a Backoff that will never wait between attempts.
// In most cases, using Backoff without delay is very dangerous.
var NoDelayBackoff Backoff = &FixedBackoff{delayMillis: 0}

// NoRetry returns a Backoff indicates that no retry.
var NoRetry Backoff = &FixedBackoff{delayMillis: -1}
