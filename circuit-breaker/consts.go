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

package cbreaker

import (
	"fmt"
	"time"
)

const (
	defaultFailureRateThreshold    = 0.8
	defaultMinimumRequestThreshold = 10
	defaultTrialRequestInterval    = time.Duration(3 * time.Second)
	defaultCircuitOpenWindow       = time.Duration(10 * time.Second)
	defaultCounterSlidingWindow    = time.Duration(20 * time.Second)
	defaultCounterUpdateInterval   = time.Duration(1 * time.Second)
)

// CircuitState represents state of the circuit breaker.
type CircuitState byte

const (
	// CircuitStateClosed initial state. All requests are sent to the remote service.
	CircuitStateClosed CircuitState = 0
	// CircuitStateOpen the circuit is tripped. All requests fail immediately without calling the remote service.
	CircuitStateOpen CircuitState = 1
	// CircuitStateHalfOpen only one trial request is sent at a time until at least one request succeeds or fails.
	// If it doesn't complete within a certain time, another trial request will be sent again.
	// All other requests fails immediately same as OPEN.
	CircuitStateHalfOpen CircuitState = 2
)

var (
	// ErrTickerDurationInvalid indicates ticker duration invalid.
	ErrTickerDurationInvalid = fmt.Errorf("Ticker duration must be > 0")
)
