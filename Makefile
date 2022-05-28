.PHONY: build
build:
	go build -v ./cmd/bot
	./bot

.DEFAULT_GOAL := build