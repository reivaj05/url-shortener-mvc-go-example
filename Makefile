GOCMD = go
PKG   = .
BIN   = url_shortener

.PHONY: %

default: fmt deps test build

all: build
build: deps
	$(GOCMD) build -a -o $(BIN) $(PKG)
fmt:
	$(GOCMD) fmt $(PKG)
test:
	$(GOCMD) test -v ./...
deps:
	wget https://raw.githubusercontent.com/pote/gpm/v1.4.0/bin/gpm
	chmod +x gpm
	./gpm
	rm gpm
