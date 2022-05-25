.PHONY: build
build:
	go build -v ./cmd/bot

.DEFAULT_GOAL := build