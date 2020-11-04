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

	resp, err := http.DefaultClient.Do(req)

	switch {
	case err != nil:
		return nil, err
	case resp.StatusCode != 200:
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	default:
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}

}

func SubmitPackage(cfg RDepotConfig, archive string, repository string, replace bool) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fr, err := os.Open(archive)
	if err != nil {
		return err
	}

	if fw, err := createFormGZip(w, "file", archive); err != nil {
		return err
	} else {
		io.Copy(fw, fr)
	}

	if err := w.WriteField("repository", repository); err != nil {
		return err
	}

	if err := w.WriteField("replace", strconv.FormatBool(replace)); err != nil {
		return err
	}

	w.Close()

	req, err := http.NewRequest("POST", "/api/manager/packages/submit", &b)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	return nil
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
