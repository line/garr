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

package httpclient

import (
	"context"
	"runtime"
	"time"

	workerpool "go.linecorp.com/garr/worker-pool"
)

var (
	numCPU = runtime.NumCPU()
)

type healthChecker struct {
	interval time.Duration
	timeout  time.Duration
	workers  *workerpool.Pool
}

func newHealthChecker(interval, timeout time.Duration) Resolver {
	if interval <= 0 {
		if interval = timeout << 1; interval <= 0 {
			interval = 200 * time.Millisecond
		}
	}

	h := &healthChecker{
		interval: interval,
		timeout:  timeout,
		workers: workerpool.NewPool(context.Background(), workerpool.Option{
			NumberWorker: numCPU,
		}),
	}
	h.workers.Start()

	return h
}

func (h *healthChecker) stopWorkers() {
	h.workers.Stop()
}

func doHealthCheck(index int, endpoint *Endpoint, timeout time.Duration, out chan int) *workerpool.Task {
	return workerpool.NewTask(context.Background(), func(context.Context) (interface{}, error) {
		if success := endpoint.Dial(timeout); success {
			out <- index
		} else {
			out <- -1
		}
		return nil, nil
	})
}

func (h *healthChecker) Resolve(in <-chan Endpoints, out chan<- Endpoints) {
	var (
		endpoints Endpoints
		ch        = make(chan int, 8)
	)

	do := func(endpoints Endpoints) {
		if n := len(endpoints); n > 0 {
			resolved := make([]*Endpoint, 0, n)

			for i := range endpoints {
				h.workers.Do(doHealthCheck(i, endpoints[i], h.timeout, ch))
			}

			for range endpoints {
				if ind := <-ch; ind >= 0 {
					resolved = append(resolved, endpoints[ind])
				}
			}

			// notify next resolver
			if len(resolved) > 0 {
				out <- resolved
			}
		}
	}

	// setup ticker
	ticker := time.NewTicker(h.interval)

	for {
		select {
		case eps, ok := <-in:
			if !ok {
				ticker.Stop()
				h.stopWorkers()
				return
			}

			if len(eps) > 0 && eps.normalize() == nil {
				endpoints = eps
				do(endpoints)
			}

		case <-ticker.C:
			do(endpoints)
		}
	}
}
