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
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestTransformer(t *testing.T) {
	v := NewTransformer(func(h *http.Response) (*http.Response, error) {
		rd, err := gzip.NewReader(h.Body)
		if err != nil {
			return nil, err
		}
		h.Body = ioutil.NopCloser(rd)
		return h, nil
	})

	resp, err := v.Transform(&http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(createGzippedPayload())),
	})
	if err != nil {
		t.FailNow()
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.FailNow()
	}
	if len(payload) != 10000 {
		t.FailNow()
	}
	for i := range payload {
		if payload[i] != byte(i) {
			t.FailNow()
		}
	}
}
