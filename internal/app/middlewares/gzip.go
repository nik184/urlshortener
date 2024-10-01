package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = encodeResp(w, r)

		next.ServeHTTP(w, r)
	})
}

func encodeResp(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		return w
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "text/html") &&
		!strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return w
	}

	gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	defer gz.Close()

	w.Header().Set("Content-Encoding", "gzip")

	return gzipWriter{ResponseWriter: w, Writer: gz}
}
