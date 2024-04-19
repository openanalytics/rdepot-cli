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
	"os"
	"strconv"
	"strings"

	"openanalytics.eu/rdepot/cli/model"
)

type RDepotConfig struct {
	Host  string
	Token string
}

func DefaultClient() *http.Client {
	return http.DefaultClient
}

func SoftDeletePackage(client *http.Client, cfg RDepotConfig, id int) error {

	patchJson := []byte(`[{"op": "replace", "path": "/deleted", "value": true}]`)

	req, err := http.NewRequest("PATCH", cfg.Host+fmt.Sprintf("/api/v2/manager/r/packages/%d", id), bytes.NewBuffer(patchJson))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json-patch+json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	return nil
}

func DeletePackage(client *http.Client, cfg RDepotConfig, id int) error {
	err := SoftDeletePackage(client, cfg, id)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"DELETE",
		cfg.Host+fmt.Sprintf("/api/v2/manager/r/packages/%d", id),
		nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	return nil
}

func ListPackagesPage(client *http.Client, cfg RDepotConfig, repository string, page int) ([]model.Package, model.Page, error) {
	var pageData model.Page
	req, err := http.NewRequest(
		"GET",
		cfg.Host+"/api/v2/manager/r/packages",
		nil)
	if err != nil {
		return nil, pageData, err
	}

	q := req.URL.Query()
	if repository != "" {
		q.Add("repositoryName", repository)
	}
	q.Add("page", strconv.Itoa(page))
	q.Add("size", "100")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	res, err := client.Do(req)

	if err != nil {
		return nil, pageData, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, pageData, fmt.Errorf("bad status: %s", res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, pageData, err
	}

	var response model.Response[model.Package]
	fmt.Print(res)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, pageData, fmt.Errorf("could not unpack response: %s", err)
	}
	return response.Data.Content, response.Data.Page, nil

}

func ListPackages(client *http.Client, cfg RDepotConfig, repository string) ([]model.Package, error) {
	pkgs, pageData, err := ListPackagesPage(client, cfg, repository, 0)
	if err != nil {
		return nil, err
	}
	for page := 1; page <= pageData.TotalPages; page++ {
		new_pkgs, _, err := ListPackagesPage(client, cfg, repository, page)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, new_pkgs...)
	}

	return pkgs, nil
}

type SubmissionResult struct {
	Status      string `json:"status"`
	Code        int    `json:"code"`
	Message     string `json:"message"`
	MessageCode string `json:"messageCode"`
}

func SubmitPackage(client *http.Client, cfg RDepotConfig, archive string, repository string, replace bool, generateManual bool) (string, error) {
	var subres SubmissionResult
	var msg string
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fr, err := os.Open(archive)
	if err != nil {
		return msg, err
	}

	if fw, err := createFormGZip(w, "file", archive); err != nil {
		return msg, err
	} else {
		io.Copy(fw, fr)
	}

	if err := w.WriteField("repository", repository); err != nil {
		return msg, err
	}

	if err := w.WriteField("replace", strconv.FormatBool(replace)); err != nil {
		return msg, err
	}

	if !generateManual { // TODO: remove in future versions
		if err := w.WriteField("generateManual", strconv.FormatBool(generateManual)); err != nil {
			return msg, err
		}
	}

	w.Close()

	req, err := http.NewRequest(
		"POST",
		cfg.Host+"/api/v2/manager/r/submissions",
		&b)
	if err != nil {
		return msg, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	res, err := client.Do(req)
	if err != nil {
		return msg, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return msg, fmt.Errorf("bad status: %s", res.Status)
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return msg, err
	}
	if err := json.Unmarshal(resBody, &subres); err != nil {
		return msg, fmt.Errorf("could not unpack response: %s", err)
	}
	msg = subres.Message
	return msg, nil

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
