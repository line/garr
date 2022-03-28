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

import (
	"fmt"
	"testing"
)

var loggedInfo string

var loggedWarn string

var loggedError string

type fakeLogger struct{}

func (f *fakeLogger) Info(i string) {
	loggedInfo = i
}

func (f *fakeLogger) Warn(title string, v interface{}) {
	loggedWarn = title + fmt.Sprintf("_%v", v)
}

func (f *fakeLogger) Error(title string, v interface{}) {
	loggedError = title + fmt.Sprintf("_%v", v)
}

func TestSetLogger(t *testing.T) {
	SetDefaultLogger(&fakeLogger{})

	if logger.Info("info"); loggedInfo != "info" {
		t.FailNow()
	}

	if logger.Warn("warn", "test"); loggedWarn != "warn_test" {
		t.FailNow()
	}

	if logger.Error("error", "test"); loggedError != "error_test" {
		t.FailNow()
	}

	if SetDefaultLogger(nil); logger != nil {
		t.FailNow()
	}
}
