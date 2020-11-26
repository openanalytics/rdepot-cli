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
	"bytes"
	// "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"openanalytics.eu/rdepot/cli/model"
	"os"
	"strconv"
	"strings"
)

type RDepotConfig struct {
	Host  string
	Token string
}

func DefaultClient() *http.Client {
	return http.DefaultClient
}

func ListPackages(client *http.Client, cfg RDepotConfig) ([]model.Package, error) {

	req, err := http.NewRequest(
		"GET",
		cfg.Host+"/api/manager/packages/list",
		nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var pkgs []model.Package
	if err := json.Unmarshal(body, &pkgs); err != nil {
		return nil, fmt.Errorf("could not unpack response: %s", err)
	}
	return pkgs, nil

}

type Pair struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

type SubmissionResult struct {
	Success Pair `json:"success"`
	Warning Pair `json:"warning"`
	Error   Pair `json:"error"`
}

type Message interface {
	Class() (string, error)
	Content() string
}

func (s SubmissionResult) Class() (string, error) {
	switch {
	case s.Success.First != "":
		return "success", nil
	case s.Warning.First != "":
		return "warning", nil
	case s.Error.First != "":
		return "error", nil
	default:
		return "", fmt.Errorf("Unrecognized response type")
	}
}

func (s SubmissionResult) Content() string {
	return s.Success.Second + s.Warning.Second + s.Error.Second
}

func SubmitPackage(client *http.Client, cfg RDepotConfig, archive string, repository string, replace bool) (Message, error) {
	var subres SubmissionResult
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fr, err := os.Open(archive)
	if err != nil {
		return subres, err
	}

	if fw, err := createFormGZip(w, "file", archive); err != nil {
		return subres, err
	} else {
		io.Copy(fw, fr)
	}

	if err := w.WriteField("repository", repository); err != nil {
		return subres, err
	}

	if err := w.WriteField("replace", strconv.FormatBool(replace)); err != nil {
		return subres, err
	}

	w.Close()

	req, err := http.NewRequest(
		"POST",
		cfg.Host+"/api/manager/packages/submit",
		&b)
	if err != nil {
		return subres, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	res, err := client.Do(req)
	if err != nil {
		return subres, err
	}

	if res.StatusCode != http.StatusOK {
		return subres, fmt.Errorf("bad status: %s", res.Status)
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return subres, err
	}

	if err := json.Unmarshal(resBody, &subres); err != nil {
		return subres, fmt.Errorf("could not unpack response: %s", err)
	}

	return subres, nil

}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createFormGZip(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", "application/gzip")
	return w.CreatePart(h)
}
