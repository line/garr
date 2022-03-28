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

// Type is the type of queue.
type Type byte

const (
	// JDKLinkedQueueType indicates JDKLinkedQueue.
	JDKLinkedQueueType Type = iota
	// MutexLinkedQueueType indicates MutexLinkedQueue.
	MutexLinkedQueueType
)

// Queue interface.
type Queue interface {
	// Offer inserts the specified element into this queue if it is possible to do so immediately
	// without violating capacity restrictions. Return nil if success.
	Offer(v interface{})
	// Poll retrieve and removes the head of this queue, or returns nil if this queue is empty.
	Poll() interface{}
	// Peek retrieve, but does not remove, the head of this queue, or returns nil if this queue is empty.
	Peek() interface{}
	// Size returns the size of current queue.
	Size() int32
	// IsEmpty returns if this queue contains no elements.
	IsEmpty() bool
	// Iterator returns an iterator over the elements in this collection.
	Iterator() Iterator
}

// Iterator interface.
type Iterator interface {
	// HasNext returns true if the iteration has more elements.
	HasNext() bool
	// Next returns the next element in the iteration.
	Next() interface{}
	// Remove from the underlying collection the last element returned
	// by this iterator.
	Remove()
}

// NewQueue create new queue based on type
func NewQueue(t Type) Queue {
	switch t {
	case MutexLinkedQueueType:
		return NewMutexLinkedQueue()
	default:
		return NewJDKLinkedQueue()
	}
}

// DefaultQueue returns jdk concurrent, non blocking queue.
func DefaultQueue() Queue {
	return NewJDKLinkedQueue()
}
