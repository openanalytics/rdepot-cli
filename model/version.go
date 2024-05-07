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
	"fmt"
	"strconv"
	"strings"
)

type VersionSegment struct {
	Digit              int
	ReleaseType        string
	ReleaseTypeVersion int
}

type Version struct {
	Epoch        int
	Segments     []VersionSegment
	CanonicalRep string
}

func CanonicalVersion(rep string) (*Version, error) {
	epoch := 0
	if strings.Contains(rep, "!") { // version contains epoch
		epochVersionSplit := strings.Split(rep, "!")
		var err error
		epoch, err = strconv.Atoi(epochVersionSplit[0])
		if err != nil {
			return nil, fmt.Errorf("invalid epoch in version %s", rep)
		}
		rep = epochVersionSplit[1]
	}
	dotted := strings.ReplaceAll(rep, "-", ".")
	parts := strings.Split(dotted, ".")
	var segments = []VersionSegment{}
	for i, part := range parts {
		digit, err := strconv.Atoi(part)
		if err == nil {
			segments = append(segments, VersionSegment{digit, "", 0})
		} else if i > 0 { // try parsing segment as a python segment
			parsedSegment, err := ParsePythonVersionSegment(part)
			if err != nil {
				return nil, err
			}
			segments = append(segments, *parsedSegment)
		} else {
			return nil, err
		}
	}
	return &Version{
		Epoch:        epoch,
		Segments:     segments,
		CanonicalRep: rep,
	}, nil
}

func ParsePythonVersionSegment(part string) (*VersionSegment, error) {
	preReleaseCycles := []string{"a", "b", "rc", "c"}
	for _, cycle := range preReleaseCycles {
		if strings.Contains(part, cycle) {
			return PreReleaseToVersionSegment(part)
		}
	}
	if strings.Contains(part, "post") || strings.Contains(part, "dev") {
		return PostOrDevReleaseToVersionSegment(part)
	}
	return nil, fmt.Errorf("invalid python version segment %s", part)
}

func PreReleaseToVersionSegment(preReleasePart string) (*VersionSegment, error) {
	var releaseType string
	if strings.Contains(preReleasePart, "a") {
		releaseType = "a"
	} else if strings.Contains(preReleasePart, "b") {
		releaseType = "b"
	} else if strings.Contains(preReleasePart, "rc") {
		releaseType = "rc"
	} else if strings.Contains(preReleasePart, "c") {
		preReleasePart = strings.Replace(preReleasePart, "c", "rc", 1)
		releaseType = "rc"
	} else {
		return nil, fmt.Errorf("invalid pre-release version %s", preReleasePart)
	}
	preReleaseParts := strings.Split(preReleasePart, releaseType)

	digit, err := strconv.Atoi(preReleaseParts[0])
	if err != nil {
		return nil, fmt.Errorf(preReleasePart)
	}
	releaseTypeVersion := 0
	if len(preReleaseParts) == 3 {
		releaseTypeVersion, err = strconv.Atoi(preReleaseParts[2])
		if err != nil {
			return nil, err
		}
	}
	return &VersionSegment{digit, releaseType, releaseTypeVersion}, nil
}

func PostOrDevReleaseToVersionSegment(part string) (*VersionSegment, error) {
	var releaseType string
	if strings.Contains(part, "post") {
		releaseType = "post"
	} else if strings.Contains(part, "dev") {
		releaseType = "dev"
	}
	preReleaseParts := strings.Split(part, releaseType)
	digit, err := strconv.Atoi(preReleaseParts[1])
	if err != nil {
		return nil, err
	}
	return &VersionSegment{digit, releaseType, 0}, nil
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
	u.Segments = v.Segments
	u.Epoch = v.Epoch
	return nil
}

func (x VersionSegment) Equals(y VersionSegment) bool {
	return !x.Less(y) && !y.Less(x)
}

func (x VersionSegment) Less(y VersionSegment) bool {
	if x.Digit == y.Digit && x.ReleaseType == y.ReleaseType {
		return x.ReleaseTypeVersion < y.ReleaseTypeVersion
	} else if x.Digit == y.Digit {
		return x.ReleaseType < y.ReleaseType
	} else {
		return x.Digit < y.Digit
	}
}

func (x Version) Equals(y Version) bool {
	return !x.Less(y) && !y.Less(x)
}

func (x Version) Less(y Version) bool {
	if x.Epoch < y.Epoch {
		return true
	}
	for i, d := range x.Segments {
		if i+1 > len(y.Segments) || y.Segments[i].Less(d) {
			return false
		}
		if d.Less(y.Segments[i]) {
			return true
		}
	}
	return len(x.Segments) < len(y.Segments)
}

type PythonVersion struct {
	CanonicalRep string
}
