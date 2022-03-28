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

// Package cbreaker (circuit-breaker) contains circuit breakers ported from line/armeria.
package cbreaker

import (
	"context"
	"fmt"
)

var (
	// ErrFailFast indicates that fail-fast is detected.
	ErrFailFast = fmt.Errorf("FailFast detected")
)

// Name fully-qualified name.
type Name struct {
	Namespace string
	Subsystem string
	Name      string
}

// Execute function.
type Execute func(ctx context.Context) (result interface{}, err error)

// CircuitBreaker tracks the number of success/failure requests and detects a remote service failure.
type CircuitBreaker interface {
	// Name return the name of the circuit breaker.
	Name() *Name
	// OnSuccess report a remote invocation success.
	OnSuccess()
	// OnFailure report a remote invocation failure.
	OnFailure()
	// CanRequest decide whether a request should be sent or failed depending on the current circuit state.
	CanRequest() bool
	// Execute delegated function.
	Execute(ctx context.Context, delegatedFn Execute) (r interface{}, err error)
}

// CircuitBreakerListener is listener interface for receiving events.
type CircuitBreakerListener interface {
	// OnStateChanged invoked when the circuit state is changed.
	OnStateChanged(cb CircuitBreaker, state CircuitState) (err error)
	// OnEventCountUpdated invoked when the circuit breaker's internal EventCount is updated.
	OnEventCountUpdated(cb CircuitBreaker, eventCount *EventCount) (err error)
	// OnRequestRejected invoked when the circuit breaker rejects a request.
	OnRequestRejected(cb CircuitBreaker) (err error)
	// Stop notify listener to stop.
	Stop()
}

// CircuitBreakerListeners is collection of CircuitBreakerListener.
type CircuitBreakerListeners []CircuitBreakerListener
