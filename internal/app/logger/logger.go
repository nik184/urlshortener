package logger

import (
	"os"
	"path"
	"runtime"

	"github.com/nik184/urlshortener/internal/app/util"
	"go.uber.org/zap"
)

var Zl *zap.SugaredLogger
var filename string = "main.log"

func init() {
	cfg := zap.NewProductionConfig()

	_, b, _, _ := runtime.Caller(0)
	rootDir := path.Dir(path.Dir(path.Dir(path.Dir(b))))
	path := path.Join(rootDir, filename)

	if exists, err := util.FileExists(path); err != nil || exists {
		_, err := os.Create(path)

		if err != nil {
			panic(err)
		}
	}

	cfg.OutputPaths = []string{path}
	logger, err := cfg.Build()
	// logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Zl = logger.Sugar()
}
