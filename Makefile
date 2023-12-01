pre:
	go mod tidy
	mkdir -p ./target

dev:
	go run ./cmd/momod

build: pre
	go build -o ./target ./cmd/momod

build_x64: pre
	GOOS=linux GOARCH=amd64 go build -o ./target/momod_linux_amd64 ./cmd/momod

.PHONY: pre dev build build_x64