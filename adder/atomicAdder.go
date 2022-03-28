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

package adder

import (
	"sync/atomic"
)

// AtomicAdder is simple atomic-based adder.
// Fastest at single routine but slower at multi routine when high-contention happens.
type AtomicAdder struct {
	value int64
}

// NewAtomicAdder creates new AtomicAdder.
func NewAtomicAdder() *AtomicAdder {
	return &AtomicAdder{}
}

// Add by the given value.
func (a *AtomicAdder) Add(x int64) {
	atomic.AddInt64(&a.value, x)
}

// Inc by 1.
func (a *AtomicAdder) Inc() {
	a.Add(1)
}

// Dec by 1.
func (a *AtomicAdder) Dec() {
	a.Add(-1)
}

// Sum returns the current sum. The returned value is NOT an
// atomic snapshot because of concurrent update.
func (a *AtomicAdder) Sum() int64 {
	return atomic.LoadInt64(&a.value)
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
func (a *AtomicAdder) Reset() {
	atomic.StoreInt64(&a.value, 0)
}

// SumAndReset equivalent in effect to sum followed by reset. Like the nature of Sum and Reset,
// this function is only effective if there are no concurrent updates.
func (a *AtomicAdder) SumAndReset() (sum int64) {
	sum = atomic.LoadInt64(&a.value)
	a.Reset()
	return
}

// Store value. This function is only effective if there are no concurrent updates.
func (a *AtomicAdder) Store(v int64) {
	atomic.StoreInt64(&a.value, v)
}
