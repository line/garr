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
	"sync/atomic"
	"testing"
)

const (
	maxNumberProducer = 8
	numberEle         = 10000
)

type ele struct {
	key   int
	value int
}

var (
	preparedJDKQueue   Queue
	preparedMutexQueue Queue
)

func init() {
	preparedJDKQueue = NewQueue(JDKLinkedQueueType)
	prepare(preparedJDKQueue)

	preparedMutexQueue = NewQueue(MutexLinkedQueueType)
	prepare(preparedMutexQueue)
}

func prepare(q Queue) {
	for i := 0; i < numberEle; i++ {
		q.Offer(&ele{key: i, value: i})
	}
}

func Benchmark_MutexLinkedQueue_50P50C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(MutexLinkedQueueType), 50, 50)
	}
}

func Benchmark_JDKLinkedQueue_50P50C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(JDKLinkedQueueType), 50, 50)
	}
}

func Benchmark_MutexLinkedQueue_50P10C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(MutexLinkedQueueType), 50, 10)
	}
}

func Benchmark_JDKLinkedQueue_50P10C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(JDKLinkedQueueType), 50, 10)
	}
}

func Benchmark_MutexLinkedQueue_10P50C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(MutexLinkedQueueType), 10, 50)
	}
}

func Benchmark_JDKLinkedQueue_10P50C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(JDKLinkedQueueType), 10, 50)
	}
}

func benchQueueMix(q Queue, numberProducer, numberConsumer int) {
	var done int32
	for i := 0; i < numberProducer; i++ {
		go func(i int) {
			for j := 0; j < numberEle; j++ {
				q.Offer(&ele{key: i, value: -j})
			}
			atomic.AddInt32(&done, 1)
		}(i)
	}

	ch := make(chan *ele, 1)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < numberConsumer; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				default:
					if item := q.Poll(); item != nil {
						ch <- item.(*ele)
					}
				}
			}
		}(ctx)
	}

	for i := 0; i < numberProducer*numberEle; i++ {
		<-ch
	}

	cancel()
	wg.Wait()
}

func Benchmark_MutexLinkedQueue_100P(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchProducer(NewQueue(MutexLinkedQueueType), 100)
	}
}
func Benchmark_JDKLinkedQueue_100P(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchProducer(NewQueue(JDKLinkedQueueType), 100)
	}
}

func Benchmark_MutexLinkedQueue_100C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchProducer(preparedMutexQueue, 100)
	}
}
func Benchmark_JDKLinkedQueue_100C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchProducer(preparedJDKQueue, 100)
	}
}

func benchProducer(q Queue, numberProducer int) {
	var wg sync.WaitGroup
	for i := 0; i < numberProducer; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < numberEle; j++ {
				q.Offer(&ele{value: -j})
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
