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

package cbreaker

import (
	"fmt"
	"sync/atomic"
	"time"
	"unsafe"

	ga "github.com/line/garr/adder"
	queue "github.com/line/garr/queue"
)

// bucket hold the count of events within {@code updateInterval}.
type bucket struct {
	timestamp int64
	s         ga.LongAdder // number of success
	f         ga.LongAdder // number of failure
}

// create new bucket
func newBucket(timestamp int64) *bucket {
	return &bucket{
		timestamp: timestamp,
		s:         ga.NewLongAdder(ga.JDKAdderType),
		f:         ga.NewLongAdder(ga.JDKAdderType),
	}
}

func (b *bucket) add(succ bool) {
	if succ {
		b.s.Add(1)
	} else {
		b.f.Add(1)
	}
}

// returns number of success operation
func (b *bucket) success() int64 {
	return b.s.Sum()
}

// returns number of failure operation
func (b *bucket) failure() int64 {
	return b.f.Sum()
}

// SlidingWindowCounter accumulates the count of events within a time window.
type SlidingWindowCounter struct {
	ticker              Ticker
	cur                 *bucket
	slidingWindowNanos  int64
	updateIntervalNanos int64
	snapshot            atomic.Value
	reservoir           queue.Queue
}

// NewSlidingWindowCounter creates new SlidingWindowCounter.
func NewSlidingWindowCounter(ticker Ticker, slidingWindowNanos, updateIntervalNanos time.Duration) (s *SlidingWindowCounter, e error) {
	if ticker == nil {
		e = fmt.Errorf("Ticker is required")
		return
	}

	s = &SlidingWindowCounter{
		ticker:              ticker,
		slidingWindowNanos:  int64(slidingWindowNanos),
		updateIntervalNanos: int64(updateIntervalNanos),
		cur:                 newBucket(ticker.Tick()),
		reservoir:           queue.DefaultQueue(),
	}
	s.snapshot.Store(EventCountZero)
	return
}

func (s *SlidingWindowCounter) current() *bucket {
	return (*bucket)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.cur))))
}

func (s *SlidingWindowCounter) casCurrent(old, new *bucket) bool {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&s.cur)),
		unsafe.Pointer(old),
		unsafe.Pointer(new),
	)
}

// Count returns the current EventCount.
func (s *SlidingWindowCounter) Count() *EventCount {
	return s.snapshot.Load().(*EventCount)
}

// OnSuccess adds success event to the current bucket.
// Return success/failure count of current sliding window on every update interval.
// Return nil if events are still being accumulated in current bucket.
func (s *SlidingWindowCounter) OnSuccess() *EventCount {
	return s.onEvent(true)
}

// OnFailure adds failure event to the current bucket.
// Return success/failure count of current sliding window on every update interval.
// Return nil if events are still being accumulated in current bucket
func (s *SlidingWindowCounter) OnFailure() *EventCount {
	return s.onEvent(false)
}

func (s *SlidingWindowCounter) onEvent(succ bool) (e *EventCount) {
	tickerNanos, currentBucket := s.ticker.Tick(), s.current()

	if tickerNanos < currentBucket.timestamp {
		// if current timestamp is older than bucket's timestamp (maybe race or GC pause?),
		// then creates an instant bucket and puts it to the reservoir not to lose event.
		bucket := newBucket(tickerNanos)
		bucket.add(succ)
		s.reservoir.Offer(bucket)
		return
	}

	if tickerNanos < currentBucket.timestamp+s.updateIntervalNanos {
		// Events are still being accumulated in current bucket.
		currentBucket.add(succ)
		return
	}

	nextBucket := newBucket(tickerNanos)
	nextBucket.add(succ)

	// replaces the bucket
	if s.casCurrent(currentBucket, nextBucket) {
		// puts old one to the reservoir
		s.reservoir.Offer(currentBucket)

		// and then updates count
		e = s.trimAndSum(tickerNanos)
		s.snapshot.Store(e)
	} else {
		// the bucket has been replaced already
		// puts new one as an instant bucket to the reservoir not to lose event
		s.reservoir.Offer(nextBucket)
	}
	return
}

func (s *SlidingWindowCounter) trimAndSum(t int64) *EventCount {
	oldLimit, iterator := t-s.slidingWindowNanos, s.reservoir.Iterator()

	var nxt interface{}
	var bck *bucket
	var success, failure int64

	for iterator.HasNext() {
		if nxt = iterator.Next(); nxt != nil {
			if bck = nxt.(*bucket); bck.timestamp < oldLimit {
				// removes old bucket
				iterator.Remove()
			} else {
				success += bck.success()
				failure += bck.failure()
			}
		}
	}

	return NewEventCount(success, failure)
}
