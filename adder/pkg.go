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

// Package adder contains a collection of thread-safe, concurrent data structures for reading and writing numeric i64/f64 counter,
// inspired by OpenJDK9 LongAdder.
//
// Beside JDKAdder, ported version of OpenJDK9 LongAdder, package also provides other alternatives for various use cases.
package adder

// Type of LongAdder.
type Type byte

const (
	// JDKAdderType is type for JDK-based LongAdder.
	JDKAdderType Type = iota
	// RandomCellAdderType is type for RandomCellAdder.
	RandomCellAdderType
	// AtomicAdderType is type for atomic-based adder.
	AtomicAdderType
	// MutexAdderType is type for MutexAdder.
	MutexAdderType
	// JDKF64AdderType is type for JDK-based DoubleAdder.
	JDKF64AdderType
	// AtomicF64AdderType is type for atomic-based float64 adder.
	AtomicF64AdderType
)

// LongAdder interface.
type LongAdder interface {
	Add(x int64)
	Inc()
	Dec()
	Sum() int64
	Reset()
	SumAndReset() int64
	Store(v int64)
}

// Float64Adder interface.
type Float64Adder interface {
	Add(x float64)
	Inc()
	Dec()
	Sum() float64
	Reset()
	SumAndReset() float64
	Store(v float64)
}

// DefaultAdder returns jdk long adder.
func DefaultAdder() LongAdder {
	return NewJDKAdder()
}

// DefaultFloat64Adder returns jdk f64 adder.
func DefaultFloat64Adder() Float64Adder {
	return NewJDKF64Adder()
}

// NewLongAdder create new long adder upon type.
func NewLongAdder(t Type) LongAdder {
	switch t {
	case MutexAdderType:
		return NewMutexAdder()
	case AtomicAdderType:
		return NewAtomicAdder()
	case RandomCellAdderType:
		return NewRandomCellAdder()
	default:
		return NewJDKAdder()
	}
}

// NewFloat64Adder create new float64 adder upon type.
func NewFloat64Adder(t Type) Float64Adder {
	switch t {
	case AtomicF64AdderType:
		return NewAtomicF64Adder()
	default:
		return NewJDKF64Adder()
	}
}

// longBinaryOperator represents an operation upon two int64-valued operands and producing an
// int64-valued result.
type longBinaryOperator interface {
	Apply(left, right int64) int64
}

// floatBinaryOperator represents an operation upon two float64-valued operands and producing an
// float64-valued result.
type floatBinaryOperator interface {
	Apply(left, right float64) float64
}
