BINDIR	:= bin/
ifeq ($(OS),Windows_NT)
	BINNAME	?= github2orgmode.exe
else
	BINNAME	?= github2orgmode
endif
INSTALL_PATH ?= /usr/local/bin

SHELL      = /usr/bin/env bash

GIT_COMMIT = $(shell git rev-parse HEAD)
ifneq ($(GIT_TAG),)
	GIT_TAG := $(GIT_TAG)
else
	GIT_TAG = $(shell git describe --tags 2>/dev/null)
endif

# go option
PKG        := ./...
PKGTESTS   := $$(go list ./... | grep -v /third-party/)
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=

LDFLAGS += $(EXT_LDFLAGS)

.PHONY: all
all: build

.PHONY: build
build: lint $(BINDIR)$(BINNAME)

# Rebuild the binary if any of these files change
SRC := $(shell find . -type f -name '*.go' -print) go.mod go.sum

$(BINDIR)$(BINNAME): $(SRC)
	go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BINDIR)$(BINNAME) ./cmd/github2orgmode

.PHONY: install
install: build
	@install "$(BINDIR)$(BINNAME)" "$(INSTALL_PATH)/$(BINNAME)"

.PHONY: test
test: lint shellcheck build
test: TESTFLAGS += -race -v
test: test-style
test: test-unit

.PHONY: test-unit
test-unit:
	@echo "==> Running unit tests <=="
	go test $(GOFLAGS) -run $(TESTS) $(PKGTESTS) $(TESTFLAGS)

.PHONY: test-style
test-style:
	@echo "==> Checking style <=="
	golangci-lint run

.PHONY: coverage
coverage:
	@echo "==> Running coverage tests <=="
	go test $(GOFLAGS) -run $(TESTS) $(PKGTESTS) -coverprofile=coverage.out --covermode=atomic

.PHONY: lint
lint: fmt vet

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: shellcheck
shellcheck:
	@echo "==> Checking shell scripts <=="
	shellcheck ./scripts/*
	@echo "==> Done <=="

.PHONY: clean
clean:
	rm $(BINDIR)$(BINNAME)
	rmdir $(BINDIR)
