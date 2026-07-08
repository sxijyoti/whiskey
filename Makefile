BINARY := whiskey
SITE := site
VERSION := v0.1.0

.PHONY: whiskey build run serve test lint clean fullclean install

whiskey:
	@make build
	@make install
	whiskey serve $(SITE)

build:
	@echo "Building Whiskey $(VERSION)..."
	@mkdir -p bin
	@go build \
		-ldflags "-X github.com/sxijyoti/whiskey/internal/cli.Version=$(VERSION)" \
		-o bin/$(BINARY) \
		./cmd/whiskey
	@echo "Done."

install:
	@echo "Installing Whiskey $(VERSION)..."
	@go install \
		-ldflags "-X github.com/sxijyoti/whiskey/internal/cli.Version=$(VERSION)" \
		./cmd/whiskey
	@echo "Done."

run:
	go run ./cmd/whiskey

serve:
	@echo "Serving Whiskey..."
	go run ./cmd/whiskey serve $(SITE)

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf $(SITE)/dist
	rm -rf bin

fullclean:
	rm -rf $(SITE)/dist
	rm -rf $(SITE)/.whiskey
	rm -rf bin