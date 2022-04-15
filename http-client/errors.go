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
	"fmt"
	"net/url"
)

var (
	// ErrNoEndpoints indicates no endpoints avaiable.
	ErrNoEndpoints = fmt.Errorf("There is no endpoints avaiable")

	// ErrEndpointsUnavailable indicates endpoints are unavailable.
	ErrEndpointsUnavailable = fmt.Errorf("All endpoints are unavailable or all circuit-breakers of endpoints are opened")
)

type errorCategory int

const (
	// indicates no error.
	errNone errorCategory = iota

	// indicates connection error category, i.e dial, unstable/terminal state, etc.
	errConnection

	// indicates decoding error category.
	errDecoding

	// indicates transform error category.
	errTransform

	// indicates request context canceled explicitly/timeout.
	errRequestCtxCanceledOrTimeout
)

// wraps over error
type errorWrap struct {
	category errorCategory
	err      error
}

func (e *errorWrap) Error() string {
	return e.err.Error()
}

func (e *errorWrap) Unwrap() error {
	return e.err
}

func connectionError(url *url.URL, err error) error {
	return &errorWrap{
		category: errConnection,
		err:      fmt.Errorf("Connection to host:[%s] got error:[%w]", url.Host, err),
	}
}

func decodingError(url *url.URL, err error) error {
	return &errorWrap{
		category: errDecoding,
		err:      fmt.Errorf("Decoding of response body from host:[%s] got error:[%w]", url.Host, err),
	}
}

func transformError(url *url.URL, err error) error {
	return &errorWrap{
		category: errTransform,
		err:      fmt.Errorf("Transformation of response body from host:[%s] got error:[%w]", url.Host, err),
	}
}

func requestCtxCanceledOrTimeout(url *url.URL, err error) error {
	return &errorWrap{
		category: errRequestCtxCanceledOrTimeout,
		err:      fmt.Errorf("Request context canceled or timeout. Detail: Host:[%s] Error:[%w]", url.Host, err),
	}
}
