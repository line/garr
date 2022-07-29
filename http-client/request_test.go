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
	"testing"
)

func TestNewRequest(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodDelete, "http://host2:9999/query", nil)
	if err != nil {
		t.FailNow()
	}

	expect := make(map[int]int)
	req := NewRequest(rawReq,
		WithTransform(Limiter(128<<10)),
		WithDecoder(JSON),
		WithTransforms([]Transformer{Limiter(128 << 8)}),
		WithExpect(&expect),
		OnResponseHeader(OnStatus5xx),
	)
	if err != nil {
		t.FailNow()
	}
	if req == nil {
		t.FailNow()
	}
	if rawReq.Host != "" {
		t.FailNow()
	}
	if &expect != req.expect {
		t.FailNow()
	}
	if req.onRespHeader == nil {
		t.FailNow()
	}

	if len(req.transforms) != 2 {
		t.FailNow()
	}
	l, ok := req.transforms[0].(*limiter)
	if !ok || l.n != 128<<10 {
		t.FailNow()
	}
	l, ok = req.transforms[1].(*limiter)
	if !ok || l.n != 128<<8 {
		t.FailNow()
	}
}
