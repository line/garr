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

package httpclient

import (
	"errors"
	"fmt"
	"net/url"
	"testing"
)

func TestErrors(t *testing.T) {
	u, _ := url.Parse("https://google.com:8091")

	e := fmt.Errorf("Decoding error")
	v := decodingError(u, e)
	if !errors.Is(v.(*errorWrap).err, e) {
		t.FailNow()
	}
	if v.(*errorWrap).category != errDecoding {
		t.FailNow()
	}

	e = fmt.Errorf("Transform error")
	v = transformError(u, e)
	if !errors.Is(v.(*errorWrap).err, e) {
		t.FailNow()
	}
	if v.(*errorWrap).category != errTransform {
		t.FailNow()
	}

	e = fmt.Errorf("Fake Connection error")
	v = connectionError(u, e)
	if !errors.Is(v.(*errorWrap).err, e) {
		t.FailNow()
	}
	if v.(*errorWrap).category != errConnection {
		t.FailNow()
	}
	if v.Error() != "Connection to host:[google.com:8091] got error:[Fake Connection error]" {
		t.FailNow()
	}
}
