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

package cbreaker

import "testing"

func TestEventCount(t *testing.T) {
	ev := NewEventCount(0, 0)
	if ev.SuccessRate() != -1 {
		t.Errorf("Fail to catch trivial case of success rate")
	}
	if ev.FailureRate() != -1 {
		t.Errorf("Fail to catch trivial case of failure rate")
	}

	ev = NewEventCount(5, 20)
	if ev.Success() != 5 || ev.success != 5 || ev.Failure() != 20 || ev.failure != 20 || ev.Total() != 25 {
		t.Errorf("Fail to create new EventCount")
	}

	if ev.SuccessRate() != 0.2 || ev.FailureRate() != 0.8 {
		t.Errorf("Fail to return rate")
	}
}
