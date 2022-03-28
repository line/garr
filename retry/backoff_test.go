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

package retry

import (
	"testing"
)

func TestBackoffBuilder(t *testing.T) {
	builder := NewBackoffBuilder()
	if _, err := builder.Build(); err == nil {
		t.Fatal()
	}
}

func TestBuilderFixedBackoff(t *testing.T) {
	builder := NewBackoffBuilder().BaseBackoffSpec("fixed=456")
	if b, err := builder.Build(); err != nil || b == nil {
		t.Fatal()
	} else {
		for i := 0; i < 10000; i++ {
			if b.NextDelayMillis(i) != 456 {
				t.Fatal()
			}
		}

		b, _ = builder.Build() // build again
		for i := 0; i < 10000; i++ {
			if b.NextDelayMillis(i) != 456 {
				t.Fatal()
			}
		}
	}

	// error spec
	builder = NewBackoffBuilder().BaseBackoffSpec("fixe=456")
	if b, err := builder.Build(); err == nil || b != nil {
		t.Fatal()
	}

	fixedBackoff, _ := NewFixedBackoff(123)
	builder = NewBackoffBuilder().
		BaseBackoff(fixedBackoff).
		WithLimit(5).
		WithJitter(0.9).
		WithJitterBound(0.9, 1.2)

	if _, err := builder.Build(); err == nil {
		t.Fatal()
	}
}

func TestBuilderNodelayBackoff(t *testing.T) {
	builder := NewBackoffBuilder().BaseBackoff(NoDelayBackoff)
	if b, err := builder.Build(); err != nil || b == nil {
		t.Fatal()
	}

	builder = NewBackoffBuilder().
		BaseBackoff(NoDelayBackoff).
		WithJitter(0.9)
	if _, err := builder.Build(); err != nil {
		t.Fatal()
	}

	builder = NewBackoffBuilder().
		BaseBackoff(NoDelayBackoff).
		WithJitterBound(0.9, 1.2)
	if _, err := builder.Build(); err == nil {
		t.Fatal()
	}

	builder = NewBackoffBuilder().
		BaseBackoff(NoDelayBackoff).
		WithJitter(0.9).
		WithJitterBound(0.9, 1.2)
	if _, err := builder.Build(); err == nil {
		t.Fatal()
	}
}

func TestBuilderFixedBackoffWithLimit(t *testing.T) {
	fixedBackoff, _ := NewFixedBackoff(123)

	builder := NewBackoffBuilder().
		BaseBackoff(fixedBackoff).
		WithLimit(5)

	if b, err := builder.Build(); err != nil {
		t.Fatal()
	} else {
		for i := 0; i < 100; i++ {
			d := b.NextDelayMillis(i)
			if i < 5 && d != 123 {
				t.Fatal()
			}

			if i >= 5 && d >= 0 {
				t.Fatal()
			}
		}
	}

	builder = NewBackoffBuilder().
		BaseBackoff(fixedBackoff).
		WithLimit(-1).
		WithJitter(0.9).
		WithJitterBound(0.9, 1.2)
	if _, err := builder.Build(); err == nil {
		t.Fatal()
	}
}

func TestParseInvalidSpec(t *testing.T) {
	// test exponential
	if _, err := parseFromSpec("exponential="); err != ErrInvalidSpecFormat {
		t.Fatal()
	}

	if _, err := parseFromSpec("exponential=1:"); err != ErrInvalidSpecFormat {
		t.Fatal()
	}

	if _, err := parseFromSpec("exponential=1:2"); err != ErrInvalidSpecFormat {
		t.Fatal()
	}

	if _, err := parseFromSpec("exponential=a:2:3"); err == nil {
		t.Fatal()
	}

	if _, err := parseFromSpec("exponential=1:a:3"); err == nil {
		t.Fatal()
	}

	if _, err := parseFromSpec("exponential=1:2:a"); err == nil {
		t.Fatal()
	}
}

func TestParseSpec(t *testing.T) {
	cases := []string{
		"exponential=1:2:3",
		"exponential=:201:3",
		"exponential=::3",
		"exponential=::",
	}

	expectedinitialDelayMillis := []int64{
		1,
		DefaultDelayMillis,
		DefaultDelayMillis,
		DefaultDelayMillis,
	}

	expectedMaxDelayMillis := []int64{
		2,
		201,
		DefaultMaxDelayMillis,
		DefaultMaxDelayMillis,
	}

	expectedMultipiler := []float64{
		3,
		3,
		3,
		DefaultMultiplier,
	}

	for i := range cases {
		if b, err := parseFromSpec(cases[i]); err != nil {
			t.Fatal(err)
		} else {
			tmp := b.(*ExponentialBackoff)
			if tmp.initialDelayMillis != expectedinitialDelayMillis[i] ||
				tmp.maxDelayMillis != expectedMaxDelayMillis[i] ||
				tmp.multiplier != expectedMultipiler[i] {
				t.Fatal()
			}
		}
	}
}
