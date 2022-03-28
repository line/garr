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
	numRoutine = 4
	delta      = 300000
)

func TestDefaultAdder(t *testing.T) {
	if _, ok := DefaultAdder().(*JDKAdder); !ok {
		t.Error("Invalid default adder expected")
	}

	if _, ok := DefaultFloat64Adder().(*JDKF64Adder); !ok {
		t.Error("Invalid default adder expected")
	}
}

func TestMaxCellsNormalize(t *testing.T) {
	if normalizeMaxCell(1<<12) != 1<<11 || normalizeMaxCell(32) != 64 || normalizeMaxCell(100) != 100 {
		t.FailNow()
	}
}

func testAdderNotRaceInc(t *testing.T, ty Type) {
	adder := NewLongAdder(ty)

	for i := 0; i < delta; i++ {
		adder.Inc()
	}

	tmp := int64(delta)
	if adder.Sum() != tmp || adder.SumAndReset() != tmp || adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testAdderRaceInc(t *testing.T, ty Type) {
	adder := NewLongAdder(ty)

	var wg sync.WaitGroup
	for i := 0; i < numRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < delta; j++ {
				adder.Inc()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	tmp := int64(delta) * int64(numRoutine)
	if adder.Sum() != tmp || adder.SumAndReset() != tmp || adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}

	// try to store
	if adder.Store(12341); adder.Sum() != 12341 {
		t.Errorf("Store(%d) logic is wrong", ty)
	}
}

func testAdderNotRaceDec(t *testing.T, ty Type) {
	adder := NewLongAdder(ty)

	for i := 0; i < delta; i++ {
		adder.Dec()
	}

	if adder.Sum() != -int64(delta) {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}

	adder.Reset()
	if adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testAdderRaceDec(t *testing.T, ty Type) {
	adder := NewLongAdder(ty)

	var wg sync.WaitGroup
	for i := 0; i < numRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < delta; j++ {
				adder.Dec()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if adder.Sum() != -int64(delta)*int64(numRoutine) {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}

	adder.Reset()
	if adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testAdderNotRaceAdd(t *testing.T, ty Type) {
	adder := NewLongAdder(ty)

	for i := 0; i < delta; i++ {
		adder.Add(int64(i))
	}

	tmp := int64(delta)
	if adder.Sum() != tmp*(tmp-1)/2 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testAdderRaceAdd(t *testing.T, ty Type) {
	adder := NewLongAdder(ty)

	var wg sync.WaitGroup
	for i := 0; i < numRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < delta; j++ {
				adder.Add(int64(j))
			}
			wg.Done()
		}()
	}
	wg.Wait()

	tmp := int64(delta)
	if adder.Sum() != (tmp*(tmp-1)/2)*int64(numRoutine) {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testF64AdderNotRaceInc(t *testing.T, ty Type) {
	adder := NewFloat64Adder(ty)

	for i := 0; i < delta; i++ {
		adder.Inc()
	}

	tmp := float64(delta)
	if adder.Sum() != tmp || adder.SumAndReset() != tmp || adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}

}

func testF64AdderRaceInc(t *testing.T, ty Type) {
	adder := NewFloat64Adder(ty)

	var wg sync.WaitGroup
	for i := 0; i < numRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < delta; j++ {
				adder.Inc()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	tmp := float64(delta) * float64(numRoutine)
	if adder.Sum() != tmp || adder.SumAndReset() != tmp || adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}

	// try to store
	if adder.Store(12341); adder.Sum() != 12341 {
		t.Errorf("Store(%d) logic is wrong", ty)
	}
}

func testF64AdderNotRaceDec(t *testing.T, ty Type) {
	adder := NewFloat64Adder(ty)

	for i := 0; i < delta; i++ {
		adder.Dec()
	}

	if adder.Sum() != -float64(delta) {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}

	adder.Reset()
	if adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testF64AdderRaceDec(t *testing.T, ty Type) {
	adder := NewFloat64Adder(ty)

	var wg sync.WaitGroup
	for i := 0; i < numRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < delta; j++ {
				adder.Dec()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if adder.Sum() != -float64(delta)*float64(numRoutine) {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}

	adder.Reset()
	if adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testF64AdderNotRaceAdd(t *testing.T, ty Type) {
	adder := NewFloat64Adder(ty)

	for i := 0; i < delta; i++ {
		adder.Add(float64(i))
	}

	tmp := float64(delta)
	if adder.Sum() != tmp*(tmp-1)/2 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testF64AdderRaceAdd(t *testing.T, ty Type) {
	adder := NewFloat64Adder(ty)

	var wg sync.WaitGroup
	for i := 0; i < numRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < delta; j++ {
				adder.Add(float64(j))
			}
			wg.Done()
		}()
	}
	wg.Wait()

	tmp := float64(delta)
	if adder.Sum() != (tmp*(tmp-1)/2)*float64(numRoutine) {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}
