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
	"bytes"
	"testing"
)

func TestDecoders(t *testing.T) {
	// expect is write closer
	req, resp := &Request{
		expect: &mockWriteCloser{Buffer: bytes.NewBuffer(make([]byte, 0, 50))},
	}, createMockResponse()
	err := decode(req, resp)
	if err != nil {
		t.FailNow()
	}
	if !req.expect.(*mockWriteCloser).closeState {
		t.FailNow()
	}

	payload := req.expect.(*mockWriteCloser).Bytes()
	if len(payload) != 100 {
		t.FailNow()
	}
	for i := range payload {
		if payload[i] != byte(i) {
			t.FailNow()
		}
	}

	// expect is an object and need json decoder
	v := make(map[string]string)
	req, resp = &Request{
		expect:  &v,
		decoder: JSON,
	}, createMockJSONResponse(map[string]string{"a": "B", "c": "D"})
	err = decode(req, resp)
	if err != nil {
		t.FailNow()
	}
	if !equalMapString(v, map[string]string{"a": "B", "c": "D"}) {
		t.FailNow()
	}

	// writer but error
	req, resp = &Request{
		expect: &mockErrWriter{},
	}, createMockResponse()
	if err = decode(req, resp); err == nil {
		t.FailNow()
	}

	// write ok, but close error
	req, resp = &Request{
		expect: &mockWriteErrCloser{},
	}, createMockResponse()
	if err = decode(req, resp); err == nil {
		t.FailNow()
	}
}

func equalMapString(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}

	for k, v1 := range map1 {
		v2, exist := map2[k]
		if !exist || v1 != v2 {
			return false
		}
	}

	return true
}
