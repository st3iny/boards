INSTALL_DIR ?= /usr/local/bin

all: build

build:
	go build cmd/boards.go

clean:
	go clean
	rm ./boards

install: build
	install ./boards ${INSTALL_DIR}

.PHONY: all build clean install
