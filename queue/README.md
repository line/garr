# Queue

High performance and thread-safe queue(s) for Go:
* `JDKLinkedQueue`: a lockless linked-list queue ported from OpenJDK ConcurrentLinkedQueue.
* `MutexLinkedQueue`: linked-list queue based on mutex.

# Usage

```go
package main

import (
    "go.linecorp.com/garr/queue"
)

func main() {
    q := queue.DefaultQueue() // default using jdk linked queue

    // push
    q.Offer(struct{}{})

    // remove and return head queue
    polled := q.Poll()

    // return head queue but not remove
    head := q.Peek()
}
```

# Benchmark

```bash
GO111MODULE=""
GOARCH="amd64"
GOBIN=""
GOEXE=""
GOEXPERIMENT=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOINSECURE=""
GONOPROXY=""
GONOSUMDB=""
GOOS="linux"
GOPRIVATE=""
GOPROXY="https://proxy.golang.org,direct"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOVCS=""
GOVERSION="go1.17.8"
GCCGO="gccgo"
AR="ar"
CC="gcc"
CXX="g++"
CGO_ENABLED="1"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build1560996522=/tmp/go-build -gno-record-gcc-switches"
```
```bash
goos: linux
goarch: amd64
cpu: AMD Ryzen 9 3950X 16-Core Processor            
Benchmark_MutexLinkedQueue_50P50C-32    	      2	562280198 ns/op	32044940 B/op	1000339 allocs/op
Benchmark_JDKLinkedQueue_50P50C-32      	      4	256514768 ns/op	32038624 B/op	1500368 allocs/op
Benchmark_MutexLinkedQueue_50P10C-32    	      2	510193642 ns/op	32020120 B/op	1000175 allocs/op
Benchmark_JDKLinkedQueue_50P10C-32      	      5	220084887 ns/op	32017083 B/op	1500175 allocs/op
Benchmark_MutexLinkedQueue_10P50C-32    	     13	110709024 ns/op	6408471 B/op	 200142 allocs/op
Benchmark_JDKLinkedQueue_10P50C-32      	     14	 89244324 ns/op	6406506 B/op	 300136 allocs/op
Benchmark_MutexLinkedQueue_100P-32      	      3	443367632 ns/op	64027277 B/op	2000195 allocs/op
Benchmark_JDKLinkedQueue_100P-32        	      2	648576225 ns/op	64003456 B/op	3000103 allocs/op
Benchmark_MutexLinkedQueue_100C-32      	      3	365702313 ns/op	64006032 B/op	2000117 allocs/op
Benchmark_JDKLinkedQueue_100C-32        	      4	333530039 ns/op	64003216 B/op	3000101 allocs/op
```
