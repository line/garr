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
	"container/list"
	"sync"
)

// MutexLinkedQueue mutex-based concurrent linked list queue.
type MutexLinkedQueue struct {
	l     *list.List
	mutex sync.RWMutex
}

// NewMutexLinkedQueue creates new MutexLinkedQueue.
func NewMutexLinkedQueue() *MutexLinkedQueue {
	return &MutexLinkedQueue{
		l: list.New(),
	}
}

// Offer inserts the specified element into this queue if it is possible to do so immediately
// without violating capacity restrictions. Return nil if success.
func (queue *MutexLinkedQueue) Offer(v interface{}) {
	if v != nil {
		queue.mutex.Lock()
		queue.l.PushBack(v)
		queue.mutex.Unlock()
	}
}

// Poll retrieve and removes the head of this queue, or returns nil if this queue is empty.
func (queue *MutexLinkedQueue) Poll() (v interface{}) {
	queue.mutex.Lock()
	if e := queue.l.Front(); e != nil {
		v = e.Value
		queue.l.Remove(e)
	}
	queue.mutex.Unlock()
	return
}

// Peek retrieve, but does not remove, the head of this queue, or returns nil if this queue is empty.
func (queue *MutexLinkedQueue) Peek() (v interface{}) {
	queue.mutex.RLock()
	if e := queue.l.Front(); e != nil {
		v = e.Value
	}
	queue.mutex.RUnlock()
	return
}

// Size returns the number of elements in this queue. If this queue
// contains more than math.MaxInt32 elements, returns math.MaxInt32.
func (queue *MutexLinkedQueue) Size() (size int32) {
	queue.mutex.RLock()
	size = int32(queue.l.Len())
	queue.mutex.RUnlock()
	return
}

// IsEmpty returns if this queue contains no elements
func (queue *MutexLinkedQueue) IsEmpty() (empt bool) {
	return queue.Size() == 0
}

// Iterator not supported. MutexLinkedQueue not support iterator.
func (queue *MutexLinkedQueue) Iterator() Iterator {
	return nil
}
