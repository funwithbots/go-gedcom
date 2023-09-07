GO ?= go

PKGS ?= ./...

BIN := bin/main

.PHONY: cli get vet test

cli:
	mkdir -pv $(dir $(BIN)) && $(GO) build -o $(BIN) ./cli

get:
	$(GO) get $(ARGS) $(PKGS)

vet:
	$(GO) vet $(ARGS) $(PKGS)

test:
	$(GO) test $(ARGS) $(PKGS)
