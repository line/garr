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

// CircuitBreakerConfig stores configurations of circuit breaker.
type CircuitBreakerConfig struct {
	name                    *Name
	failureRateThreshold    float64
	minimumRequestThreshold int64
	trialRequestInterval    time.Duration
	circuitOpenWindow       time.Duration
	counterSlidingWindow    time.Duration
	counterUpdateInterval   time.Duration
	listeners               CircuitBreakerListeners
}

// GetName returns name of CircuitBreaker.
func (c *CircuitBreakerConfig) GetName() *Name {
	return c.name
}

// GetFailureRateThreshold returns the threshold of failure rate to detect a remote service fault.
func (c *CircuitBreakerConfig) GetFailureRateThreshold() float64 {
	return c.failureRateThreshold
}

// GetMinimumRequestThreshold returns the minimum number of requests within a time window necessary
// to detect a remote service fault.
func (c *CircuitBreakerConfig) GetMinimumRequestThreshold() int64 {
	return c.minimumRequestThreshold
}

// GetTrialRequestInterval returns the trial request interval in HalfOpen state.
func (c *CircuitBreakerConfig) GetTrialRequestInterval() time.Duration {
	return c.trialRequestInterval
}

// GetCircuitOpenWindow returns the duration of Open state.
func (c *CircuitBreakerConfig) GetCircuitOpenWindow() time.Duration {
	return c.circuitOpenWindow
}

// GetCounterSlidingWindow returns the time length of sliding window to accumulate the count of events.
func (c *CircuitBreakerConfig) GetCounterSlidingWindow() time.Duration {
	return c.counterSlidingWindow
}

// GetCounterUpdateInterval returns the interval that a circuit breaker can see the latest accumulated count of events.
func (c *CircuitBreakerConfig) GetCounterUpdateInterval() time.Duration {
	return c.counterUpdateInterval
}

// Getlisteners returns CircuitBreakerListener(s)
func (c *CircuitBreakerConfig) Getlisteners() CircuitBreakerListeners {
	return c.listeners
}

// Validate current configuration.
func (c *CircuitBreakerConfig) Validate() (err error) {
	if c.failureRateThreshold <= 0 || 1 < c.failureRateThreshold {
		err = fmt.Errorf("failureRateThreshold: %.3f (expected: > 0 and <= 1)", c.failureRateThreshold)
		return
	}

	if c.trialRequestInterval <= 0 {
		err = fmt.Errorf("trialRequestInterval: %d (expected: > 0)", c.trialRequestInterval)
		return
	}

	if c.circuitOpenWindow <= 0 {
		err = fmt.Errorf("circuitOpenWindow: %d (expected: > 0)", c.circuitOpenWindow)
		return
	}

	if c.counterSlidingWindow <= 0 {
		err = fmt.Errorf("counterSlidingWindow: %d (expected: > 0)", c.counterSlidingWindow)
		return
	}

	if c.counterUpdateInterval <= 0 {
		err = fmt.Errorf("counterUpdateInterval: %d (expected: > 0)", c.counterUpdateInterval)
		return
	}

	if c.counterSlidingWindow <= c.counterUpdateInterval {
		err = fmt.Errorf("counterSlidingWindow: %d (expected: > counterUpdateInterval)", c.counterSlidingWindow)
		return
	}

	return
}

// String is stringer interface of CircuitBreakerConfig.
func (c *CircuitBreakerConfig) String() string {
	return fmt.Sprintf("name: %s, failureRateThreshold: %.3f, minimumRequestThreshold: %d, trialRequestInterval: %d, circuitOpenWindow: %d, counterSlidingWindow: %d, counterUpdateInterval: %d",
		c.name, c.failureRateThreshold, c.minimumRequestThreshold,
		c.trialRequestInterval, c.circuitOpenWindow,
		c.counterSlidingWindow, c.counterUpdateInterval,
	)
}
