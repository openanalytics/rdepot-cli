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
	"strconv"
	"strings"
)

type Version struct {
	Digits       []int
	CanonicalRep string
}

func CanonicalVersion(rep string) (*Version, error) {
	dotted := strings.ReplaceAll(rep, "-", ".")
	parts := strings.Split(dotted, ".")
	var digits = []int{}
	for _, part := range parts {
		digit, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		digits = append(digits, digit)
	}
	return &Version{
		Digits:       digits,
		CanonicalRep: rep,
	}, nil
}

func (u *Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.CanonicalRep)
}

func (u *Version) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	v, err := CanonicalVersion(str)
	if err != nil {
		return err
	}
	u.CanonicalRep = v.CanonicalRep
	u.Digits = v.Digits
	return nil
}

func (x Version) Equals(y Version) bool {
	return !x.Less(y) && !y.Less(x)
}

func (x Version) Less(y Version) bool {
	for i, d := range x.Digits {
		if i+1 > len(y.Digits) || d > y.Digits[i] {
			return false
		}
		if d < y.Digits[i] {
			return true
		}
	}
	return len(x.Digits) < len(y.Digits)
}
