PREFIX=github.com/lox/parfait
VERSION=$(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
FLAGS=-X main.Version=$(VERSION)

test:
	go get github.com/kardianos/govendor
	govendor test +local

setup:
	go get github.com/mitchellh/gox
	go get github.com/kardianos/govendor

build:
	go build -ldflags="$(FLAGS)" $(PREFIX)

install:
	go install -ldflags="$(FLAGS)" $(PREFIX)
