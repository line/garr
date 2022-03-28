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
	"sync"
)

// MutexAdder is mutex-based LongAdder. Slowest compared to other alternatives.
type MutexAdder struct {
	value int64
	lock  sync.RWMutex
}

// NewMutexAdder creates new MutexAdder.
func NewMutexAdder() *MutexAdder {
	return &MutexAdder{}
}

// Add by the given value.
func (m *MutexAdder) Add(x int64) {
	m.lock.Lock()
	m.value += x
	m.lock.Unlock()
}

// Inc by 1.
func (m *MutexAdder) Inc() {
	m.Add(1)
}

// Dec by 1.
func (m *MutexAdder) Dec() {
	m.Add(-1)
}

// Sum returns the current sum.
func (m *MutexAdder) Sum() (sum int64) {
	m.lock.RLock()
	sum = m.value
	m.lock.RUnlock()
	return
}

// Reset variables maintaining the sum to zero.
func (m *MutexAdder) Reset() {
	m.lock.Lock()
	m.value = 0
	m.lock.Unlock()
}

// SumAndReset equivalent in effect to sum followed by reset.
func (m *MutexAdder) SumAndReset() (sum int64) {
	m.lock.Lock()
	sum = m.value
	m.value = 0
	m.lock.Unlock()
	return
}

// Store value.
func (m *MutexAdder) Store(v int64) {
	m.lock.Lock()
	m.value = v
	m.lock.Unlock()
}
