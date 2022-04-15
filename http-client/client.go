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
	"net/http"
	"sync/atomic"
	"time"

	cb "go.linecorp.com/garr/circuit-breaker"
	"go.linecorp.com/garr/retry"
)

// Client instruments http.Client with feature rich.
type Client struct {
	// native http client
	c *http.Client

	// resolvers
	resolvers     []Resolver
	healthChecker Resolver
	chain         resolverChain
	chainReady    int32

	// load balancer
	lbBuilder LBBuilder
	lb        atomic.Value // LB

	// circuit breaker builder
	cbBuilder *cb.CircuitBreakerBuilder

	// retry-backoff strategy/algorithm
	backoff retry.Backoff
}

// NewClient returns new Client with user defined options.
//
// Initialization timeout specifies wait duration for client to resolve endpoint(s)
// through resolver(s).
func NewClient(initializationTimeout time.Duration, endpoints Endpoints, opts ...ClientOption) (c *Client, err error) {
	// create new client
	c = &Client{
		c: &http.Client{},
	}

	for i := range opts {
		opts[i](c)
	}

	// validate endpoint(s)
	if len(endpoints) > 0 {
		err = endpoints.normalize()
	} else if len(c.resolvers) == 0 {
		err = ErrNoEndpoints
	}

	if err == nil {
		// if not defined -> using default circuit breaker builder
		if c.cbBuilder == nil {
			c.cbBuilder = defaultCircuitBreakerBuilder()
		}

		// try to build one breaker to test setting(s)
		_, err = c.cbBuilder.Build()
	}

	if err == nil {
		// check balancer builder
		// if not define -> using RoundRobin
		if c.lbBuilder == nil {
			c.lbBuilder = defaultLBBuilder()
		}

		// if not defined -> using default backoff
		if c.backoff == nil {
			c.backoff = defaultBackoff()
		}

		// if not defined -> using default health checker
		if c.healthChecker == nil {
			c.healthChecker = defaultHealthChecker()
		}

		// initialize chain of resolver(s)
		c.resolvers = append(c.resolvers, c.healthChecker, c)
		c.chain = newResolverChain(c.resolvers)

		// warm up chain if need
		if len(endpoints) > 0 {
			c.chain.push(endpoints)
		}

		// get first resolved endpoint(s) through chain
		c.chain.wait(initializationTimeout)
	}

	return
}

// Close stops client and underlying daemons.
func (c *Client) Close() (err error) {
	c.chain.close()
	return
}

// Resolve endpoints and push to next resolver in chain.
func (c *Client) Resolve(in <-chan Endpoints, out chan<- Endpoints) {
	var currentEndpoints Endpoints
	for endpoints := range in {
		if len(endpoints) > 0 {
			// create new LB
			lb := c.lbBuilder.Build()
			lb.Initialize(endpoints)

			// compare with current endpoints/LB
			if endpoints = lb.Endpoints(); !endpoints.Equal(currentEndpoints) {
				// there are changes -> assign
				currentEndpoints = endpoints

				// inject new circuit breaker
				for i := range endpoints {
					endpoints[i].setupCB(c.cbBuilder)
				}

				// update LB
				c.lb.Store(lb)

				// notify chain ready
				if atomic.CompareAndSwapInt32(&c.chainReady, 0, 1) {
					out <- nil
				}
			}
		}
	}
}

func (c *Client) loadLB() LB {
	lb, _ := c.lb.Load().(LB)
	return lb
}
