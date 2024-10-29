package error

import (
	"time"

	"github.com/nik184/urlshortener/internal/app/logger"
)

type Er struct {
	err  error
	time time.Time
}

func (e *Er) Error() string {
	return e.err.Error()
}

func newErr(err error) error {
	var nErr = Er{
		err:  err,
		time: time.Now(),
	}

	nErr.log()

	return &nErr
}

func (e Er) log() {
	logger.Zl.Error(
		e.err.Error(),
		"time | ", e.time,
	)
}
