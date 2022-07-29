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
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	numReq          = 10000
	concurrency     = 8
	initTimeout     = 500 * time.Millisecond
	idleConnPerHost = 20
)

func init() {
	runtime.GOMAXPROCS(numCPU << 3)
}

func TestClientWithInvalidParams(t *testing.T) {
	_, err := NewClient(100, nil)
	if err == nil {
		t.FailNow()
	}

	client := &Client{}
	resp := client.Do(nil)
	if resp.err != ErrNoEndpoints {
		t.FailNow()
	}

	client = &Client{}
	client.lb.Store(defaultLBBuilder().Build())
	resp = client.Do(nil)
	if resp.err != ErrNoEndpoints {
		t.FailNow()
	}
}

func TestClient(t *testing.T) {
	testClient(t, func(endpoints Endpoints) (*Client, error) {
		return NewClient(initTimeout, endpoints,
			WithTransport(&http.Transport{
				MaxIdleConnsPerHost: idleConnPerHost,
			}),
		)
	}, doRequests)
}

func TestClientWithFailedServer(t *testing.T) {
	testClient(t, func(endpoints Endpoints) (*Client, error) {
		return NewClient(initTimeout, endpoints,
			WithTransport(&http.Transport{
				MaxIdleConnsPerHost: idleConnPerHost,
			}),
		)
	}, doRequestsWithServerFailure)
}

func TestClientWithFailedServers(t *testing.T) {
	testClient(t, func(endpoints Endpoints) (*Client, error) {
		return NewClient(initTimeout, endpoints,
			WithTransport(&http.Transport{
				MaxIdleConnsPerHost: idleConnPerHost,
			}),
			WithHealthChecker(2*time.Second, 100*time.Millisecond),
		)
	}, doRequestsWithAllServersFailure)
}

func TestClientWithRequestTimeout(t *testing.T) {
	testClient(t, func(endpoints Endpoints) (*Client, error) {
		return NewClient(initTimeout, endpoints,
			WithTransport(&http.Transport{
				MaxIdleConnsPerHost: idleConnPerHost,
			}),
		)
	}, doRequestsWithTimeout)

	testClient(t, func(endpoints Endpoints) (*Client, error) {
		return NewClient(initTimeout, endpoints,
			WithTransport(&http.Transport{
				MaxIdleConnsPerHost: idleConnPerHost,
			}),
		)
	}, doRequestsWithTimeoutAndMockAction)
}

func testClient(t *testing.T, f func(Endpoints) (*Client, error), exec func([]*http.Server, *Client)) {
	addresses, servers := newServers(3)
	defer stopServers(servers)

	// setup endpoints
	endpoints, err := ParseFromURLs(addresses)
	if err != nil {
		t.FailNow()
	}

	client, err := f(endpoints)
	if err != nil {
		t.FailNow()
	}
	if client == nil {
		t.FailNow()
	}

	start := time.Now()
	exec(servers, client)
	t.Log("Execution time:", time.Since(start).Seconds())

	err = client.Close()
	if err != nil {
		t.FailNow()
	}
}

func doRequests(servers []*http.Server, client *Client) {
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numReq; j++ {
				req, err := http.NewRequest(http.MethodGet, "/json", nil)
				if err != nil {
					panic(err)
				}

				expect := make(map[string]string)
				r := NewRequest(req, WithDecoder(JSON), WithExpect(&expect))

				resp := client.Do(r)
				if resp.StatusCode() != http.StatusOK {
					panic(resp.Error())
				}
			}
		}()
	}
	wg.Wait()
}

func doRequestsWithServerFailure(servers []*http.Server, client *Client) {
	go func() {
		time.Sleep(time.Second)
		stopServer(servers[0])
	}()

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numReq; j++ {
				req, err := http.NewRequest(http.MethodGet, "/json", nil)
				if err != nil {
					panic(err)
				}

				expect := make(map[string]string)
				r := NewRequest(req, WithDecoder(JSON), WithExpect(&expect))

				resp := client.Do(r)
				if resp.StatusCode() != http.StatusOK {
					panic(resp.Error())
				}
			}
		}()
	}
	wg.Wait()
}

func doRequestsWithAllServersFailure(servers []*http.Server, client *Client) {
	go func() {
		time.Sleep(time.Second)
		stopServers(servers)
	}()

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numReq; j++ {
				req, _ := http.NewRequest(http.MethodGet, "/json", nil)
				expect := make(map[string]string)
				r := NewRequest(req, WithDecoder(JSON), WithExpect(&expect))
				resp := client.Do(r)
				if resp.StatusCode() != http.StatusOK && resp.StatusCode() != -1 {
					panic(resp.StatusCode())
				}
			}
		}()
	}
	wg.Wait()
}

func doRequestsWithTimeout(servers []*http.Server, client *Client) {
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numReq; j++ {
				req, err := http.NewRequest(http.MethodGet, "/json", nil)
				if err != nil {
					panic(err)
				}

				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Nanosecond)
				req = req.WithContext(ctx)

				expect := make(map[string]string)
				r := NewRequest(req, WithDecoder(JSON), WithExpect(&expect))
				cancel()

				resp := client.Do(r)
				if !resp.IsRequestCtxCanceledOrTimeout() {
					panic(resp.Error())
				}
			}
		}()
	}
	wg.Wait()
}

func doRequestsWithTimeoutAndMockAction(servers []*http.Server, client *Client) {
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numReq; j++ {
				req, err := http.NewRequest(http.MethodGet, "/json", nil)
				if err != nil {
					panic(err)
				}

				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Nanosecond)
				req = req.WithContext(ctx)

				expect := make(map[string]string)
				r := NewRequest(req, WithDecoder(JSON), WithExpect(&expect), OnRequestCtxCanceledOrTimeout(func() EndpointAction {
					return None
				}))
				cancel()

				resp := client.Do(r)
				if !resp.IsRequestCtxCanceledOrTimeout() {
					panic(resp.Error())
				}
			}
		}()
	}
	wg.Wait()
}

var (
	port int32 = 19909
)

func newServer() (addr string, server *http.Server) {
	mux := http.NewServeMux()
	mux.Handle("/json", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"hello":"world"}`))
	}))

	port := atomic.AddInt32(&port, 1)
	addr = fmt.Sprintf(":%d", port)
	server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	addr = fmt.Sprintf("http://127.0.0.1:%d", port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	return
}

func newServers(num int) (addrs []string, servers []*http.Server) {
	for i := 0; i < num; i++ {
		addr, server := newServer()
		addrs = append(addrs, addr)
		servers = append(servers, server)
	}
	return
}

func stopServers(servers []*http.Server) {
	for i := range servers {
		stopServer(servers[i])
	}
}

func stopServer(server *http.Server) {
	if err := server.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}
