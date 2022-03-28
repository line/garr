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
	"context"
	"sync"
	"testing"
	"time"
)

func testProducer(t *testing.T, queue Queue) {
	// try offer nil
	queue.Offer(nil)

	var wg sync.WaitGroup
	for i := 0; i < maxNumberProducer; i++ {
		wg.Add(1)
		go func(producer int) {
			for j := 0; j < numberEle; j++ {
				queue.Offer(&ele{key: producer, value: j})
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	if queue.Peek() == nil || int(queue.Size()) != maxNumberProducer*numberEle || queue.IsEmpty() {
		t.Fatal()
	}

	m := make(map[int]map[int]struct{})
	for i := 0; i < maxNumberProducer*numberEle; i++ {
		if polled := queue.Poll(); polled == nil {
			t.Fatal()
		} else {
			e := polled.(*ele)
			if _, ok := m[e.key]; !ok {
				m[e.key] = make(map[int]struct{})
			}
			m[e.key][e.value] = struct{}{}
		}
	}

	for i := 0; i < maxNumberProducer; i++ {
		if len(m[i]) != numberEle {
			t.Fatal()
		}

		for k := range m[i] {
			if k < 0 || k >= numberEle {
				t.Fatal()
			}
		}
	}
}

func testMix(t *testing.T, q Queue, numberProducer, numberConsumer int) {
	for i := 0; i < numberProducer; i++ {
		go func(producer int) {
			time.Sleep(100 * time.Millisecond)
			for j := 0; j < numberEle; j++ {
				q.Offer(&ele{key: producer, value: j})
			}
		}(i)
	}

	ch := make(chan *ele, 1)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < numberConsumer; i++ {
		wg.Add(1)
		go testConsumer(ctx, &wg, q, ch)
	}

	m := make(map[int]map[int]struct{})
	for i := 0; i < numberProducer*numberEle; i++ {
		if polled := <-ch; polled == nil {
			t.Fatal()
		} else {
			if _, ok := m[polled.key]; !ok {
				m[polled.key] = make(map[int]struct{})
			}
			m[polled.key][polled.value] = struct{}{}
		}
	}

	cancel()
	wg.Wait()

	for i := 0; i < numberProducer; i++ {
		if len(m[i]) != numberEle {
			t.Fatal(len(m[i]), numberEle)
		}

		for k := range m[i] {
			if k < 0 || k >= numberEle {
				t.Fatal()
			}
		}
	}
}

func testConsumer(ctx context.Context, wg *sync.WaitGroup, q Queue, ch chan *ele) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		default:
			if item := q.Peek(); item != nil {
				if q.IsEmpty() {
					continue
				}

				if item = q.Poll(); item != nil {
					ch <- item.(*ele)
				}
			}
		}
	}
}
