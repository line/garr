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
	"context"
	"sync"
	"testing"
	"time"
)

func Test_NewSlidingWindow(t *testing.T) {
	if _, err := NewSlidingWindowCounter(nil, 456, 123); err == nil {
		t.Errorf("Fail to create new sliding window counter")
	}

	if s, err := NewSlidingWindowCounter(SystemTicker, 456, 123); err != nil || s.slidingWindowNanos != 456 || s.updateIntervalNanos != 123 {
		t.Errorf("Fail to create new sliding window counter")
	}

	// test casCurrent
	counter, _ := NewSlidingWindowCounter(
		SystemTicker,
		50*time.Nanosecond,
		5*time.Second,
	)

	nextBucket := newBucket(SystemTicker.Tick())
	current := counter.current()
	if !counter.casCurrent(current, nextBucket) || counter.current() != nextBucket {
		t.Fatal()
	}
}

func Test_StressSlidingWindow(t *testing.T) {
	counter, _ := NewSlidingWindowCounter(
		SystemTicker,
		50*time.Millisecond,
		6*time.Second,
	)

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					counter.OnSuccess()
				}
			}
		}()
	}

	time.Sleep(20 * time.Second)
	cancel()
	wg.Wait()
}
