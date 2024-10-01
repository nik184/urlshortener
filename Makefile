test:
	go test internal/app/handlers/*.go -v /
	go test internal/app/storage/*.go -v

shortener_compile:
	go build -o cmd/shortener/shortener cmd/shortener/*.go