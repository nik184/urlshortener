package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logger(next http.Handler) http.Handler {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

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
		sugar.Infoln(
			"req |",
			"uri:", uri,
			"meth:", meth,
			"dur:", dur,
		)

		sugar.Infoln(
			"resp |",
			"size:", lw.data.size,
			"code:", lw.data.code,
			"dur:", dur,
		)
	})
}
