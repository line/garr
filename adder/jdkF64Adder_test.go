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
	"testing"
)

func TestJDKF64AdderNotRaceInc(t *testing.T) {
	testF64AdderNotRaceInc(t, JDKF64AdderType)
}

func TestJDKF64AdderRaceInc(t *testing.T) {
	testF64AdderRaceInc(t, JDKF64AdderType)
}

func TestJDKF64AdderNotRaceDec(t *testing.T) {
	testF64AdderNotRaceDec(t, JDKF64AdderType)
}

func TestJDKF64AdderRaceDec(t *testing.T) {
	testF64AdderRaceDec(t, JDKF64AdderType)
}

func TestJDKF64AdderNotRaceAdd(t *testing.T) {
	testF64AdderNotRaceAdd(t, JDKF64AdderType)
}

func TestJDKF64AdderRaceAdd(t *testing.T) {
	testF64AdderRaceAdd(t, JDKF64AdderType)
}
