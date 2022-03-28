# Retry Backoff

Inspired by [line/armeria](https://github.com/line/armeria) retry-backoff package, controls backoff between attempts in a single retry operation. 

Supporting backoff strategies:
- Exponential
- Fixed
- Jitter
- Random

# Usage

```go
package main

import (
    "github.com/line/garr/retry"
)

func main() {
	builder := retry.NewBackoffBuilder()
	if _, err := builder.Build(); err == nil {
		panic(err)
	}

	// build a fixed attempt backoff
	// remember that base backoff is mandatory
	builder = NewBackoffBuilder().BaseBackoffSpec("fixed=456")

	backoff, err := builder.Build()
	if err != nil || backoff == nil {
		panic(err)
	}
    
	for i := 0; i < 10000; i++ {
		if backoff.NextDelayMillis(i) != 456 {
			t.FailNow()
		}
	}
    
	// backoff with jitter
	builder = retry.NewBackoffBuilder().
			BaseBackoff(fixedBackoff).
			WithLimit(5).
			WithJitter(0.9).
			WithJitterBound(0.9, 1.2)
	backoff, err = builder.Build()
    
	if backoff.NextDelayMillis(6) < 0 {
		// should stop retrying
		stopRetrying()
	} else {
		panic("wrong backoff")
	}
}
```
