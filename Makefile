# Makefile for file2qr
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/share/man
GO ?= go
GOFMT ?= gofmt
VERSION := $(shell grep "const VERSION" file2qr.go | cut -d'"' -f2)

.PHONY: all build clean install uninstall fmt test

all: build

build:
	$(GO) build -o file2qr

clean:
	rm -f file2qr

fmt:
	$(GOFMT) -w *.go

test:
	$(GO) test -v ./...

install: build
	install -D -m 755 file2qr $(DESTDIR)$(BINDIR)/file2qr
	install -D -m 644 file2qr.1 $(DESTDIR)$(MANDIR)/man1/file2qr.1
	@echo "Installed file2qr $(VERSION) to $(DESTDIR)$(BINDIR)/file2qr"

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/file2qr
	rm -f $(DESTDIR)$(MANDIR)/man1/file2qr.1
	@echo "Uninstalled file2qr from $(DESTDIR)$(BINDIR)/file2qr"

# Create distribution archive
dist: clean
	mkdir -p file2qr-$(VERSION)
	cp file2qr.go file2qr.1 README.md LICENSE Makefile file2qr-$(VERSION)/
	tar -czf file2qr-$(VERSION).tar.gz file2qr-$(VERSION)
	rm -rf file2qr-$(VERSION)
	@echo "Created file2qr-$(VERSION).tar.gz"