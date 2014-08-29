test:
	go test -v .

cover:
	go test -cover .

format:
	gofmt -l -w *.go

vet:
	go tool vet -v *.go

.PHONY: cover test format vet
