package storage

type stor interface {
	SetBatch(banch []ShortenURLRow) error
	Set(url, short string) error
	GetByShort(short string) (*ShortenURLRow, error)
	GetByURL(url string) (*ShortenURLRow, error)
}

type ShortenURLRow struct {
	UUID  string `json:"uuid"`
	Short string `json:"shorten_url"`
	URL   string `json:"original_url"`
}

type NotUniqErr struct {
	error
	OldVal *ShortenURLRow
}

func (e NotUniqErr) Error() string {
	return e.error.Error()
}

func NewNotUniqErr(err error, oldVal *ShortenURLRow) error {
	return &NotUniqErr{
		error:  err,
		OldVal: oldVal,
	}
}
