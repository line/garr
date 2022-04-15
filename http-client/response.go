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
)

// Response instruments http.Response with features.
type Response struct {
	req   *Request
	_resp *http.Response
	data  interface{}
	err   error
}

func (r *Response) reset() {
	r.req = nil
	r._resp = nil
	r.data = nil
	r.err = nil
}

// Raw returns received raw http.Response.
//
// Response body is always instrumented by the client
// and processed through transformers and decoder defined
// with belonged request.
func (r *Response) Raw() *http.Response {
	return r._resp
}

// Request returns Request that this response belongs to.
func (r *Response) Request() *Request {
	return r.req
}

// StatusCode returns response status code (if have).
// On connection failure, return -1.
func (r *Response) StatusCode() (code int) {
	if r._resp != nil {
		code = r._resp.StatusCode
	} else {
		code = -1
	}
	return
}

// Status returns response status (if have).
// On connection failure, return empty.
func (r *Response) Status() (status string) {
	if r._resp != nil {
		status = r._resp.Status
	}
	return
}

// Data returns (parsed) response data.
func (r *Response) Data() interface{} {
	return r.data
}

func (r *Response) errorCategory() errorCategory {
	if ew, ok := r.err.(*errorWrap); ok {
		return ew.category
	}
	return errNone
}

// IsDecodingError indicates decoding error.
func (r *Response) IsDecodingError() (v bool) {
	return r.errorCategory() == errDecoding
}

// IsConnectionError indicates connection error.
func (r *Response) IsConnectionError() bool {
	return r.errorCategory() == errConnection
}

// IsTransformError indicates transformation error.
func (r *Response) IsTransformError() bool {
	return r.errorCategory() == errTransform
}

// IsRequestCtxCanceledOrTimeout indicates request context canceled explicitly or timeout occured.
func (r *Response) IsRequestCtxCanceledOrTimeout() bool {
	return r.errorCategory() == errRequestCtxCanceledOrTimeout
}

// Error returns response error (in detail).
func (r *Response) Error() error {
	return r.err
}
