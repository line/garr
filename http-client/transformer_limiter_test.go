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

func TestLimiter(t *testing.T) {
	l := Limiter(1024)

	resp, err := l.Transform(&http.Response{
		ContentLength: 1024,
	})
	if err != nil {
		t.FailNow()
	}
	if resp.Body == nil {
		t.FailNow()
	}

	resp, err = l.Transform(&http.Response{
		ContentLength: 4096,
	})
	if err == nil {
		t.FailNow()
	}
	if resp != nil {
		t.FailNow()
	}
}
