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
	"net"
	"net/url"
	"strconv"
	"time"

	cb "go.linecorp.com/garr/circuit-breaker"
)

// EndpointMetadata represents endpoint's metadata.
type EndpointMetadata struct {
	Weight uint `json:"weight" yaml:"weight"`
}

// Equal performs deep equal checking.
func (e *EndpointMetadata) Equal(other *EndpointMetadata) (eq bool) {
	eq = other != nil &&
		e.Weight == other.Weight
	return
}

// RawEndpoint is composite of raw url alongs and its metadata.
type RawEndpoint struct {
	RawURL   string           `json:"url" yaml:"url"`
	Metadata EndpointMetadata `json:"meta" yaml:"meta"`
}

// ToEndpoint converts to Endpoint.
func (r *RawEndpoint) ToEndpoint() (endp *Endpoint, err error) {
	u, err := url.Parse(r.RawURL)
	if err == nil {
		endp = &Endpoint{
			URL:      *u,
			Metadata: r.Metadata,
		}
		err = endp.normalize()
	}
	return
}

// Endpoint is composite of url.URL and its metadata.
type Endpoint struct {
	url.URL
	network  string
	breaker  cb.CircuitBreaker
	Metadata EndpointMetadata
}

// Equal performs deep equal checking.
func (e *Endpoint) Equal(other *Endpoint) (eq bool) {
	eq = other != nil &&
		e.URL.Host == other.URL.Host &&
		e.URL.Scheme == other.URL.Scheme &&
		e.URL.Opaque == other.URL.Opaque &&
		e.equalUser(other) &&
		e.Metadata.Equal(&other.Metadata)
	return
}

func (e *Endpoint) equalUser(other *Endpoint) (eq bool) {
	if e.URL.User == nil {
		eq = other.URL.User == nil
	} else if eq = other.URL.User != nil && e.URL.User.Username() == other.URL.User.Username(); eq {
		p1, s1 := e.URL.User.Password()
		p2, s2 := other.URL.User.Password()
		eq = s1 == s2 && p1 == p2
	}
	return
}

func (e *Endpoint) canRequest() bool {
	return e.breaker.CanRequest()
}

func (e *Endpoint) onConnectFailure() {
	e.breaker.OnFailure()
}

func (e *Endpoint) onConnectSuccess() {
	e.breaker.OnSuccess()
}

// setup circuit breaker
func (e *Endpoint) setupCB(builder *cb.CircuitBreakerBuilder) {
	e.breaker, _ = builder.Build()
}

func (e *Endpoint) normalize() (err error) {
	network, p, err := lookupPortByScheme(e.Scheme)
	if err == nil {
		if e.Port() == "" {
			e.Host = net.JoinHostPort(e.Host, strconv.Itoa(p))
		}
		e.network = network
	}
	return
}

// Dial endpoint.
func (e *Endpoint) Dial(timeout time.Duration) (success bool) {
	var (
		conn net.Conn
		err  error
	)

	if timeout > 0 {
		conn, err = net.DialTimeout(e.network, e.Host, timeout)
	} else {
		conn, err = net.Dial(e.network, e.Host)
	}

	if err == nil {
		if conn != nil {
			_ = conn.Close()
		}
		success = true
	}

	return
}

// Endpoints represents a collection of endpoint(s).
type Endpoints []*Endpoint

// Equal performs deep equal checking.
func (e Endpoints) Equal(other Endpoints) (eq bool) {
	if len(e) == len(other) {
		for i := range e {
			if !e[i].Equal(other[i]) {
				return
			}
		}
		eq = true
	}
	return
}

// Clone endpoints.
func (e Endpoints) Clone() Endpoints {
	eps := make([]*Endpoint, len(e))
	copy(eps, e)
	return eps
}

func (e Endpoints) normalize() (err error) {
	for i := range e {
		if err = e[i].normalize(); err != nil {
			return
		}
	}
	return
}

// Parse endpoint(s) from raw(s).
func Parse(raws []RawEndpoint) (eps Endpoints, err error) {
	eps = make([]*Endpoint, len(raws))

	for i := range raws {
		if eps[i], err = raws[i].ToEndpoint(); err != nil {
			return
		}
	}

	return
}

// ParseFromURLs endpoint(s) from url(s).
func ParseFromURLs(urls []string) (eps Endpoints, err error) {
	raws := make([]RawEndpoint, len(urls))
	for i := range urls {
		raws[i].RawURL = urls[i]
	}
	return Parse(raws)
}
