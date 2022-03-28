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

const (
	randomCellSize      = 1 << 7 // 128
	randomCellSizeMinus = randomCellSize - 1
)

// RandomCellAdder takes idea from JDKAdder by preallocating a fixed number of Cells. Unlike JDKAdder, in each update,
// RandomCellAdder assign a random-fixed Cell to invoker instead of retry/reassign Cell when contention.
//
// RandomCellAdder is often faster than JDKAdder in multi routine race benchmark
// but slower in case of single routine (no race).
//
// RandomCellAdder consume ~1KB for storing cells, which is often larger than JDKAdder which number of cells is dynamic.
type RandomCellAdder struct {
	cells []int64
}

// NewRandomCellAdder create new RandomCellAdder
func NewRandomCellAdder() *RandomCellAdder {
	return &RandomCellAdder{
		cells: make([]int64, randomCellSize),
	}
}

// Add the given value
func (r *RandomCellAdder) Add(x int64) {
	atomic.AddInt64(&r.cells[getRandomInt()&randomCellSizeMinus], x)
}

// Inc by 1
func (r *RandomCellAdder) Inc() {
	r.Add(1)
}

// Dec by 1
func (r *RandomCellAdder) Dec() {
	r.Add(-1)
}

// Sum return the current sum. The returned value is NOT an
// atomic snapshot; invocation in the absence of concurrent
// updates returns an accurate result, but concurrent updates that
// occur while the sum is being calculated might not be
// incorporated.
func (r *RandomCellAdder) Sum() (sum int64) {
	for i := range r.cells {
		sum += atomic.LoadInt64(&r.cells[i])
	}
	return
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
// Because this method is intrinsically racy
func (r *RandomCellAdder) Reset() {
	for i := range r.cells {
		atomic.StoreInt64(&r.cells[i], 0)
	}
}

// SumAndReset equivalent in effect to sum followed by reset.
// This method may apply for example during quiescent
// points between multithreaded computations. If there are
// updates concurrent with this method, the returned value is
// guaranteed to be the final value occurring before
// the reset.
func (r *RandomCellAdder) SumAndReset() (sum int64) {
	for i := range r.cells {
		sum += atomic.LoadInt64(&r.cells[i])
		atomic.StoreInt64(&r.cells[i], 0)
	}
	return
}

// Store value. This function is only effective if there are no concurrent updates.
func (r *RandomCellAdder) Store(v int64) {
	atomic.StoreInt64(&r.cells[0], v)
	for i := 1; i < randomCellSize; i++ {
		atomic.StoreInt64(&r.cells[i], 0)
	}
}
