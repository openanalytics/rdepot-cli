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

type Response[C any] struct {
	Status      string  `json:"status"`
	Code        int     `json:"code"`
	Message     string  `json:"message"`
	MessageCode string  `json:"messageCode"`
	Data        Data[C] `json:"data"`
}

type Data[C any] struct {
	Links   []Link `json:"links"`
	Content []C    `json:"content"`
	Page    Page   `json:"page"`
}

type Page struct {
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	Number        int `json:"number"`
}

type Package struct {
	Id                 int        `json:"id"`
	Username           User       `json:"user"`
	Repository         Repository `json:"repository"`
	Version            Version    `json:"version"`
	Name               string     `json:"name"`
	SubmissionId       int        `json:"submissionId"`
	Description        string     `json:"description"`
	Author             string     `json:"author"`
	Title              string     `json:"title"`
	Url                string     `json:"url"`
	Source             string     `json:"source"`
	Active             bool       `json:"active"`
	Deleted            bool       `json:"deleted"`
	Depends            string     `json:"depends"`
	Imports            string     `json:"imports"`
	Suggests           string     `json:"suggests"`
	SystemRequirements string     `json:"systemRequirements"`
	License            string     `json:"license"`
	Md5Sum             string     `json:"md5sum"`
	Links              []Link     `json:"links"`
}

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type Repository struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	PublicationUri string `json:"publicationUri"`
}

type Link struct {
	Rel         string `json:"rel"`
	Href        string `json:"href"`
	Hreflang    string `json:"hreflang"`
	Media       string `json:"media"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Deprecation string `json:"deprecation"`
	Profile     string `json:"profile"`
	Name        string `json:"name"`
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
