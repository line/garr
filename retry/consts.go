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

import "fmt"

const (
	// DefaultDelayMillis is default delay millis.
	DefaultDelayMillis int64 = 200
	// DefaultInitialDelayMillis is default initial delay millis.
	DefaultInitialDelayMillis int64 = 200
	// DefaultMinDelayMillis is default min delay millis.
	DefaultMinDelayMillis int64 = 0
	// DefaultMaxDelayMillis is default max delay millis.
	DefaultMaxDelayMillis int64 = 10000
	// DefaultMultiplier is default multiplier.
	DefaultMultiplier float64 = 2.0
	// DefaultMinJitterRate is default min jitter rate.
	DefaultMinJitterRate float64 = -0.2
	// DefaultMaxJitterRate is default max jitter rate.
	DefaultMaxJitterRate float64 = 0.2
)

var (
	// ErrInvalidSpecFormat indicates invalid specification format.
	ErrInvalidSpecFormat = fmt.Errorf("Invalid format of specification")
)
