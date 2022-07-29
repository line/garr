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
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"

	cb "go.linecorp.com/garr/circuit-breaker"
	"go.linecorp.com/garr/retry"
)

var (
	knownNetworks = []string{"tcp", "udp", ""}
)

func lookupPortByScheme(scheme string) (network string, port int, err error) {
	for i := range knownNetworks {
		if port, err = net.LookupPort(knownNetworks[i], scheme); err == nil {
			network = knownNetworks[i]
			return
		}
	}
	return
}

func injectTarget(endpoint *Endpoint, req *Request) (originalURL *url.URL) {
	u := req.r.URL

	originalURL = &url.URL{
		Scheme: u.Scheme,
		Opaque: u.Opaque,
		User:   u.User,
		Host:   u.Host,
		Path:   u.Path,
	}

	u.Scheme = endpoint.Scheme
	u.Opaque = endpoint.Opaque
	u.User = endpoint.User
	u.Host = endpoint.Host
	u.Path = path.Join(endpoint.Path, u.Path)

	return
}

func revert(req *Request, originalURL *url.URL) {
	u := req.r.URL

	u.Scheme = originalURL.Scheme
	u.Opaque = originalURL.Opaque
	u.User = originalURL.User
	u.Host = originalURL.Host
	u.Path = originalURL.Path
}

func decode(req *Request, resp *http.Response) (err error) {
	if req.expect != nil {
		if w, ok := req.expect.(io.Writer); ok {
			if _, err = io.Copy(w, resp.Body); err == nil {
				if closer, ok := req.expect.(io.Closer); ok {
					err = closer.Close()
				}
			}
		} else if req.decoder != nil {
			err = req.decoder(resp.Body, req.expect)
		}
	}
	return
}

func drainAndClose(r io.ReadCloser) {
	_, _ = io.Copy(ioutil.Discard, r)
	_ = r.Close()
}

// default circuit breaker settings:
//   defaultFailureRateThreshold    = 0.8
// 	 defaultMinimumRequestThreshold = 10
// 	 defaultTrialRequestInterval    = time.Duration(3 * time.Second)
// 	 defaultCircuitOpenWindow       = time.Duration(10 * time.Second)
// 	 defaultCounterSlidingWindow    = time.Duration(20 * time.Second)
// 	 defaultCounterUpdateInterval   = time.Duration(1 * time.Second)
//
// See also: https://line.github.io/armeria/client-circuit-breaker.html
func defaultCircuitBreakerBuilder() *cb.CircuitBreakerBuilder {
	return cb.NewCircuitBreakerBuilder()
}

// default backoff, using:
//
// exponential:
// - initialDelayMillis: 50
// - maxDelayMillis: 5000 (5 seconds)
// - multipiler: 1.15
// jitter: 0.1
// limit: 3 (try 3 times)
//
// See also: https://line.github.io/armeria/client-retry.html#backoff
func defaultBackoff() (b retry.Backoff) {
	baseBackoff, _ := retry.NewExponentialBackoff(50, 5000, 1.15)
	b, _ = retry.NewBackoffBuilder().
		BaseBackoff(baseBackoff).
		WithJitter(0.1).
		WithLimit(3).
		Build()
	return
}

// default health check resolver:
// - interval: 500 Millis
// - timeout: 100 Millis
func defaultHealthChecker() Resolver {
	return newHealthChecker(500*time.Millisecond, 100*time.Millisecond)
}
