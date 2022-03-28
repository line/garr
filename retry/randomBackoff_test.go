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

import "testing"

func TestRandomBackoff(t *testing.T) {
	if _, err := NewRandomBackoff(-1, 2); err == nil {
		t.FailNow()
	}

	if _, err := NewRandomBackoff(3, 2); err == nil {
		t.FailNow()
	}

	if r, err := NewRandomBackoff(1000, 1000); err != nil {
		t.Error(err)
		t.FailNow()
	} else if d := r.NextDelayMillis(1); d != 1000 {
		t.FailNow()
	}

	if r, err := NewRandomBackoff(1000, 1200); err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		for i := 0; i < 1000; i++ {
			if d := r.NextDelayMillis(i); d < 1000 || d > 1200 {
				t.FailNow()
			}
		}
	}
}
