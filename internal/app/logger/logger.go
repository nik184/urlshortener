package logger

import "go.uber.org/zap"

var Zl *zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Zl = logger.Sugar()
}
