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
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestResponseReset(t *testing.T) {
	r := &Response{
		req:   &Request{},
		_resp: &http.Response{},
		data:  123,
		err:   fmt.Errorf("Fake error"),
	}
	r.reset()
	if r.req != nil {
		t.FailNow()
	}
	if r._resp != nil {
		t.FailNow()
	}
	if r.data != nil {
		t.FailNow()
	}
	if r.err != nil {
		t.FailNow()
	}
}

func TestResponseMethods(t *testing.T) {
	req := &Request{}
	_resp := &http.Response{StatusCode: 404, Status: "fake"}

	r := &Response{}
	if r.StatusCode() != -1 {
		t.FailNow()
	}

	r = &Response{
		req:   req,
		_resp: _resp,
		data:  123,
		err:   fmt.Errorf("Mock error"),
	}
	if r.Raw() != _resp {
		t.FailNow()
	}
	if r.Request() != req {
		t.FailNow()
	}
	if r.StatusCode() != 404 {
		t.FailNow()
	}
	if r.Status() != "fake" {
		t.FailNow()
	}
	if r.Data().(int) != 123 {
		t.FailNow()
	}
	if r.errorCategory() != errNone {
		t.FailNow()
	}
	if r.Error().Error() != "Mock error" {
		t.FailNow()
	}
	if r.IsConnectionError() != false {
		t.FailNow()
	}
	if r.IsDecodingError() != false {
		t.FailNow()
	}
	if r.IsTransformError() != false {
		t.FailNow()
	}

	// inject decoding error
	r.err = &errorWrap{
		category: errDecoding,
		err:      fmt.Errorf("Decoding ERROR"),
	}
	if r.IsConnectionError() != false {
		t.FailNow()
	}
	if r.IsDecodingError() != true {
		t.FailNow()
	}
	if r.IsTransformError() != false {
		t.FailNow()
	}
	if r.Error().Error() != "Decoding ERROR" {
		t.FailNow()
	}

	// inject connection error
	r.err = &errorWrap{
		category: errConnection,
		err:      fmt.Errorf("Connection ERROR"),
	}
	if r.IsConnectionError() != true {
		t.FailNow()
	}
	if r.IsDecodingError() != false {
		t.FailNow()
	}
	if r.IsTransformError() != false {
		t.FailNow()
	}
	if r.Error().Error() != "Connection ERROR" {
		t.FailNow()
	}

	// inject transform error
	er := fmt.Errorf("Transform ERROR")
	r.err = &errorWrap{
		category: errTransform,
		err:      er,
	}
	if r.IsConnectionError() != false {
		t.FailNow()
	}
	if r.IsDecodingError() != false {
		t.FailNow()
	}
	if r.IsTransformError() != true {
		t.FailNow()
	}
	if !errors.Is(r.Error(), er) {
		t.FailNow()
	}
}
