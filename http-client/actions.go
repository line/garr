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

import "net/http"

// EndpointAction specifies in which cases a request should be passed
// to the next endpoint or retrying on current endpoint.
//
// Similar to: http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_next_upstream
type EndpointAction byte

const (
	// None indicates that there is no need for taking any actions.
	None EndpointAction = iota

	// NextEndpoint indicates that client should retry request on next endpoint.
	NextEndpoint

	// Retrying indicates that client could retry request on same endpoint.
	Retrying
)

// OnStatus5xx is builtin decider which judge on 5xx status code.
func OnStatus5xx(statusCode int, _ http.Header) (action EndpointAction) {
	switch statusCode {
	case http.StatusBadGateway:
		action = Retrying
	case http.StatusInternalServerError, http.StatusServiceUnavailable:
		action = NextEndpoint
	default:
		action = None
	}
	return
}
