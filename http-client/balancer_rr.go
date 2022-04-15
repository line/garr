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
	"sync/atomic"

	"github.com/valyala/fastrand"
)

// RoundRobinLB is round robin strategy.
type RoundRobinLB struct {
	endpoints Endpoints
	index     uint32
	n         uint32
}

// Initialize endpoints for RoundRobinLB.
func (p *RoundRobinLB) Initialize(endpoints Endpoints) {
	p.endpoints = endpoints
	p.n = uint32(len(endpoints))
}

// Endpoints returned saved endpoints inside RoundRobinLB
func (p *RoundRobinLB) Endpoints() Endpoints {
	return p.endpoints
}

// Pick returns index of picked endpoint.
func (p *RoundRobinLB) Pick() (chosen int) {
	if p.n > 0 {
		chosen = int(atomic.AddUint32(&p.index, 1) % p.n)
	}
	return
}

// RoundRobinLBBuilder is builder for RoundRobin load-balancer
type RoundRobinLBBuilder struct{}

// Build round robin balancer.
func (p *RoundRobinLBBuilder) Build() LB {
	return &RoundRobinLB{index: fastrand.Uint32()}
}
