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

package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"strings"

	"openanalytics.eu/rdepot/cli/model"
)

type RDepotConfig struct {
	Host       string
	Token      string
	Username   string
	Technology string
}

func DefaultClient() *http.Client {
	return http.DefaultClient
}

func basicAuth(username string, token string) string {
	var auth string
	if username == "" {
		auth = token
	} else {
		auth = username + ":" + token
	}
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func SoftDeletePackage(client *http.Client, cfg RDepotConfig, id int) error {

	patchJson := []byte(`[{"op": "replace", "path": "/deleted", "value": true}]`)
	if cfg.Technology != "python" && cfg.Technology != "r" {
		return fmt.Errorf("invalid technology provided for deleting only Python and R are supported")
	}
	path, err := technologyToPath(cfg.Technology)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", cfg.Host+fmt.Sprintf("/api/v2/manager/"+path+"packages/%d", id), bytes.NewBuffer(patchJson))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json-patch+json")
	req.Header.Set("Authorization", "Basic "+basicAuth(cfg.Username, cfg.Token))
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

func DeletePackage(client *http.Client, cfg RDepotConfig, pkg model.Package) error {
	if !pkg.Deleted {
		err := SoftDeletePackage(client, cfg, pkg.Id)
		if err != nil {
			return err
		}
	}

	if cfg.Technology != "python" && cfg.Technology != "r" {
		return fmt.Errorf("invalid technology provided for deleting only Python and R are supported")
	}
	path, err := technologyToPath(strings.ToLower(pkg.Technology))
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"DELETE",
		cfg.Host+fmt.Sprintf("/api/v2/manager/"+path+"packages/%d", pkg.Id),
		nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(cfg.Username, cfg.Token))

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

func ListPackagesPage(client *http.Client, cfg RDepotConfig, repository string, page int) ([]byte, error) {
	path, err := technologyToPath(cfg.Technology)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"GET",
		cfg.Host+"/api/v2/manager/"+path+"packages",
		nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if repository != "" {
		q.Add("repository", repository)
	}
	q.Add("page", strconv.Itoa(page))
	q.Add("size", "100")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(cfg.Username, cfg.Token))

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func technologyToPath(s string) (string, error) {
	switch s {
	case "r":
		return "r/", nil
	case "python":
		return "python/", nil
	case "all":
		return "", nil
	default:
		return "", fmt.Errorf("undefined technology %s", s)
	}
}

func ListPackages(client *http.Client, cfg RDepotConfig, repository string, archivedFilter bool, nameFilter string) ([]model.Package, error) {
	return ListGenericPackages[model.Package](client, cfg, repository, archivedFilter, nameFilter)
}

func ListGenericPackages[G model.GenericPackage](client *http.Client, cfg RDepotConfig, repository string, archivedFilter bool, nameFilter string) ([]G, error) {
	body, err := ListPackagesPage(client, cfg, repository, 0)

	if err != nil {
		return nil, err
	}

	var response model.Response[G]
	err = response.Unmarshal(body)
	if err != nil {
		return nil, err
	}

	var pkgs = response.Data.Content
	for page := 1; page <= response.Data.Page.TotalPages; page++ {
		new_body, err := ListPackagesPage(client, cfg, repository, page)
		if err != nil {
			return nil, err
		}
		err = response.Unmarshal(new_body)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, response.Data.Content...)
	}

	if archivedFilter {
		pkgs = model.FilterArchived(pkgs)
	}
	if nameFilter != "" {
		if pkgs, err = model.FilterByName(pkgs, nameFilter); err != nil {
			return nil, err
		}
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

	if cfg.Technology != "python" && cfg.Technology != "r" {
		return msg, fmt.Errorf("invalid technology provided for deleting only Python and R are supported")
	}
	path, err := technologyToPath(cfg.Technology)
	if err != nil {
		return msg, err
	}

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
		cfg.Host+"/api/v2/manager/"+path+"submissions",
		&b)
	if err != nil {
		return msg, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(cfg.Username, cfg.Token))

	res, err := client.Do(req)
	if err != nil {
		return msg, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return msg, fmt.Errorf("bad status: %s", res.Status)
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
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

func createFormGZip(w *multipart.Writer, fieldname, path string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	splitPath := strings.Split(path, "/")
	filename := splitPath[len(splitPath)-1]
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", "application/gzip")
	return w.CreatePart(h)
}
