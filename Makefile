test:
	@go test -v ./...

build:
	@go build -o bin/blocker

run: build
	./bin/blocker