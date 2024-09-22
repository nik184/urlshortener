package middlewares

import (
	"net/http"
)

type (
	responceData struct {
		code int
		size int
		body string
	}

	loggerRespWriter struct {
		http.ResponseWriter
		data *responceData
	}
)

func (r *loggerRespWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.data.size += size
	r.data.body += string(b)
	return size, err
}

func (r *loggerRespWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.data.code = statusCode
}
