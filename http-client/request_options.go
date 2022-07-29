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

// RequestOption is setter for request's option.
type RequestOption func(r *Request)

// WithTransform adds a transformation for received response.
// Response would be transformed in order of RequestOption arguments.
//
// Default: empty transformers
func WithTransform(t Transformer) RequestOption {
	return func(r *Request) {
		r.transforms = append(r.transforms, t)
	}
}

// WithTransforms adds transformation(s) for received response.
// Response would be transformed in order of RequestOption arguments.
//
// Default: empty transformers
func WithTransforms(t []Transformer) RequestOption {
	return func(r *Request) {
		r.transforms = append(r.transforms, t...)
	}
}

// WithDecoder attaches a decoder for received response.
//
// Default: nil
func WithDecoder(d Decoder) RequestOption {
	return func(r *Request) {
		r.decoder = d
	}
}

// WithExpect indicates expectation of decoded response data.
// This option is not mandatory.
//
// Default: nil
func WithExpect(expect interface{}) RequestOption {
	return func(r *Request) {
		r.expect = expect
	}
}

// OnResponseHeader specifies action to take according to received response header.
// This option is not mandatory.
//
// Default: nil
func OnResponseHeader(f func(statusCode int, header http.Header) EndpointAction) RequestOption {
	return func(r *Request) {
		r.onRespHeader = f
	}
}

// OnRequestCtxCanceledOrTimeout specifies action to take in case request context is canceled
// or timeout.
//
// Default: nil
func OnRequestCtxCanceledOrTimeout(f func() EndpointAction) RequestOption {
	return func(r *Request) {
		r.onRequestCtxCanceledOrTimeout = f
	}
}
