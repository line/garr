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

// EventCounter event counter interface.
type EventCounter interface {
	// Count return the current EventCount.
	Count() *EventCount
	// OnSuccess count success events.
	OnSuccess() *EventCount
	// OnFailure count failure events.
	OnFailure() *EventCount
}

// EventCount stores the count of events.
type EventCount struct {
	success int64
	failure int64
}

// NewEventCount creates new event count.
func NewEventCount(success, failure int64) *EventCount {
	return &EventCount{success: success, failure: failure}
}

// Success returns number of success events.
func (e *EventCount) Success() int64 {
	return e.success
}

// Failure returns number of failure events.
func (e *EventCount) Failure() int64 {
	return e.failure
}

// Total returns total number of events.
func (e *EventCount) Total() int64 {
	return e.success + e.failure
}

// SuccessRate return number of success rate
func (e *EventCount) SuccessRate() float64 {
	total := e.Total()
	if total == 0 {
		return -1
	}
	return float64(e.success) / float64(total)
}

// FailureRate return number of failure rate
func (e *EventCount) FailureRate() float64 {
	total := e.Total()
	if total == 0 {
		return -1
	}
	return float64(e.failure) / float64(total)
}

// EventCountZero event count with zero
var EventCountZero = NewEventCount(0, 0)
