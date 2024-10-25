package storage

type URLWithShort struct {
	URL   string
	Short string
}

type stor interface {
	SetBatch(banch []URLWithShort) error
	Set(url, short string) error
	Get(id string) (string, error)
}
