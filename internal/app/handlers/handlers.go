package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nik184/urlshortener/internal/app/config"
)

func readBody(rw http.ResponseWriter, r *http.Request) ([]byte, error) {

	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}

		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	return io.ReadAll(reader)
}

func isURLValid(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)

	return err == nil && parsedURL.Host != "" && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func concatPathToAddr(uri string) (result string) {
	result = config.BaseURL + "/" + uri
	if !strings.Contains(config.BaseURL, "") {
		result = "http://" + result
	}

	return
}
