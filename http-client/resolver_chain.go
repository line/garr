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
	"sync"
	"sync/atomic"
	"time"
)

// runtime for chain of resolver(s)
type resolverChain struct {
	wg         *sync.WaitGroup
	nResolvers int
	resolvers  []Resolver
	pipes      []chan Endpoints
	ready      int32
}

func newResolverChain(resolvers []Resolver) (chain resolverChain) {
	// setup pipeline(s)
	n := len(resolvers)
	pipes := make([]chan Endpoints, n+1)
	for i := range pipes {
		pipes[i] = make(chan Endpoints, 1)
	}

	// setup chain
	chain = resolverChain{
		wg:         &sync.WaitGroup{},
		nResolvers: n,
		resolvers:  resolvers,
		pipes:      pipes,
	}

	// run resolver(s)
	chain.wg.Add(n)
	for i, r := range resolvers {
		go func(r Resolver, in, out chan Endpoints) {
			r.Resolve(in, out)
			close(out)
			chain.wg.Done()
		}(r, pipes[i], pipes[i+1])
	}

	return
}

func (r *resolverChain) close() {
	if r.nResolvers > 0 {
		close(r.pipes[0])
		r.wg.Wait()
	}
}

func (r *resolverChain) push(endpoints Endpoints) {
	r.pipes[0] <- endpoints
}

func (r *resolverChain) wait(timeout time.Duration) (ready bool) {
	if atomic.CompareAndSwapInt32(&r.ready, 0, 1) {
		tm := time.NewTimer(timeout)

		select {
		case <-r.pipes[r.nResolvers]:
			tm.Stop()
			atomic.StoreInt32(&r.ready, 2)
			ready = true

		case <-tm.C:
			atomic.StoreInt32(&r.ready, 0)
		}

	} else {
		ready = atomic.LoadInt32(&r.ready) == 2
	}
	return
}
