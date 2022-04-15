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
	"testing"
	"time"
)

type noopResolver struct{}

func (l *noopResolver) Resolve(in <-chan Endpoints, out chan<- Endpoints) {
	for v := range in {
		time.Sleep(300 * time.Millisecond)
		out <- v
	}
}

func TestResolverChainWait(t *testing.T) {
	chain := newResolverChain([]Resolver{&noopResolver{}})
	chain.push(nil)

	readies := make(chan bool)
	go chainWait(chain, 50*time.Millisecond, readies)
	go chainWait(chain, 50*time.Millisecond, readies)
	go chainWait(chain, 50*time.Millisecond, readies)

	time.Sleep(300 * time.Millisecond)
	if !chain.wait(50 * time.Millisecond) {
		t.FailNow()
	}
	for i := 0; i < 3; i++ {
		if <-readies {
			t.FailNow()
		}
	}
	close(readies)
	chain.close()
}

func chainWait(chain resolverChain, timeout time.Duration, results chan bool) {
	results <- chain.wait(timeout)
}
