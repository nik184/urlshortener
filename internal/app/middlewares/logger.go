package middlewares

import (
	"net/http"
	"time"

	"github.com/nik184/urlshortener/internal/app/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		starttime := time.Now()
		uri := r.URL
		meth := r.Method

		respData := responceData{}
		lw := &loggerRespWriter{
			ResponseWriter: w,
			data:           &respData,
		}
		next.ServeHTTP(lw, r)

		dur := time.Since(starttime)
		logger.Zl.Infoln(
			"req |",
			"uri:", uri,
			"meth:", meth,
			"dur:", dur,
		)

		logger.Zl.Infoln(
			"resp |",
			"size:", lw.data.size,
			"code:", lw.data.code,
			"body:", lw.data.body,
			"dur:", dur,
		)
	})
}
