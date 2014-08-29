test:
	go test -v .

cover:
	go test -cover .

format:
	gofmt -l -w *.go

.PHONY: cover test format
