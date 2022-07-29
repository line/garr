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

func TestRoundRobin(t *testing.T) {
	builder := defaultLBBuilder()

	lb := builder.Build()
	lb.Initialize(validEndpoints())

	endpoints := lb.Endpoints()
	valid := validEndpoints()
	if len(endpoints) != len(valid) {
		t.FailNow()
	}
	for i := range endpoints {
		if !endpoints[i].Equal(valid[i]) {
			t.FailNow()
		}
	}
	if lb.Pick()+lb.Pick()+lb.Pick() != 3 {
		t.FailNow()
	}
}

func TestRoundRobinRace(t *testing.T) {
	builder := defaultLBBuilder()

	lb := builder.Build()
	lb.Initialize(validEndpoints())

	type counter struct {
		v [3]int
	}
	ch := make(chan counter, 5)
	for i := 0; i < 7; i++ {
		go func() {
			time.Sleep(200 * time.Millisecond)

			var counting counter
			for j := 0; j < 18000; j++ {
				counting.v[lb.Pick()]++
			}

			ch <- counting
		}()
	}

	var sum counter
	for i := 0; i < 7; i++ {
		v := <-ch
		sum.v[0] += v.v[0]
		sum.v[1] += v.v[1]
		sum.v[2] += v.v[2]
	}

	// 42000 = 18000 / 3 * 7
	if sum.v[0] != 42000 || sum.v[1] != 42000 || sum.v[2] != 42000 {
		t.FailNow()
	}
}
