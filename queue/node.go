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
	"sync/atomic"
	"unsafe"
)

type linkedListNode struct {
	_v interface{}    // real value
	_i unsafe.Pointer // wrapper over value
	_n unsafe.Pointer // next
}

func newLinkedListNode(v interface{}) *linkedListNode {
	return &linkedListNode{
		_v: v,
		_i: unsafe.Pointer(&v),
		_n: nil,
	}
}

func (n *linkedListNode) value() interface{} {
	return n._v
}

func (n *linkedListNode) next() unsafe.Pointer {
	return atomic.LoadPointer(&n._n)
}

func (n *linkedListNode) item() unsafe.Pointer {
	return atomic.LoadPointer(&n._i)
}

func (n *linkedListNode) setItemNil() {
	atomic.StorePointer(&n._i, nil)
}

func (n *linkedListNode) casItemNil(old unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&n._i, old, nil)
}

func (n *linkedListNode) casNext(old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&n._n, old, new)
}

func (n *linkedListNode) setNext(new unsafe.Pointer) {
	atomic.StorePointer(&n._n, new)
}
