// Copyright 2020-2024 Open Analytics
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

	"path/filepath"
)

type Response[C any] struct {
	Status      string  `json:"status"`
	Code        int     `json:"code"`
	Message     string  `json:"message"`
	MessageCode string  `json:"messageCode"`
	Data        Data[C] `json:"data"`
}

func (r *Response[C]) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, r); err != nil {
		return fmt.Errorf("could not unpack response: %s", err)
	}
	return nil
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
	Id          int        `json:"id"`
	User        User       `json:"user"`
	Repository  Repository `json:"repository"`
	Submission  Submission `json:"submission"`
	Name        string     `json:"name"`
	Version     Version    `json:"version"`
	Source      string     `json:"source"`
	Active      bool       `json:"active"`
	Deleted     bool       `json:"deleted"`
	Technology  string     `json:"technology"`
	Description string     `json:"description"`
	Author      string     `json:"author"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Links       []Link     `json:"links"`
}

type RPackage struct {
	Package
	Depends            string `json:"depends"`
	Imports            string `json:"imports"`
	Suggests           string `json:"suggests"`
	SystemRequirements string `json:"systemRequirements"`
	License            string `json:"license"`
	Md5sum             string `json:"md5sum"`
}

type PythonPackage struct {
	Package
	AuthorEmail            string `json:"authorEmail"`
	Classifiers            string `json:"classifiers"`
	DescriptionContentType string `json:"descriptionContentType"`
	HomePage               string `json:"homePage"`
	Keywords               string `json:"keywords"`
	License                string `json:"license"`
	Maintainer             string `json:"maintainer"`
	MaintainerEmail        string `json:"maintainerEmail"`
	Platform               string `json:"platform"`
	ProjectUrl             string `json:"projectUrl"`
	ProvidesExtra          string `json:"providesExtra"`
	RequiresDist           string `json:"requiresDist"`
	RequiresExternal       string `json:"requiresExternal"`
	RequiresPython         string `json:"requiresPython"`
	SummaryField           string `json:"summary"`
	Hash                   string `json:"hash"`
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
	Published      bool   `json:"published"`
	Technology     string `json:"technology"`
}

type Submission struct {
	Id    int    `json:"id"`
	State string `json:"state"`
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

type GenericPackage interface {
	GetName() string
	GetVersion() Version
	Summary() string
	GetId() int
}

func (p Package) GetName() string {
	return p.Name
}

func (p Package) GetVersion() Version {
	return p.Version
}

func (p Package) GetId() int {
	return p.Id
}

func (pkg Package) Summary() string {
	return fmt.Sprintf("%s %s", pkg.Name, pkg.Version.CanonicalRep)
}

// Filter packages matching a name glob pattern
func FilterByName[G GenericPackage](packages []G, name string) ([]G, error) {
	filtered := make([]G, 0)

	for _, pkg := range packages {
		matched, err := filepath.Match(name, pkg.GetName())
		if err != nil {
			return nil, err
		} else if matched {
			filtered = append(filtered, pkg)
		}
	}

	return filtered, nil
}

// Retain only archived packages
func FilterArchived[G GenericPackage](packages []G) []G {

	newest := make(map[string]Version)

	for _, pkg := range packages {
		if newest[pkg.GetName()].Less(pkg.GetVersion()) {
			newest[pkg.GetName()] = pkg.GetVersion()
		}
	}

	filtered := make([]G, 0, len(newest))

	for _, pkg := range packages {
		if pkg.GetVersion().Less(newest[pkg.GetName()]) {
			filtered = append(filtered, pkg)
		}
	}

	return filtered
}
