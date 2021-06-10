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
	"fmt"

	"path/filepath"
)

type Package struct {
	Id                 int        `json:"id"`
	Version            Version    `json:"version"`
	Submission         Submission `json:"submission"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Author             string     `json:"author"`
	Depends            string     `json:"depends"`
	Imports            string     `json:"imports"`
	Suggests           string     `json:"suggests"`
	SystemRequirements string     `json:"systemRequirements"`
	License            string     `json:"license"`
	Title              string     `json:"title"`
	Url                string     `json:"url"`
	Source             string     `json:"source"`
	Md5Sum             string     `json:"md5sum"`
	Active             bool       `json:"active"`
	Deleted            bool       `json:"deleted"`
}

type Submission struct {
	Id       int    `json:"id"`
	Changes  string `json:"changes"`
	Accepted bool   `json:"accepted"`
	Deleted  bool   `json:"deleted"`
}

func (pkg Package) Summary() string {
	return fmt.Sprintf("%s %s", pkg.Name, pkg.Version.CanonicalRep)
}

// Filter packages matching a name glob pattern
func FilterByName(packages []Package, name string) ([]Package, error) {
	filtered := make([]Package, 0)

	for _, pkg := range packages {
		matched, err := filepath.Match(name, pkg.Name)
		if err != nil {
			return nil, err
		} else if matched {
			filtered = append(filtered, pkg)
		}
	}

	return filtered, nil
}

// Retain only archived packages
func FilterArchived(packages []Package) []Package {

	newest := make(map[string]Version)

	for _, pkg := range packages {
		if newest[pkg.Name].Less(pkg.Version) {
			newest[pkg.Name] = pkg.Version
		}
	}

	filtered := make([]Package, 0, len(newest))

	for _, pkg := range packages {
		if pkg.Version.Less(newest[pkg.Name]) {
			filtered = append(filtered, pkg)
		}
	}

	return filtered
}
