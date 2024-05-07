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

package model

import (
	"encoding/json"
	"testing"
)

func TestCanonicalVersion(t *testing.T) {
	if _, err := CanonicalVersion("alpha.1"); err == nil {
		t.Errorf("expected error for alpha.1")
	}
}

func TestMarshal(t *testing.T) {
	v := pkgv("1.2.3")
	res, _ := json.MarshalIndent(v, "", "  ")
	if string(res) != "\"1.2.3\"" {
		t.Errorf("expected 1.2.3, got %s", string(res))
	}
}

func pkgv(str string) *Version {
	if v, err := CanonicalVersion(str); err != nil {
		panic(err)
	} else {
		return v
	}
}

func TestEquals(t *testing.T) {

	var tests = []struct {
		x     *Version
		y     *Version
		equal bool
	}{
		{x: pkgv("1.2.3"), y: pkgv("1.2.3"), equal: true},
		{x: pkgv("1.1"), y: pkgv("1-1"), equal: true},
		{x: pkgv("1.0"), y: pkgv("2.0"), equal: false},
	}

	for _, test := range tests {

		equal := test.x.Equals(*test.y)
		if equal != test.equal {
			t.Errorf("%s == %s : expected %t, got %t",
				test.x.CanonicalRep,
				test.y.CanonicalRep,
				test.equal,
				equal)
		}
	}

}

func TestLess(t *testing.T) {

	var tests = []struct {
		x    *Version
		y    *Version
		less bool
	}{
		{x: pkgv("1.0"), y: pkgv("0.9"), less: false},
		{x: pkgv("1.0"), y: pkgv("1.0"), less: false},
		{x: pkgv("1.0.0"), y: pkgv("2.0.0"), less: true},
		{x: pkgv("2.0.0"), y: pkgv("1.0.0"), less: false},
		{x: pkgv("0.99.0"), y: pkgv("1.0.0"), less: true},
		{x: pkgv("1.1"), y: pkgv("1-1"), less: false},
		{x: pkgv("1-1"), y: pkgv("1.1"), less: false},
		{x: pkgv("1"), y: pkgv("1.0"), less: true},
		{x: pkgv("1.0"), y: pkgv("1"), less: false},
		{x: pkgv("1.0a1"), y: pkgv("1.0b10"), less: true},
		{x: pkgv("1.0rc1"), y: pkgv("1.0b10"), less: false},
		{x: pkgv("1.0c1"), y: pkgv("1.0a10"), less: false},
	}

	for _, test := range tests {

		less := test.x.Less(*test.y)
		if less != test.less {
			t.Errorf("%s < %s : expected %t, got %t",
				test.x.CanonicalRep,
				test.y.CanonicalRep,
				test.less,
				less)
		}
	}

}
