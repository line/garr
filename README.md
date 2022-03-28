# garr

![Test](https://github.com/line/garr/actions/workflows/test.yml/badge.svg)

Collection of high performance, thread-safe, lock-free go data structures.

* [adder](./adder/README.md) - Data structure to perform highly-performant sum under high contention. Inspired by [OpenJDK LongAdder](https://openjdk.java.net/)
* [circuit-breaker](./circuit-breaker/README.md) - Data structure to implement circuit breaker pattern to detect remote service failure/alive status.
* [queue](./queue/README.md) - Queue data structure, go implementation of `JDKLinkedQueue` and `MutexLinkedQueue` from `OpenJDK`.
* [retry](./retry/README.md) - Controls backoff between attempts in a retry operation.
* [worker-pool](./worker-pool/README.md) - Worker pool implementation in go to help perform multiple tasks concurrently with a fixed-but-expandable amount of workers.

# Usage

## Getting started

```bash
go get -u github.com/line/garr
```

## Examples

Please find detailed examples in each sub-package.

### Adder

```go
package main

import (
	"fmt"
	"time"

	ga "github.com/line/garr/adder"
)

func main() {
	// or ga.DefaultAdder() which uses jdk long-adder as default
	adder := ga.NewLongAdder(ga.JDKAdderType) 

	for i := 0; i < 100; i++ {
		go func() {
			adder.Add(123)
		}()
	}

	time.Sleep(3 * time.Second)

	// get total added value
	fmt.Println(adder.Sum()) 
}
```

#### Build your own Prometheus counter with Adder

```go
package prom

import (
	ga "github.com/line/garr/adder"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// NewCounterI64 creates a new CounterI64 based on the provided prometheus.CounterOpts.
func NewCounterI64(opts prometheus.CounterOpts) CounterI64 {
	return CounterI64{counter: prometheus.NewCounter(opts)}
}

// CounterI64 is optimized Prometheus Counter for int64 value type.
type CounterI64 struct {
	val     ga.JDKAdder
	counter prometheus.Counter
}

// Value returns current value.
func (c *CounterI64) Value() int64 {
	return c.val.Sum()
}

// Reset value.
func (c *CounterI64) Reset() {
	c.val.Reset()
}

// Desc returns metric desc.
func (c *CounterI64) Desc() *prometheus.Desc {
	return c.counter.Desc()
}

// Inc by 1.
func (c *CounterI64) Inc() {
	c.val.Add(1)
}

// Add by variant.
func (c *CounterI64) Add(val int64) {
	if val > 0 {
		c.val.Add(val)
	}
}

// Write implements prometheus.Metric interface.
func (c *CounterI64) Write(out *dto.Metric) (err error) {
	if err = c.counter.Write(out); err == nil {
		value := float64(c.val.Sum())
		out.Counter.Value = &value
	}
	return
}

// Collect implements prometheus.Collector interface.
func (c *CounterI64) Collect(ch chan<- prometheus.Metric) {
	ch <- c
}

// Describe implements prometheus.Collector interface.
func (c *CounterI64) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.counter.Desc()
}
```

### Queue

```go
package main

import (
    "fmt"

    "github.com/line/garr/queue"
)

func main() {
    q := queue.DefaultQueue() // default using jdk linked queue

    // push
    q.Offer(123)

    // return head queue but not remove
    head := q.Peak()
    fmt.Println(head)

    // remove and return head queue
    polled := q.Poll()
    fmt.Println(polled)
}
```

### Circuit Breaker

```go
package main

import (
    cbreaker "github.com/line/garr/circuit-breaker"
)

func makeRequest() error {
	return nil
}

func main() {
    cb := cbreaker.NewCircuitBreakerBuilder().
                        SetTicker(cbreaker.SystemTicker).
                        SetFailureRateThreshold(validFailureRateThreshold).
                        Build()

    if cb.CanRequest() {
        err := makeRequest()
        if err != nil {
            cb.OnFailure()
        } else {
            cb.OnSuccess()
        }
    }
}
```
