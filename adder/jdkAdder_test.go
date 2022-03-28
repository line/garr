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

func TestJDKAdderNotRaceInc(t *testing.T) {
	testAdderNotRaceInc(t, JDKAdderType)
}

func TestJDKAdderRaceInc(t *testing.T) {
	testAdderRaceInc(t, JDKAdderType)
}

func TestJDKAdderNotRaceDec(t *testing.T) {
	testAdderNotRaceDec(t, JDKAdderType)
}

func TestJDKAdderRaceDec(t *testing.T) {
	testAdderRaceDec(t, JDKAdderType)
}

func TestJDKAdderNotRaceAdd(t *testing.T) {
	testAdderNotRaceAdd(t, JDKAdderType)
}

func TestJDKAdderRaceAdd(t *testing.T) {
	testAdderRaceAdd(t, JDKAdderType)
}
