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

package queue

import (
	"math"
	"sync/atomic"
	"unsafe"

	"go.linecorp.com/garr/internal"
)

// JDKLinkedQueue represents jdk-based concurrent non blocking linked-list queue.
type JDKLinkedQueue struct {
	// The padding members 1 to 3 below are here to ensure each item is on a separate cache line.
	// This prevents false sharing and hence improves performance.
	_ internal.CacheLinePad
	h unsafe.Pointer // head
	_ internal.CacheLinePad
	t unsafe.Pointer // tail
	_ internal.CacheLinePad
}

// NewJDKLinkedQueue creates new JDKLinkedQueue.
func NewJDKLinkedQueue() *JDKLinkedQueue {
	q := &JDKLinkedQueue{
		t: unsafe.Pointer(&linkedListNode{}),
	}
	q.h = q.t
	return q
}

func (queue *JDKLinkedQueue) head() unsafe.Pointer {
	return atomic.LoadPointer(&queue.h)
}

func (queue *JDKLinkedQueue) tail() unsafe.Pointer {
	return atomic.LoadPointer(&queue.t)
}

// Offer inserts the specified element at the tail of this queue.
func (queue *JDKLinkedQueue) Offer(v interface{}) {
	if v != nil {
		newNode := unsafe.Pointer(newLinkedListNode(v))

		var oldT unsafe.Pointer

		t := queue.tail()
		p := t
		for {
			_p := (*linkedListNode)(p)
			if q := _p.next(); q == nil {
				// p is last node
				if _p.casNext(nil, newNode) {
					// Successful CAS is the linearization point
					// for e to become an element of this queue,
					// and for newNode to become "live".
					if p != t { // hop two nodes at a time
						queue.casTail(t, newNode) // Failure is OK.
					}
					return
				}
				// Lost CAS race to another thread; re-read next
			} else if p == q {
				// We have fallen off list.  If tail is unchanged, it
				// will also be off-list, in which case we need to
				// jump to head, from which all live nodes are always
				// reachable.  Else the new tail is a better bet.
				if oldT, t = t, queue.tail(); oldT != t { // t != (t = tail)?
					p = t
				} else {
					p = queue.head()
				}
			} else if p != t { // Check for tail updates after two hops.
				if oldT, t = t, queue.tail(); oldT != t {
					p = t
				} else {
					p = q
				}
			} else {
				p = q
			}
		}
	}
}

// Poll head element.
func (queue *JDKLinkedQueue) Poll() interface{} {
loop:
	for {
		var q unsafe.Pointer

		h := queue.head()
		p := h
		for ; ; p = q {
			_p := (*linkedListNode)(p)
			if item := _p.item(); item != nil && _p.casItemNil(item) {
				v := _p.value()
				// Successful CAS is the linearization point
				// for item to be removed from this queue.
				if p != h { // hop two nodes at a time
					q = _p.next()
					if q != nil {
						queue.updateHead(h, q)
					} else {
						queue.updateHead(h, p)
					}
				}
				return v
			}

			if q = _p.next(); q == nil {
				queue.updateHead(h, p)
				return nil
			}

			if p == q {
				goto loop
			}
		}
	}
}

// Peek return head element
func (queue *JDKLinkedQueue) Peek() interface{} {
loop:
	for {
		var q unsafe.Pointer

		h := queue.head()
		p := h
		for ; ; p = q {
			_p := (*linkedListNode)(p)
			if item := _p.item(); item != nil {
				v := _p.value()
				queue.updateHead(h, p)
				return v
			}

			if q = _p.next(); q == nil {
				queue.updateHead(h, p)
				return nil
			}

			if p == q {
				goto loop
			}
		}
	}
}

// Returns the first live (non-deleted) node on list, or null if none.
// This is yet another variant of poll/peek; here returning the
// first node, not element.  We could make peek() a wrapper around
// first(), but that would cost an extra volatile read of item,
// and the need to add a retry loop to deal with the possibility
// of losing a race to a concurrent poll().
func (queue *JDKLinkedQueue) first() unsafe.Pointer {
loop:
	for {
		var q unsafe.Pointer

		h := queue.head()
		p := h
		for ; ; p = q {
			_p := (*linkedListNode)(p)

			hasItem := _p.item() != nil
			if hasItem {
				queue.updateHead(h, p)
				return p
			}

			if q = _p.next(); q == nil {
				queue.updateHead(h, p)
				return nil
			}

			if p == q {
				goto loop
			}
		}
	}
}

// IsEmpty returns if this queue contains no elements.
func (queue *JDKLinkedQueue) IsEmpty() bool {
	return queue.first() == nil
}

// Size returns the number of elements in this queue. If this queue
// contains more than math.MaxInt32 elements, returns math.MaxInt32.
// Beware that, unlike in most collections, this method is
// NOT a constant-time operation. Because of the
// asynchronous nature of these queues, determining the current
// number of elements requires an O(n) traversal.
// Additionally, if elements are added or removed during execution
// of this method, the returned result may be inaccurate. Thus,
// this method is typically not very useful in concurrent applications.
func (queue *JDKLinkedQueue) Size() int32 {
loop:
	for {
		var count int32
		for p := queue.first(); p != nil; {
			_p := (*linkedListNode)(p)
			if _p.item() != nil {
				count++
				if count == math.MaxInt32 {
					return count
				}
			}

			oldP := p
			p = _p.next()
			if oldP == p {
				goto loop
			}
		}
		return count
	}
}

// Iterator returns iterator of underlying elements.
func (queue *JDKLinkedQueue) Iterator() Iterator {
	return newJdkLinkedQueueIter(queue)
}

func (queue *JDKLinkedQueue) casTail(old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&queue.t, old, new)
}

func (queue *JDKLinkedQueue) casHead(old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&queue.h, old, new)
}

func (queue *JDKLinkedQueue) updateHead(h, p unsafe.Pointer) {
	if h != p && queue.casHead(h, p) {
		(*linkedListNode)(h).setNext(h)
	}
}

func (queue *JDKLinkedQueue) succ(node unsafe.Pointer) unsafe.Pointer {
	old := node
	if node = (*linkedListNode)(node).next(); old == node {
		node = queue.head()
	}
	return node
}

type jdkLinkedQueueIter struct {
	q        *JDKLinkedQueue
	nextNode unsafe.Pointer
	nextItem unsafe.Pointer
	nextVal  interface{}
	lastRet  unsafe.Pointer
}

func newJdkLinkedQueueIter(queue *JDKLinkedQueue) (iter *jdkLinkedQueueIter) {
	iter = &jdkLinkedQueueIter{
		q: queue,
	}

loop:
	for {
		var q unsafe.Pointer

		p := iter.q.head()
		h := p
		for ; ; p = q {
			_p := (*linkedListNode)(p)

			item := _p.item()
			if item != nil {
				iter.nextNode = p
				iter.nextItem = item
				iter.nextVal = _p.value()
				break
			}

			q = _p.next()
			if q == nil {
				break
			}

			if p == q {
				goto loop
			}
		}

		iter.q.updateHead(h, p)
		return
	}
}

// HasNext returns true if has next.
func (i *jdkLinkedQueueIter) HasNext() bool {
	return i.nextItem != nil
}

// Next return next elements. There is no guarantee that hasNext and next are atomically due to data racy.
func (i *jdkLinkedQueueIter) Next() interface{} {
	pred := i.nextNode
	if pred == nil {
		return nil
	}

	i.lastRet = pred

	var q, item unsafe.Pointer
	var val interface{}
	for p := i.q.succ(pred); ; p = q {
		if p == nil {
			i.nextNode = p
			x := i.nextVal

			i.nextItem = item
			i.nextVal = val

			return x
		}

		_p := (*linkedListNode)(p)
		item, val = _p.item(), _p.value()
		if item != nil {
			i.nextNode = p
			x := i.nextVal

			i.nextItem = item
			i.nextVal = val

			return x
		}

		// unlink deleted nodes
		q = i.q.succ(p)
		if q != nil {
			(*linkedListNode)(pred).casNext(p, q)
		}
	}
}

// Remove from the underlying collection the last element returned
// by this iterator.
func (i *jdkLinkedQueueIter) Remove() {
	l := i.lastRet
	if l == nil {
		return
	}

	// rely on a future traversal to relink.
	_l := (*linkedListNode)(l)
	_l.setItemNil()
	i.lastRet = nil
}
