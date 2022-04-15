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

// Request instruments http.Request with additional properties/hooks
// that help performing RoundTrip action.
type Request struct {
	// original http request
	r *http.Request

	// transformations chain for response
	transforms []Transformer

	// response body decoder
	decoder Decoder

	// action on received response header
	onRespHeader func(statusCode int, header http.Header) EndpointAction

	// action when request context canceled or timeout
	onRequestCtxCanceledOrTimeout func() EndpointAction

	// expect output
	expect interface{}
}

// NewRequest creates new Request.
func NewRequest(r *http.Request, opts ...RequestOption) *Request {
	v := &Request{
		r: r,
	}

	// make empty to trigger using url.URL
	r.Host = ""

	for i := range opts {
		opts[i](v)
	}

	return v
}
