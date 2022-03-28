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

package adder

import (
	"runtime"
	"sync"
	"testing"
)

const (
	benchNumRoutine       = 32
	benchDelta            = 100000
	benchDeltaSingleRoute = 1000000
	wrRatio               = 5000 // 5000 write - 1 read
)

var (
	atomicAdder1     = NewLongAdder(AtomicAdderType)
	mutexAdder1      = NewLongAdder(MutexAdderType)
	jdkAdder1        = NewLongAdder(JDKAdderType)
	randomCellAdder1 = NewLongAdder(RandomCellAdderType)

	atomicAdder2     = NewLongAdder(AtomicAdderType)
	mutexAdder2      = NewLongAdder(MutexAdderType)
	jdkAdder2        = NewLongAdder(JDKAdderType)
	randomCellAdder2 = NewLongAdder(RandomCellAdderType)

	atomicAdder3     = NewLongAdder(AtomicAdderType)
	mutexAdder3      = NewLongAdder(MutexAdderType)
	jdkAdder3        = NewLongAdder(JDKAdderType)
	randomCellAdder3 = NewLongAdder(RandomCellAdderType)
)

func init() {
	// set max procs to thread contention
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func BenchmarkJDKAdderSingleRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderSingleRoutine(jdkAdder1)
	}
}

func BenchmarkAtomicAdderSingleRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderSingleRoutine(atomicAdder1)
	}
}

func BenchmarkRandomCellAdderSingleRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderSingleRoutine(randomCellAdder1)
	}
}

func BenchmarkMutexAdderSingleRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderSingleRoutine(mutexAdder1)
	}
}

func BenchmarkJDKAdderMultiRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutine(jdkAdder2)
	}
}

func BenchmarkAtomicAdderMultiRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutine(atomicAdder2)
	}
}

func BenchmarkRandomCellAdderMultiRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutine(randomCellAdder2)
	}
}

func BenchmarkMutexAdderMultiRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutine(mutexAdder2)
	}
}

func BenchmarkJDKAdderMultiRoutineMix(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutineMix(jdkAdder3)
	}
}

func BenchmarkAtomicAdderMultiRoutineMix(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutineMix(atomicAdder3)
	}
}
func BenchmarkRandomCellAdderMultiRoutineMix(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutineMix(randomCellAdder3)
	}
}

func BenchmarkMutexAdderMultiRoutineMix(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutineMix(mutexAdder3)
	}
}

func benchAdderSingleRoutine(adder LongAdder) {
	for i := 0; i < benchDeltaSingleRoute; i++ {
		adder.Add(1)
	}
}

func benchAdderMultiRoutine(adder LongAdder) {
	var wg sync.WaitGroup
	for i := 0; i < benchNumRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < benchDelta; j++ {
				adder.Add(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func benchAdderMultiRoutineMix(adder LongAdder) {
	var wg sync.WaitGroup
	for i := 0; i < benchNumRoutine; i++ {
		wg.Add(1)
		go func() {
			var sum int64
			for j := 0; j < benchDelta; j++ {
				adder.Add(1)
				if j%wrRatio == 0 {
					sum += adder.Sum()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
