# Nova Project - Makefile(Golang)

# Golang
GO = go

# Binary
BIN = $(shell pwd)/bin

# Sub Directory
SUB_DIR = app

# Build
all: test build

build:
	$(shell mkdir -p $(BIN))
	$(GO) build -o $(BIN)

test:
	for dir in $(SUBDIRS); do \
        $(MAKE) -C $$dir test; \
    done