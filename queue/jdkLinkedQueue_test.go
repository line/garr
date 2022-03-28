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
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestJDKLinkedQueue_Producer(t *testing.T) {
	testProducer(t, NewQueue(JDKLinkedQueueType))
}

func TestJDKLinkedQueue_Mix(t *testing.T) {
	testMix(t, DefaultQueue(), 10, 10)
}

type anonymousStruct struct {
	v int
}

func TestJDKLinkedQueue_Iterator(t *testing.T) {
	q := NewJDKLinkedQueue()

	as := make([]*anonymousStruct, 100)
	for i := range as {
		as[i] = &anonymousStruct{v: i}
		q.Offer(as[i])
	}

	count := 0

	iter := q.Iterator()
	for iter.HasNext() {
		count++
		v := iter.Next().(*anonymousStruct)
		if v.v < 50 {
			iter.Remove()
		}
	}

	if count == 0 {
		t.Fatal(count)
	}

	if q.Size() != 50 {
		t.Fatal(q.Size())
	}
}

func TestJDKLinkedQueueRace(t *testing.T) {
	q := NewJDKLinkedQueue()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(200 * time.Millisecond)
			for j := 0; j < 100000; j++ {
				q.Offer(j)
			}
		}()
	}

	var counter uint32
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(400 * time.Millisecond)
			for j := 0; j < 30000; j++ {
				if polled := q.Poll(); polled != nil {
					v := polled.(int)
					if v >= 20000 {
						atomic.AddUint32(&counter, 1)
					}
				}
			}
		}()
	}

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(200 * time.Millisecond)
			for iter := q.Iterator(); iter.HasNext(); {
				next := iter.Next()
				if next != nil {
					v := next.(int)
					if v < 20000 {
						iter.Remove()
					}
				}
			}
		}()
	}
	wg.Wait()

	t.Log(counter)
}
