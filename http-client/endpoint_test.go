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
	"testing"
	"time"
)

func TestParseEndpoints(t *testing.T) {
	_, err := ParseFromURLs([]string{})
	if err != nil {
		t.FailNow()
	}

	_, err = ParseFromURLs([]string{"://github.com"})
	if err == nil {
		t.FailNow()
	}

	_, err = ParseFromURLs([]string{"phantom://github.com"})
	if err == nil {
		t.FailNow()
	}

	endpoints, err := ParseFromURLs([]string{"https://github.com"})
	if err != nil {
		t.FailNow()
	}
	if len(endpoints) != 1 {
		t.FailNow()
	}
	if endpoints[0].Scheme != "https" ||
		endpoints[0].Host != "github.com:443" ||
		endpoints[0].Port() != "443" {
		t.FailNow()
	}
}

func TestEndpointsNormalization(t *testing.T) {
	endpoints, err := ParseFromURLs([]string{"https://github.com", "https://google.com"})
	if err != nil {
		t.FailNow()
	}
	if len(endpoints) != 2 {
		t.FailNow()
	}

	// inject wrong scheme
	endpoints[0].Scheme = "phantom"
	if endpoints.normalize() == nil {
		t.FailNow()
	}
}

func TestEndpointsClone(t *testing.T) {
	endpoints, _ := ParseFromURLs([]string{"https://github.com", "https://google.com"})
	cloned := endpoints.Clone()
	if !endpoints.Equal(cloned) {
		t.FailNow()
	}
}

func TestEndpointsEqual(t *testing.T) {
	endpoints1, _ := ParseFromURLs([]string{"https://github.com", "https://google.com"})
	endpoints2, _ := ParseFromURLs([]string{"https://linxGnu@github.com", "https://google.com"})
	if endpoints1.Equal(endpoints2) {
		t.FailNow()
	}

	endpoints1, _ = ParseFromURLs([]string{"https://test@github.com", "https://google.com"})
	endpoints2, _ = ParseFromURLs([]string{"https://linxGnu@github.com", "https://google.com"})
	if endpoints1.Equal(endpoints2) {
		t.FailNow()
	}

	endpoints1, _ = ParseFromURLs([]string{"https://test@github.com", "https://google.com"})
	endpoints2, _ = ParseFromURLs([]string{"https://github.com", "https://google.com"})
	if endpoints1.Equal(endpoints2) {
		t.FailNow()
	}

	endpoints1, _ = ParseFromURLs([]string{"https://test@github.com", "https://google.com"})
	endpoints2, _ = ParseFromURLs([]string{"https://test@github.com", "https://google.com"})
	if !endpoints1.Equal(endpoints2) {
		t.FailNow()
	}
}

func TestEndpointDial(t *testing.T) {
	dialEndpoint := func(rawURL string, validCase bool) {
		r := &RawEndpoint{RawURL: rawURL}
		e, err := r.ToEndpoint()
		if err != nil {
			t.FailNow()
		}
		if e.Dial(0) != validCase {
			t.FailNow()
		}
		if e.Dial(100*time.Millisecond) != validCase {
			t.FailNow()
		}
	}

	// valid cases
	valids := []string{
		"https://github.com",
		"https://google.com",
	}
	for i := range valids {
		dialEndpoint(valids[i], true)
	}

	// invalid cases
	invalids := []string{
		"https://google1.com",
		"http://127.0.0.1:9578",
	}
	for i := range invalids {
		dialEndpoint(invalids[i], false)
	}
}
