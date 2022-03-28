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
	"math"
	"sync/atomic"
)

// AtomicF64Adder is simple atomic-based adder.
// Fastest at single routine but slower at multi routine when high-contention happens.
type AtomicF64Adder struct {
	value uint64
}

// NewAtomicF64Adder creates new AtomicF64Adder.
func NewAtomicF64Adder() *AtomicF64Adder {
	return &AtomicF64Adder{}
}

// Add by the given value.
func (a *AtomicF64Adder) Add(v float64) {
	for {
		old := a.Sum()
		new := old + v
		if atomic.CompareAndSwapUint64(&a.value, math.Float64bits(old), math.Float64bits(new)) {
			return
		}
	}
}

// Inc by 1.
func (a *AtomicF64Adder) Inc() {
	a.Add(1)
}

// Dec by 1.
func (a *AtomicF64Adder) Dec() {
	a.Add(-1)
}

// Sum returns the current sum. The returned value is NOT an
// atomic snapshot because of concurrent update.
func (a *AtomicF64Adder) Sum() float64 {
	return math.Float64frombits(atomic.LoadUint64(&a.value))
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
func (a *AtomicF64Adder) Reset() {
	atomic.StoreUint64(&a.value, 0)
}

// SumAndReset equivalent in effect to sum followed by reset. Like the nature of Sum and Reset,
// this function is only effective if there are no concurrent updates.
func (a *AtomicF64Adder) SumAndReset() (sum float64) {
	sum = a.Sum()
	a.Reset()
	return
}

// Store value. This function is only effective if there are no concurrent updates.
func (a *AtomicF64Adder) Store(v float64) {
	atomic.StoreUint64(&a.value, math.Float64bits(v))
}
