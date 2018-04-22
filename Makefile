PREFIX=github.com/lox/parfait
VERSION=v$(shell awk -F\" '/const Version/ {print $$2}' version/version.go)
FLAGS=-X main.Version=$(VERSION) -s -w
ARCHS=linux/amd64 darwin/amd64 windows/amd64
test:
	go get github.com/kardianos/govendor
	govendor test +local

setup:
	go get github.com/mitchellh/gox
	go get github.com/kardianos/govendor

clean:
	-rm -rf build/

build:
	mkdir -p build/
	gox -osarch="$(ARCHS)" -ldflags="$(FLAGS)" -output="build/parfait_{{.OS}}_{{.Arch}}" $(PREFIX)

install:
	go install -ldflags="$(FLAGS)" $(PREFIX)
