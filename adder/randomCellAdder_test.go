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

func TestRandomCellAdderNotRaceInc(t *testing.T) {
	testAdderNotRaceInc(t, RandomCellAdderType)
}

func TestRandomCellAdderRaceInc(t *testing.T) {
	testAdderRaceInc(t, RandomCellAdderType)
}

func TestRandomCellAdderNotRaceDec(t *testing.T) {
	testAdderNotRaceDec(t, RandomCellAdderType)
}

func TestRandomCellAdderRaceDec(t *testing.T) {
	testAdderRaceDec(t, RandomCellAdderType)
}

func TestRandomCellAdderNotRaceAdd(t *testing.T) {
	testAdderNotRaceAdd(t, RandomCellAdderType)
}

func TestRandomCellAdderRaceAdd(t *testing.T) {
	testAdderRaceAdd(t, RandomCellAdderType)
}
