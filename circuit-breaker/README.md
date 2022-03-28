# circuit breaker

Inspired by [line/armeria](https://github.com/line/armeria) circuit breaker. This helps tracks success/failure rate of a request/operation (e.g. requests to a remote service) in a configurable time window and manage circuit breaker state for you.

* `CLOSED` Failure rate is below threshold, allow all requests to pass.
* `OPEN` Failure rate is above threshold, reject all requests.
* `HALF_OPEN` Allow a single trial request to pass to check if circuit breaker should be kept opened or can be closed.

# Usage

## Quick Start

```go
package main

import (
    cbreaker "go.linecorp.com/garr/circuit-breaker"
)

func main() {
    // Initialize circuit breaker once in your setup code
    cb := cbreaker.NewCircuitBreakerBuilder().Build()

    // ...

    // Safeguard your request by checking against circuit breaker and report success/failure
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

## Customizing Circuit Breaker Parameter

Circuit breaker builder comes with set of default parameters that you can customize when building the circuit breaker.

```go
package main

import (
    cbreaker "go.linecorp.com/garr/circuit-breaker"
)

func main() {
    builder := cbreaker.NewCircuitBreakerBuilder()

    // Set unique namespace for the circuit breaker
    // useful to differentiate multiple circuit breakers in logs
    builder.Name(&Name{
        Namespace: "my service",
        Subsystem: "API",
        Name: "external service circuit breaker",
    })

    // Set a custom ticker.
    // Default to cbreaker.SystemTicker
    builder.SetTicker(cbreaker.SystemTicker)

    // Set counter sliding window.
    // Requests are accumulated and counted to determine success/failure rate within this sliding window.
    // Default to 20 seconds.
    builder.SetCounterSlidingWindow(20 * time.Second)

    // Set counter update interval.
    // Request counter is updated every this interval.
    // Default to 1 second.
    builder.SetCounterUpdateInterval(1 * time.Second)

    // Set minimum request threshold.
    // Only checks against failure rate threshold if number of requests are above this threshold within a sliding window.
    // Default to 10.
    builder.SetMinimumRequestThreshold(10)

    // Set failure rate threshold.
    // If circuit breaker is closed and failure rate surpasses this threshold, transition state to opened.
    // Default to 0.8
    builder.SetFailureRateThreshold(0.8)

    // Set circuit open window.
    // After circuit breaker is opened, wait for this open window to pass, then transition state to half-opened.
    // Default to 10 seconds.
    builder.SetCircuitOpenWindow(10 * time.Second)

    // Set trial request interval.
    // After circuit breaker is half-opened, allow a single trial request to pass.
    // If that trial request doesn't report success/failure after this interval, allow another trial request to pass.
    // Default to 3 seconds.
    builder.SetTrialRequestInterval(3 * time.Second)

    // Build the circuit breaker
    cb := builder.Build()
}
```

## Listening to circuit breaker events

You can make a custom listener and hook it into the circuit breaker so that the listeners get invoked when certain events happen. This listener can be added when building the circuit breaker.


```go
package main

import (
    cbreaker "go.linecorp.com/garr/circuit-breaker"
)

func main() {
    builder := cbreaker.NewCircuitBreakerBuilder()

    // Add a custom listener
    builder.AddListener(&dummyCircuitBreakerListener{})

    // Build the circuit breaker
    cb := builder.Build()
}

type dummyCircuitBreakerListener struct{}

// OnStateChanged invoked when the circuit state is changed.
func (d *dummyCircuitBreakerListener) OnStateChanged(cb CircuitBreaker, state CircuitState) (err error) {
	return
}

// OnEventCountUpdated invoked when the circuit breaker's internal EventCount is updated.
func (d *dummyCircuitBreakerListener) OnEventCountUpdated(cb CircuitBreaker, eventCount *EventCount) (err error) {
	return
}

// OnRequestRejected invoked when the circuit breaker rejects a request.
func (d *dummyCircuitBreakerListener) OnRequestRejected(cb CircuitBreaker) (err error) {
	return
}

// Stop listening
func (d *dummyCircuitBreakerListener) Stop() {
}
```
