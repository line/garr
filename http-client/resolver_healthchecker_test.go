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

func TestHealthChecker(t *testing.T) {
	h := newHealthChecker(200*time.Millisecond, 500*time.Millisecond)
	testHealthChecker(t, h, 3)
}

func TestHealthCheckerInvalidInput(t *testing.T) {
	h := newHealthChecker(-1, 100).(*healthChecker)
	if h.interval != 200 {
		t.FailNow()
	}
	h.stopWorkers()

	h = newHealthChecker(-1, -2).(*healthChecker)
	if h.interval != 200*time.Millisecond {
		t.FailNow()
	}
	h.stopWorkers()
}

func TestHealthCheckerInfiniteTimeout(t *testing.T) {
	h := newHealthChecker(100*time.Millisecond, 0)
	testHealthChecker(t, h, 3)
}

func testHealthChecker(t *testing.T, h Resolver, expectNumResolvedEndpoints int) {
	// initialize input and output chan
	in, out := make(chan Endpoints, 1), make(chan Endpoints, 1)

	// start resolving
	go h.Resolve(in, out)

	in <- validEndpoints()
	v := <-out
	if v == nil {
		t.FailNow()
	}
	if len(v) != expectNumResolvedEndpoints {
		t.FailNow()
	}

	close(in) // to notify health checker
}
