PREFIX=github.com/lox/parfait
VERSION=$(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
FLAGS=-X main.Version=$(VERSION)
ARCHS="linux/amd64 darwin/amd64 windows/amd64"

test:
	go get github.com/kardianos/govendor
	govendor test +local

setup:
	go get github.com/mitchellh/gox
	go get github.com/kardianos/govendor

build:
	mkdir -p build/
	gox -osarch="$(ARCHS)" -ldflags="$(FLAGS)" -output="build/parfait_{{.OS}}_{{.Arch}}" $(PREFIX)

install:
	go install -ldflags="$(FLAGS)" $(PREFIX)
