BINARY := whiskey

.PHONY: whiskey build run serve test lint clean distclean install

whiskey:
	@make build
	@make install
	whiskey serve

build:
	@mkdir -p bin
	go build -o bin/$(BINARY) ./cmd/whiskey

run:
	go run ./cmd/whiskey

serve:
	go run ./cmd/whiskey serve

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/whiskey

clean:
	rm -rf dist
	rm -rf bin

fullclean:
	rm -rf dist
	rm -rf bin
	rm -rf .whiskey