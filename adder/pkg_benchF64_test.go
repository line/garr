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
	"sync"
	"testing"
)

var (
	atomicF64Adder1 = NewFloat64Adder(AtomicF64AdderType)
	jdkF64Adder1    = NewFloat64Adder(JDKF64AdderType)

	atomicF64Adder2 = NewFloat64Adder(AtomicF64AdderType)
	jdkF64Adder2    = NewFloat64Adder(JDKF64AdderType)

	atomicF64Adder3 = NewFloat64Adder(AtomicF64AdderType)
	jdkF64Adder3    = NewFloat64Adder(JDKF64AdderType)
)

func BenchmarkJDKF64AdderSingleRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchF64AdderSingleRoutine(jdkF64Adder1)
	}
}

func BenchmarkAtomicF64AdderSingleRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchF64AdderSingleRoutine(atomicF64Adder1)
	}
}

func BenchmarkJDKF64AdderMultiRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchF64AdderMultiRoutine(jdkF64Adder2)
	}
}

func BenchmarkAtomicF64AdderMultiRoutine(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchF64AdderMultiRoutine(atomicF64Adder2)
	}
}

func BenchmarkJDKF64AdderMultiRoutineMix(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchF64AdderMultiRoutineMix(jdkF64Adder3)
	}
}

func BenchmarkAtomicF64AdderMultiRoutineMix(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchF64AdderMultiRoutineMix(atomicF64Adder3)
	}
}

func benchF64AdderSingleRoutine(adder Float64Adder) {
	for i := 0; i < benchDeltaSingleRoute; i++ {
		adder.Add(1.1)
	}
}

func benchF64AdderMultiRoutine(adder Float64Adder) {
	var wg sync.WaitGroup
	for i := 0; i < benchNumRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < benchDelta; j++ {
				adder.Add(1.1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func benchF64AdderMultiRoutineMix(adder Float64Adder) {
	var wg sync.WaitGroup
	for i := 0; i < benchNumRoutine; i++ {
		wg.Add(1)
		go func() {
			var sum float64
			for j := 0; j < benchDelta; j++ {
				adder.Add(1.1)
				if j%wrRatio == 0 {
					sum += adder.Sum()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
