.PHONY: cli-build
cli-build:
	go build ./cmd/cli

.PHOHY: cli-run
cli-run:
	go run ./cmd/cli
