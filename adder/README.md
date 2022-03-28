# Adder 

Thread-safe, high performance, contention-awareness `LongAdder` and `DoubleAdder` for Go, inspired by OpenJDK9.
Beside JDK-based `LongAdder` and `DoubleAdder`, the library also includes other adders for various usage.

# Usage

## JDKAdder (recommended)

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

## RandomCellAdder

* A `LongAdder` with simple strategy by preallocating atomic cell and select random cell to update.
* Slower than JDK LongAdder but faster than AtomicAdder on contention.
* Consume ~1KB to store cells.

```go
adder := ga.NewLongAdder(ga.RandomCellAdderType)
```

## AtomicAdder

* A `LongAdder` based on atomic variable. All routines share this variable.

```go
adder := ga.NewLongAdder(ga.AtomicAdderType)
```

## MutexAdder

* A `LongAdder` based on mutex. All routines share same value and mutex.

```go
adder := ga.NewLongAdder(ga.MutexAdderType)
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
BenchmarkJDKF64AdderSingleRoutine-32                 212           5778275 ns/op               0 B/op          0 allocs/op
BenchmarkAtomicF64AdderSingleRoutine-32              240           4862050 ns/op               0 B/op          0 allocs/op
BenchmarkJDKF64AdderMultiRoutine-32                   78          16965361 ns/op           10984 B/op         55 allocs/op
BenchmarkAtomicF64AdderMultiRoutine-32                18          64129718 ns/op            1042 B/op         33 allocs/op
BenchmarkJDKF64AdderMultiRoutineMix-32                61          18043656 ns/op            1168 B/op         33 allocs/op
BenchmarkAtomicF64AdderMultiRoutineMix-32             16          65600563 ns/op            1046 B/op         33 allocs/op
BenchmarkJDKAdderSingleRoutine-32                    223           5410781 ns/op               0 B/op          0 allocs/op
BenchmarkAtomicAdderSingleRoutine-32                 272           4485851 ns/op               0 B/op          0 allocs/op
BenchmarkRandomCellAdderSingleRoutine-32              75          16295095 ns/op              54 B/op          0 allocs/op
BenchmarkMutexAdderSingleRoutine-32                   63          18643397 ns/op               0 B/op          0 allocs/op
BenchmarkJDKAdderMultiRoutine-32                      96          16589125 ns/op            1090 B/op         33 allocs/op
BenchmarkAtomicAdderMultiRoutine-32                   49          24409391 ns/op            1041 B/op         33 allocs/op
BenchmarkRandomCellAdderMultiRoutine-32               64          20095099 ns/op            1161 B/op         33 allocs/op
BenchmarkMutexAdderMultiRoutine-32                     4         318551677 ns/op            1808 B/op         41 allocs/op
BenchmarkJDKAdderMultiRoutineMix-32                   72          16713561 ns/op            1123 B/op         33 allocs/op
BenchmarkAtomicAdderMultiRoutineMix-32                46          25109417 ns/op            1040 B/op         33 allocs/op
BenchmarkRandomCellAdderMultiRoutineMix-32            68          20661803 ns/op            1153 B/op         33 allocs/op
BenchmarkMutexAdderMultiRoutineMix-32                  4         417323064 ns/op            1784 B/op         40 allocs/op
```
