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
	"time"
)

// CircuitBreakerBuilder is the builder for CircuitBreaker, using builder-pattern.
type CircuitBreakerBuilder struct {
	name                    *Name
	ticker                  Ticker
	failureRateThreshold    float64
	minimumRequestThreshold int64
	trialRequestInterval    time.Duration
	circuitOpenWindow       time.Duration
	counterSlidingWindow    time.Duration
	counterUpdateInterval   time.Duration
	listeners               CircuitBreakerListeners
}

// NewCircuitBreakerBuilder creates new circuit breaker builder.
func NewCircuitBreakerBuilder() (c *CircuitBreakerBuilder) {
	c = &CircuitBreakerBuilder{
		ticker:                  SystemTicker,
		failureRateThreshold:    defaultFailureRateThreshold,
		minimumRequestThreshold: defaultMinimumRequestThreshold,
		trialRequestInterval:    defaultTrialRequestInterval,
		circuitOpenWindow:       defaultCircuitOpenWindow,
		counterSlidingWindow:    defaultCounterSlidingWindow,
		counterUpdateInterval:   defaultCounterUpdateInterval,
	}
	return
}

// Name sets name for circuit breaker.
func (c *CircuitBreakerBuilder) Name(name *Name) *CircuitBreakerBuilder {
	c.name = name
	return c
}

// SetTicker sets ticker.
func (c *CircuitBreakerBuilder) SetTicker(t Ticker) *CircuitBreakerBuilder {
	c.ticker = t
	return c
}

// SetFailureRateThreshold sets the threshold of failure rate to detect a remote service fault.
func (c *CircuitBreakerBuilder) SetFailureRateThreshold(failureRateThreshold float64) *CircuitBreakerBuilder {
	c.failureRateThreshold = failureRateThreshold
	return c
}

// SetMinimumRequestThreshold sets the minimum number of requests within a time window necessary to detect a remote service fault.
func (c *CircuitBreakerBuilder) SetMinimumRequestThreshold(minimumRequestThreshold int64) *CircuitBreakerBuilder {
	c.minimumRequestThreshold = minimumRequestThreshold
	return c
}

// SetTrialRequestInterval sets the trial request interval in HalfOpen state.
func (c *CircuitBreakerBuilder) SetTrialRequestInterval(trialRequestInterval time.Duration) *CircuitBreakerBuilder {
	c.trialRequestInterval = trialRequestInterval
	return c
}

// SetCircuitOpenWindow sets the duration of Open state.
func (c *CircuitBreakerBuilder) SetCircuitOpenWindow(circuitOpenWindow time.Duration) *CircuitBreakerBuilder {
	c.circuitOpenWindow = circuitOpenWindow
	return c
}

// SetCounterSlidingWindow sets the time length of sliding window to accumulate the count of events.
func (c *CircuitBreakerBuilder) SetCounterSlidingWindow(counterSlidingWindow time.Duration) *CircuitBreakerBuilder {
	c.counterSlidingWindow = counterSlidingWindow
	return c
}

// SetCounterUpdateInterval sets the interval that a circuit breaker can see the latest accumulated count of events.
func (c *CircuitBreakerBuilder) SetCounterUpdateInterval(counterUpdateInterval time.Duration) *CircuitBreakerBuilder {
	c.counterUpdateInterval = counterUpdateInterval
	return c
}

// AddListener adds a CircuitBreakerListener.
func (c *CircuitBreakerBuilder) AddListener(listener CircuitBreakerListener) *CircuitBreakerBuilder {
	if listener != nil {
		c.listeners = append(c.listeners, listener)
	}
	return c
}

// Build returns a newly-created CircuitBreaker based on the properties of this builder.
func (c *CircuitBreakerBuilder) Build() (cb CircuitBreaker, err error) {
	cb, err = NewNonBlockingCircuitBreaker(c.ticker, &CircuitBreakerConfig{
		name:                    c.name,
		failureRateThreshold:    c.failureRateThreshold,
		minimumRequestThreshold: c.minimumRequestThreshold,
		trialRequestInterval:    c.trialRequestInterval,
		circuitOpenWindow:       c.circuitOpenWindow,
		counterSlidingWindow:    c.counterSlidingWindow,
		counterUpdateInterval:   c.counterUpdateInterval,
		listeners:               c.listeners,
	})
	return
}
