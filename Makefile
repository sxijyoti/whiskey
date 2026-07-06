BINARY := whiskey
SITE := site

.PHONY: whiskey build run serve test lint clean fullclean install

whiskey:
	@make build
	@make install
	whiskey serve $(SITE)

build:
	@mkdir -p bin
	go build -o bin/$(BINARY) ./cmd/whiskey

run:
	go run ./cmd/whiskey

serve:
	go run ./cmd/whiskey serve $(SITE)

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/whiskey

clean:
	rm -rf $(SITE)/dist
	rm -rf bin

fullclean:
	rm -rf $(SITE)/dist
	rm -rf $(SITE)/.whiskey
	rm -rf bin