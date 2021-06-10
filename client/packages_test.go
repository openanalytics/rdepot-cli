// Copyright 2020 Open Analytics
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestListPackages(t *testing.T) {

	var tests = []struct {
		body  []byte
		nPkgs int
	}{
		{
			body:  []byte(`[]`),
			nPkgs: 0,
		},
		{
			body:  []byte(`[{"id":4,"version":"0.0.1","submission":{"id":4,"changes":null,"accepted":true,"deleted":false},"name":"foo","description":"foo description","author":"Foo Author","depends":null,"imports":null,"suggests":null,"systemRequirements":null,"license":"GPL-2","title":"foo title","source":"/opt/rdepot/repositories/3/36142023/foo_0.0.1.tar.gz","md5sum":"de696d506f435e040f0b215da4d3c643","active":true,"deleted":false,"packageEvents":null}]`),
			nPkgs: 1,
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			want := "/api/manager/packages/list"
			if req.URL.String() != want {
				t.Errorf("Expected %s, got %s", want, req.URL.String())
			}
			rw.Write(test.body)
		}))
		defer server.Close()

		config := RDepotConfig{Host: server.URL, Token: "validtoken"}

		res, err := ListPackages(server.Client(), config, "")

		if err != nil {
			t.Errorf("Got error: %s", err)
		}

		if len(res) != test.nPkgs {
			t.Errorf("Expected %d packages, got %d", test.nPkgs, len(res))
		}

	}

}

func expectEqual(t *testing.T, expected interface{}, actual interface{}) {
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestSubmitPackage(t *testing.T) {
	var tests = []struct {
		body    []byte
		replace bool
	}{
		{
			body:    []byte(`{"success": {"first": "oaColors_0.0.4.tar.gz", "second": "submission created successfully"}}`),
			replace: false,
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if actual, expect := req.URL.String(), "/api/manager/packages/submit"; actual != expect {
				t.Errorf("Expected %s, got %s", expect, actual)
			}
			expectEqual(t, "/api/manager/packages/submit", req.URL.String())
			expectEqual(t, strconv.FormatBool(test.replace), req.FormValue("replace"))
			expectEqual(t, "test", req.FormValue("repository"))
			if _, fh, err := req.FormFile("file"); err != nil {
				t.Errorf("Error: %s", err)
			} else if actual, expect := fh.Header.Get("Content-Type"), "application/gzip"; actual != expect {
				t.Errorf("Expected content type: %s, got %s", expect, actual)
			}

			rw.Write(test.body)
		}))
		defer server.Close()

		config := RDepotConfig{Host: server.URL, Token: "validtoken"}

		res, err := SubmitPackage(server.Client(), config, "testdata/oaColors_0.0.4.tar.gz", "test", test.replace, true)

		if err != nil {
			t.Errorf("Error: %s", err)
		}

		if mc, err := res.Class(); err != nil {
			t.Errorf("Error: %s", err)
			if mc != "success" {
				t.Errorf("Unexpected message class: %s", mc)
			}
		}

	}
}
