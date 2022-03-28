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
	"runtime"
	"sync/atomic"

	"go.linecorp.com/garr/internal"
)

var maxCells int

func init() {
	maxCells = normalizeMaxCell(runtime.NumCPU() << 2)
}

func normalizeMaxCell(v int) int {
	switch {
	case v > (1 << 11):
		return (1 << 11)
	case v < 64:
		return 64
	default:
		return v
	}
}

type cells []atomic.Value

type cell struct {
	_   internal.CacheLinePad
	val int64
	_   internal.CacheLinePad
}

func (c *cell) cas(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.val, old, new)
}

// striped64 is ported version of OpenJDK9 striped64.
// It maintains a lazily-initialized table of atomically
// updated variables, plus an extra "base" field. The table size
// is a power of two. Indexing uses masked per-routine hash codes.
// Nearly all declarations in this class are package-private,
// accessed directly by subclasses.
type striped64 struct {
	cells     atomic.Value
	cellsBusy int32
	base      int64
}

func (s *striped64) casBase(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&s.base, old, new)
}

func (s *striped64) casCellsBusy() bool {
	return atomic.CompareAndSwapInt32(&s.cellsBusy, 0, 1)
}

func (s *striped64) accumulate(index int, x int64, fn longBinaryOperator, wasUncontended bool) {
	if index == 0 {
		index = getRandomInt()
		wasUncontended = true
	}

	collide := false
	var v, newV int64
	var as, rs cells
	var a, r *cell
	var m, n, j int

	var _a, _as interface{}

	for {
		_as = s.cells.Load()
		if _as != nil {
			as = _as.(cells)

			n = len(as) - 1
			if n < 0 {
				goto checkCells
			}

			if _a = as[index&n].Load(); _a != nil {
				a = _a.(*cell)
			} else {
				a = nil
			}

			if a == nil {
				if atomic.LoadInt32(&s.cellsBusy) == 0 { // Try to attach new Cell
					r = &cell{val: x} // Optimistically create
					if atomic.LoadInt32(&s.cellsBusy) == 0 && s.casCellsBusy() {
						rs = s.cells.Load().(cells)
						if m = len(rs) - 1; rs != nil && m >= 0 { // Recheck under lock
							if j = index & m; rs[j].Load() == nil {
								rs[j].Store(r)
								atomic.StoreInt32(&s.cellsBusy, 0)
								break
							}
						}
						atomic.StoreInt32(&s.cellsBusy, 0)
						continue
					}
				}
				collide = false
			} else if !wasUncontended { // CAS already known to fail
				wasUncontended = true // Continue after rehash
			} else {
				index &= n
				if v = atomic.LoadInt64(&a.val); fn == nil {
					newV = v + x
				} else {
					newV = fn.Apply(v, x)
				}
				if a.cas(v, newV) {
					break
				} else if n >= maxCells || &as[0] != &s.cells.Load().(cells)[0] { // At max size or stale
					collide = false
				} else if !collide {
					collide = true
				} else if atomic.LoadInt32(&s.cellsBusy) == 0 && s.casCellsBusy() {
					rs = s.cells.Load().(cells)
					if &as[0] == &rs[0] { // double size of cells
						if n = cap(as); len(as) < n {
							s.cells.Store(rs[:n])
						} else {
							// slice is full, n == len(as) then we just x4 size for buffering
							// Note: this trick is different from jdk source code
							rs = make(cells, n<<1, n<<2)
							copy(rs, as)
							s.cells.Store(rs)
						}
					}
					atomic.StoreInt32(&s.cellsBusy, 0)
					collide = false
					continue
				}
			}

			index ^= index << 13 // xorshift
			index ^= index >> 17
			index ^= index << 5
			continue
		}

	checkCells:
		if _as == nil {
			if atomic.LoadInt32(&s.cellsBusy) == 0 && s.cells.Load() == nil && s.casCellsBusy() {
				if s.cells.Load() == nil { // Initialize table
					rs = make(cells, 2, 4)
					rs[index&1].Store(&cell{val: x})
					s.cells.Store(rs)
					atomic.StoreInt32(&s.cellsBusy, 0)
					break
				}
				atomic.StoreInt32(&s.cellsBusy, 0)
			} else { // Fall back on using base
				if v = atomic.LoadInt64(&s.base); fn == nil {
					newV = v + x
				} else {
					newV = fn.Apply(v, x)
				}
				if s.casBase(v, newV) {
					break
				}
			}
		}
	}
}
