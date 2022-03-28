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
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestNoopCounter(t *testing.T) {
	noop := &noopCounter{}
	if noop.Count() != EventCountZero {
		t.FailNow()
	}
	if noop.OnSuccess() != nil || noop.OnFailure() != nil {
		t.FailNow()
	}
}

func TestNewNonBlockingCircuitBreaker(t *testing.T) {
	if _, err := NewNonBlockingCircuitBreaker(nil, nil); err == nil {
		t.Errorf("Invalid NewNonBlockingCircuitBreaker")
	} else if _, err = NewNonBlockingCircuitBreaker(nil, &CircuitBreakerConfig{}); err == nil {
		t.Errorf("Invalid NewNonBlockingCircuitBreaker")
	} else if _, err = NewNonBlockingCircuitBreaker(SystemTicker, nil); err == nil {
		t.Errorf("Invalid NewNonBlockingCircuitBreaker")
	} else if _, err = NewNonBlockingCircuitBreaker(SystemTicker, &CircuitBreakerConfig{}); err == nil {
		t.Errorf("Invalid NewNonBlockingCircuitBreaker")
	}
}

func TestNewNonBlockingCircuitBreakerWithValidConfigs(t *testing.T) {
	name := &Name{Name: "dummy-breaker"}
	validConfig := &CircuitBreakerConfig{
		name:                    name,
		failureRateThreshold:    0.7,
		minimumRequestThreshold: 19,
		trialRequestInterval:    time.Second,
		circuitOpenWindow:       time.Second * 2,
		counterSlidingWindow:    time.Second * 30,
		counterUpdateInterval:   time.Second * 4,
		listeners:               make(CircuitBreakerListeners, 2, 10),
	}

	nbc, err := NewNonBlockingCircuitBreaker(SystemTicker, validConfig)
	if err != nil || nbc == nil || nbc.config != validConfig || nbc.Name() != name || nbc.ticker != SystemTicker {
		if err != nil {
			t.Error(err)
		} else {
			t.Errorf("Invalid NewNonBlockingCircuitBreaker")
		}
	}

	if validConfig.GetName() != name ||
		validConfig.GetFailureRateThreshold() != 0.7 ||
		validConfig.GetMinimumRequestThreshold() != 19 ||
		validConfig.GetTrialRequestInterval() != time.Second ||
		validConfig.GetCircuitOpenWindow() != 2*time.Second ||
		validConfig.GetCounterSlidingWindow() != 30*time.Second ||
		validConfig.GetCounterUpdateInterval() != 4*time.Second ||
		len(validConfig.Getlisteners()) != 2 {
		t.Errorf("Invalid NewNonBlockingCircuitBreaker")
	}
}

type cbListenerMock struct{}

// OnStateChanged invoked when the circuit state is changed.
func (c *cbListenerMock) OnStateChanged(cb CircuitBreaker, state CircuitState) (err error) {
	return fmt.Errorf("Fake error")
}

// OnEventCountUpdated invoked when the circuit breaker's internal EventCount is updated.
func (c *cbListenerMock) OnEventCountUpdated(cb CircuitBreaker, eventCount *EventCount) (err error) {
	return fmt.Errorf("Fake error")
}

// OnRequestRejected invoked when the circuit breaker rejects a request.
func (c *cbListenerMock) OnRequestRejected(cb CircuitBreaker) (err error) {
	return fmt.Errorf("Fake error")
}

func (c *cbListenerMock) Stop() {}

func TestNonBlockingCircuitBreaker_MoreSuccess(t *testing.T) {
	validConfig := &CircuitBreakerConfig{
		name:                    &Name{Name: "dummy-breaker"},
		failureRateThreshold:    0.3,
		minimumRequestThreshold: 19,
		trialRequestInterval:    time.Second,
		circuitOpenWindow:       time.Second * 2,
		counterSlidingWindow:    time.Second * 10,
		counterUpdateInterval:   time.Millisecond * 50,
		listeners:               make(CircuitBreakerListeners, 3, 10),
	}
	validConfig.listeners[2] = &cbListenerMock{}

	nbc, _ := NewNonBlockingCircuitBreaker(SystemTicker, validConfig)

	// set start state is opened
	nbc.s = nbc.newOpenState()

	// fake notify
	nbc.notifyCountUpdated(&EventCount{})
	nbc.notifyRequestRejected()
	nbc.notifyStateChanged(CircuitStateOpen)

	numberWorker := 100
	var wg sync.WaitGroup
	for i := 0; i < numberWorker; i++ {
		wg.Add(1)
		go func() {
			var count int64
			for j := 0; j < 100000; j++ {
				if _, err := nbc.Execute(context.Background(), func(context.Context) (result interface{}, err error) {
					if x := rand.Int() % 100000; x > 2000 {
						err = nil
					} else {
						err = fmt.Errorf("Fake error")
					}
					return
				}); err != nil {
					count++
					nbc.OnFailure()
					time.Sleep(time.Millisecond >> 4)
				} else if nbc.CanRequest() {
					count++
					nbc.OnSuccess()
					time.Sleep(time.Millisecond >> 6)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestNonBlockingCircuitBreaker_Failure(t *testing.T) {
	validConfig := &CircuitBreakerConfig{
		name:                    &Name{Name: "dummy-breaker"},
		failureRateThreshold:    0.3,
		minimumRequestThreshold: 19,
		trialRequestInterval:    time.Second,
		circuitOpenWindow:       time.Second * 2,
		counterSlidingWindow:    time.Minute * 10,
		counterUpdateInterval:   time.Second * 5,
		listeners:               make(CircuitBreakerListeners, 2, 10),
	}
	nbc, _ := NewNonBlockingCircuitBreaker(SystemTicker, validConfig)

	// set start state is closed
	nbc.s = nbc.newClosedState()

	numberWorker := 50
	var wg sync.WaitGroup
	for i := 0; i < numberWorker; i++ {
		wg.Add(1)
		go func() {
			var count int64
			for j := 0; j < 100000; j++ {
				if _, err := nbc.Execute(context.Background(), func(context.Context) (result interface{}, err error) {
					if x := rand.Int() % 100000; x < 2000 {
						err = nil
					} else {
						err = fmt.Errorf("Fake error")
					}
					return
				}); err == nil {
					count++
					nbc.OnSuccess()
					time.Sleep(time.Millisecond >> 8)
				} else if nbc.CanRequest() {
					count++
					nbc.OnFailure()
					time.Sleep(time.Millisecond >> 4)
				}

				// fake call to counter
				nbc.state().counter.Count()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
