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
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	//"log"
	"net/http"
	"net/textproto"
	//"net/http/httptest"
	//"net/http/httputil"
)

type RDepotConfig struct {
	Host  string
	Token string
}

func ListPackages(cfg RDepotConfig) ([]byte, error) {

	req, err := http.NewRequest(
		"GET",
		cfg.Host+"/api/manager/packages/list",
		nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	res, err := http.DefaultClient.Do(req)

	switch {
	case err != nil:
		return nil, err
	case res.StatusCode != 200:
		defer res.Body.Close()
		return ioutil.ReadAll(res.Body)
	default:
		defer res.Body.Close()
		return ioutil.ReadAll(res.Body)
	}

}

func SubmitPackage(cfg RDepotConfig, archive string, repository string, replace bool) ([]byte, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fr, err := os.Open(archive)
	if err != nil {
		return nil, err
	}

	if fw, err := createFormGZip(w, "file", archive); err != nil {
		return nil, err
	} else {
		io.Copy(fw, fr)
	}

	if err := w.WriteField("repository", repository); err != nil {
		return nil, err
	}

	if err := w.WriteField("replace", strconv.FormatBool(replace)); err != nil {
		return nil, err
	}

	w.Close()

	req, err := http.NewRequest(
		"POST",
		cfg.Host+"/api/manager/packages/submit",
		&b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)

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
