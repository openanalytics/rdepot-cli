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
	"strings"
	"testing"
)

func TestListPackages(t *testing.T) {

	var tests = []struct {
		body  []byte
		nPkgs int
	}{
		{
			body:  []byte(`{ "status": "SUCCESS", "code": 200, "message": "Your request has been processed successfully.", "messageCode": "success.request.processed", "data": { "links": [ { "rel": "self", "href": "http://localhost:8017/api/v2/manager/r/packages?page=0&size=1" } ], "content": [], "page": { "size": 0, "totalElements": 0, "totalPages": 0, "number": 0 } }}`),
			nPkgs: 0,
		},
		{
			body:  []byte(`{ "status": "SUCCESS", "code": 200, "message": "Your request has been processed successfully.", "messageCode": "success.request.processed", "data": { "links": [ { "rel": "first", "href": "http://localhost:8017/api/v2/manager/r/packages?page=0&size=1" }, { "rel": "self", "href": "http://localhost:8017/api/v2/manager/r/packages?page=0&size=1" }, { "rel": "next", "href": "http://localhost:8017/api/v2/manager/r/packages?page=1&size=1" }, { "rel": "last", "href": "http://localhost:8017/api/v2/manager/r/packages?page=19&size=1" } ], "content": [ { "id": 8, "user": { "id": 4, "name": "Albert Einstein", "login": "einstein", "email": "einstein@ldap.forumsys.com" }, "repository": { "id": 3, "name": "testrepo2", "publicationUri": "http://localhost/repo/testrepo2" }, "submissionId": 6, "name": "accrued", "version": "1.2", "description": "Package for visualizing data quality of partially accruing time series.", "author": "Julie Eaton and Ian Painter", "title": "Visualization tools for partially accruing data", "url": null, "source": "/opt/rdepot/repositories/3/83118397/accrued_1.2.tar.gz", "active": true, "deleted": false, "depends": "R (>= 3.0), grid", "imports": null, "suggests": null, "systemRequirements": null, "license": "GPL-3", "md5sum": "70d295115295a4718593f6a39d77add9", "links": [ { "rel": "self", "href": "http://localhost:8017/api/v2/manager/r/packages/8" }, { "rel": "packageList", "href": "http://localhost:8017/api/v2/manager/r/packages" } ] } ], "page": { "size": 1, "totalElements": 1, "totalPages": 0, "number": 0 } }}`),
			nPkgs: 1,
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			want := "/api/v2/manager/packages"
			// Remove the query parameters from the url
			expect := strings.Split(req.URL.String(), "?")[0]
			if expect != want {
				t.Errorf("Expected %s, got %s", want, expect)
			}
			rw.Write(test.body)
		}))
		defer server.Close()

		config := RDepotConfig{Host: server.URL, Token: "validtoken", Technology: "all"}

		res, err := ListPackages(server.Client(), config, "", false, "")

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
			body:    []byte(`{"status": "SUCCESS", "code": 201, "message": "Your resource has been created successfully.", "messageCode": "success.resource.created", "data": {}}`),
			replace: false,
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if actual, expect := req.URL.String(), "/api/v2/manager/r/submissions"; actual != expect {
				t.Errorf("Expected %s, got %s", expect, actual)
			}
			expectEqual(t, "/api/v2/manager/r/submissions", req.URL.String())
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

		config := RDepotConfig{Host: server.URL, Token: "validtoken", Technology: "r"}

		_, err := SubmitPackage(server.Client(), config, "testdata/oaColors_0.0.4.tar.gz", "test", test.replace, true)

		if err != nil {
			t.Errorf("Error: %s", err)
		}
	}
}
