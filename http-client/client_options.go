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
	"time"

	cb "go.linecorp.com/garr/circuit-breaker"
	"go.linecorp.com/garr/retry"
)

// ClientOption is setter for client's option(s).
type ClientOption func(c *Client)

// CheckRedirect specifies the policy for handling redirects.
// If CheckRedirect is not nil, the client calls it before
// following an HTTP redirect. The arguments req and via are
// the upcoming request and the requests made already, oldest
// first. If CheckRedirect returns an error, the Client's Get
// method returns both the previous Response (with its Body
// closed) and CheckRedirect's error (wrapped in a url.Error)
// instead of issuing the Request req.
// As a special case, if CheckRedirect returns ErrUseLastResponse,
// then the most recent response is returned with its body
// unclosed, along with a nil error.
//
// If CheckRedirect is nil, the Client uses its default policy,
// which is to stop after 10 consecutive requests.
// 	CheckRedirect func(req *Request, via []*Request) error

// WithTransport attaches transport with client.
//
// Transport specifies the mechanism by which individual
// HTTP requests are made.
// If nil, DefaultTransport is used.
func WithTransport(t http.RoundTripper) ClientOption {
	return func(c *Client) {
		c.c.Transport = t
	}
}

// WithTimeout specifies a time limit for requests made by this
// Client. The timeout includes connection time, any
// redirects, and reading the response body. The timer remains
// running after Get, Head, Post, or Do return and will
// interrupt reading of the response body.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.c.Timeout = timeout
	}
}

// WithCookieJar specifies the cookie jar.
//
// The Jar is used to insert relevant cookies into every
// outbound Request and is updated with the cookie values
// of every inbound Response. The Jar is consulted for every
// redirect that the Client follows.
//
// If Jar is nil, cookies are only sent if they are explicitly
// set on the Request.
func WithCookieJar(jar http.CookieJar) ClientOption {
	return func(c *Client) {
		c.c.Jar = jar
	}
}

// WithCheckRedirect specifies the policy for handling redirects.
//
// If checkRedirect is not nil, the client calls it before
// following an HTTP redirect. The arguments req and via are
// the upcoming request and the requests made already, oldest
// first. If checkRedirect returns an error, the Client's Get
// method returns both the previous Response (with its Body
// closed) and checkRedirect's error (wrapped in a url.Error)
// instead of issuing the Request req.
//
// As a special case, if checkRedirect returns http.ErrUseLastResponse,
// then the most recent response is returned with its body
// unclosed, along with a nil error.
//
// If checkRedirect is nil, the Client uses its default policy,
// which is to stop after 10 consecutive requests.
func WithCheckRedirect(checkRedirect func(req *http.Request, via []*http.Request) error) ClientOption {
	return func(c *Client) {
		c.c.CheckRedirect = checkRedirect
	}
}

// WithResolver appends endpoints resolver to chain.
//
// Chain of resolvers are used to resolve endpoints' url.
// It's disabled when endpoint factory presents.
//
// Chain respects the order of ClientOption arguments.
// Thus, if client is built with:
//
//  client := NewClient(WithResolver(r3), WithResolver(r1), WithResolver(r2))
//  // order of resolving will be
//  r3 -> r1 -> r2
func WithResolver(r Resolver) ClientOption {
	return func(c *Client) {
		c.resolvers = append(c.resolvers, r)
	}
}

// WithResolvers appends endpoints resolvers to chain.
func WithResolvers(r []Resolver) ClientOption {
	return func(c *Client) {
		c.resolvers = append(c.resolvers, r...)
	}
}

// WithHealthChecker attachs health check resolver. Interval indicates
// the duration between each check and timeout indicates for tcp dial timeout.
//
// Default:
// - interval: 500 milllis
// - timeout: 100 millis
func WithHealthChecker(interval, timeout time.Duration) ClientOption {
	return func(c *Client) {
		if c.healthChecker != nil {
			c.healthChecker.(*healthChecker).stopWorkers()
		}
		c.healthChecker = newHealthChecker(interval, timeout)
	}
}

// WithLoadBalanceBuilder specifies load balance builder which is
// used to generate LB on demand.
//
// Default: PickfirstLB
func WithLoadBalanceBuilder(builder LBBuilder) ClientOption {
	return func(c *Client) {
		c.lbBuilder = builder
	}
}

// WithCircuitBreakerBuilder specifies circuit breaker builder which
// is used to generate CB on deman.
//
// Default using circuit breaker with settings:
//   FailureRateThreshold    = 0.8
// 	 MinimumRequestThreshold = 10
// 	 TrialRequestInterval    = time.Duration(3 * time.Second)
// 	 CircuitOpenWindow       = time.Duration(10 * time.Second)
// 	 CounterSlidingWindow    = time.Duration(20 * time.Second)
// 	 CounterUpdateInterval   = time.Duration(1 * time.Second)
func WithCircuitBreakerBuilder(builder *cb.CircuitBreakerBuilder) ClientOption {
	return func(c *Client) {
		c.cbBuilder = builder
	}
}

// WithBackoff specifies retry-backoff.
//
// Default:
// exponential:
// - initialDelayMillis: 50
// - maxDelayMillis: 5000 (5 seconds)
// - multipiler: 1.15
// jitter: 0.1
// limit: 3 (try 3 times)
func WithBackoff(b retry.Backoff) ClientOption {
	return func(c *Client) {
		c.backoff = b
	}
}
