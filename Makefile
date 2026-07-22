BINARY := whiskey
SITE := site
VERSION := v0.1.0

.PHONY: whiskey build run serve test lint clean fullclean install release-check release-snapshot release completions

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
	rm -rf dist
	rm -rf completions

fullclean:
	rm -rf $(SITE)/dist
	rm -rf $(SITE)/.whiskey
	rm -rf bin
	rm -rf dist
	rm -rf completions

docs:
	go build -o bin/whiskey ./cmd/whiskey
	./bin/whiskey build site

completions:
	@echo "Generating shell completions..."
	@mkdir -p completions
	@go run ./cmd/whiskey completion bash > completions/whiskey.bash
	@go run ./cmd/whiskey completion zsh > completions/whiskey.zsh
	@go run ./cmd/whiskey completion fish > completions/whiskey.fish

release-check:
	@echo "Checking GoReleaser configuration..."
	goreleaser check

release-snapshot: completions
	@echo "Building GoReleaser snapshot..."
	goreleaser release --snapshot --clean

release: completions
	@echo "Running GoReleaser build & release..."
	@if [ -z "$$GH_TOKEN" ] && [ -z "$$GITHUB_TOKEN" ]; then \
		echo "Warning: Neither GH_TOKEN nor GITHUB_TOKEN is set. Real release will fail."; \
	fi
	goreleaser release --clean