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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type mockReadCloser struct {
	*bytes.Buffer
	closeState bool
}

func (m *mockReadCloser) Close() error {
	m.closeState = true
	return nil
}

type mockWriteCloser struct {
	*bytes.Buffer
	closeState bool
}

func (m *mockWriteCloser) Close() error {
	m.closeState = true
	return nil
}

type mockErrWriter struct {
}

func (m *mockErrWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("Fake error")
}

type mockWriteErrCloser struct {
}

func (m *mockWriteErrCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockWriteErrCloser) Close() error {
	return fmt.Errorf("Fake error")
}

func createMockResponse() *http.Response {
	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(createBytes())),
	}
}

func createMockJSONResponse(payload interface{}) *http.Response {
	b, _ := json.Marshal(payload)
	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(b)),
	}
}

func createBytes() []byte {
	v := make([]byte, 100)
	for i := range v {
		v[i] = byte(i)
	}
	return v
}

func createGzippedPayload() []byte {
	v := make([]byte, 10000)
	for i := range v {
		v[i] = byte(i)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write(v)
	_ = gz.Close()

	return buf.Bytes()
}

func TestDefaultCB(t *testing.T) {
	if b := defaultCircuitBreakerBuilder(); b == nil {
		t.FailNow()
	}
}

func TestDefaultBackoff(t *testing.T) {
	if b := defaultBackoff(); b == nil {
		t.FailNow()
	}
}

func TestDrainAndClose(t *testing.T) {
	r := &mockReadCloser{
		Buffer: bytes.NewBuffer(make([]byte, 100)),
	}

	drainAndClose(r)

	if r.closeState != true {
		t.FailNow()
	}
	if r.Buffer.Len() != 0 {
		t.FailNow()
	}
}

func TestInjectAndRevert(t *testing.T) {
	// create an endpoint
	rawEndpoint := &RawEndpoint{RawURL: "https://host1:9999/path"}
	endpoint, err := rawEndpoint.ToEndpoint()
	if err != nil {
		t.FailNow()
	}
	if endpoint.Path != "/path" {
		t.FailNow()
	}

	// create request
	rawReq, err := http.NewRequest(http.MethodDelete, "http://host2:9999/query", nil)
	if err != nil {
		t.FailNow()
	}

	// try to injection
	req := NewRequest(rawReq)
	original := injectTarget(endpoint, req)
	if rawReq.URL.Scheme != "https" ||
		rawReq.URL.Host != "host1:9999" ||
		rawReq.URL.Path != "/path/query" {
		t.FailNow()
	}

	// try to revert
	revert(req, original)
	if rawReq.URL.Scheme != "http" ||
		rawReq.URL.Host != "host2:9999" ||
		rawReq.URL.Path != "/query" {
		t.FailNow()
	}
}

func TestLookupPortByScheme(t *testing.T) {
	network, p, err := lookupPortByScheme("http")
	if err != nil || p != 80 || network != "tcp" {
		t.FailNow()
	}

	network, p, err = lookupPortByScheme("https")
	if err != nil || p != 443 || network != "tcp" {
		t.FailNow()
	}

	network, p, err = lookupPortByScheme("ftp")
	if err != nil || p != 21 || network != "tcp" {
		t.FailNow()
	}

	network, p, err = lookupPortByScheme("ssh")
	if err != nil || p != 22 || network != "tcp" {
		t.FailNow()
	}

	network, p, err = lookupPortByScheme("ftps")
	if err != nil || p != 990 || network != "tcp" {
		t.FailNow()
	}

	network, p, err = lookupPortByScheme("unknown")
	if err == nil || p != 0 || network != "" {
		t.FailNow()
	}
}

func validEndpoints() Endpoints {
	valids := []string{
		"https://github.com",
		"https://google.com",
		"https://golang.org",
	}
	endpoints, err := ParseFromURLs(valids)
	if err != nil {
		panic(err)
	}
	return endpoints
}
