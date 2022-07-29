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
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
)

// Do sends a request and returns a response, following
// policy (such as redirects, cookies, auth) as configured on the
// client.
func (c *Client) Do(req *Request) (resp *Response) {
	resp = &Response{}

	lb := c.loadLB()
	if lb == nil {
		resp.err = ErrNoEndpoints
		return
	}

	endpoints := lb.Endpoints()
	if len(endpoints) == 0 {
		resp.err = ErrNoEndpoints
		return
	}

	// executing request
	picked := lb.Pick()
	need := c.instrument(endpoints[picked], req, resp)

	// we need some more acts, retry or try on next endpoint?
	if need != None && len(endpoints) > 1 {
		c.judge(need, endpoints, picked, req, resp)
	}

	return
}

func (c *Client) judge(need EndpointAction, endpoints Endpoints, pickedEndpoint int, req *Request, resp *Response) {
	var (
		n               = len(endpoints)
		ind             = pickedEndpoint
		lastErr         *multierror.Error
		err             error
		retryCount      int
		nextDelayMillis int64
	)

loop:
	switch need {
	case Retrying:
		// inc retry counter
		retryCount++

		if nextDelayMillis = c.backoff.NextDelayMillis(retryCount); nextDelayMillis >= 0 {
			// wait before retrying
			if nextDelayMillis > 0 {
				time.Sleep(time.Duration(nextDelayMillis) * time.Millisecond)
			}

			// reset response
			resp.reset()

			// retrying
			need = c.instrument(endpoints[ind], req, resp)

			goto loop
		} else if resp.err == nil {
			resp.err = fmt.Errorf("Retried host:[%v] url:[%v] but failed. Attempts so far: %d", endpoints[ind].URL.Host, req.r.URL, retryCount)
		}

	case NextEndpoint:
		// reset retry count
		retryCount = 0

		if ind++; ind == n {
			ind = 0
		}

		if ind != pickedEndpoint {
			// recording last error
			if resp.err != nil {
				lastErr = multierror.Append(lastErr, resp.err)
			}

			// reset response
			resp.reset()

			// try on next endpoint
			need = c.instrument(endpoints[ind], req, resp)

			goto loop
		} else { // loop all over but still failed

			// aggregating errors
			if lastErr != nil {
				err = lastErr.ErrorOrNil()
			}

			// build up final error
			if err != nil {
				resp.err = fmt.Errorf("Retrying request on all endpoints but failed. Last error: %v", err)
			} else {
				resp.err = ErrEndpointsUnavailable
			}
		}
	}
}

func (c *Client) instrument(endpoint *Endpoint, req *Request, resp *Response) (action EndpointAction) {
	// check if endpoint could make request
	if endpoint.canRequest() {
		originalURL := injectTarget(endpoint, req)
		action = c.exec(endpoint, req, resp)
		revert(req, originalURL)
	} else {
		// notify that we need to try on next endpoint
		action = NextEndpoint
	}
	return
}

func (c *Client) exec(endpoint *Endpoint, req *Request, resp *Response) (action EndpointAction) {
	// mark instrumented request that response belongs to
	resp.req = req

	// mark as noop (default)
	action = None

	// execute request
	if _resp, err := c.c.Do(req.r); err == nil {
		// report connect success to CB
		endpoint.onConnectSuccess()

		// mark original response for later draining
		originalResponseBody := _resp.Body

		// verify with header-judger
		if req.onRespHeader != nil {
			if action = req.onRespHeader(_resp.StatusCode, _resp.Header); action != None {
				drainAndClose(originalResponseBody)
				return
			}
		}

		// do transformation(s)
		for i := range req.transforms {
			if _resp, err = req.transforms[i].Transform(_resp); err != nil {
				resp.err = transformError(req.r.URL, err)
				drainAndClose(originalResponseBody)
				return
			}
		}
		resp._resp = _resp

		// do decode
		if err = decode(req, _resp); err == nil { // optimistic branching
			resp.data = req.expect
		} else {
			resp.err = decodingError(req.r.URL, err)
		}

		// drain up and close request body for connection reusing
		drainAndClose(originalResponseBody)
	} else {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			resp.err = requestCtxCanceledOrTimeout(req.r.URL, err)
			if req.onRequestCtxCanceledOrTimeout == nil {
				action = None
			} else {
				action = req.onRequestCtxCanceledOrTimeout()
			}

		default:
			// report connect failure to CB
			endpoint.onConnectFailure()

			// there is something wrong with the connection, should retry on next endpoint
			resp.err = connectionError(req.r.URL, err)
			action = NextEndpoint
		}
	}

	return
}
