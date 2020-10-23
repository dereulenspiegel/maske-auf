VERSION                 ?= $(shell git describe --tags --always --dirty)
GIT_COMMIT          		?= $(shell git rev-list -1 HEAD)
RELEASE_VERSION     		?= $(shell git describe --abbrev=0 --tag)
GO_BUILD_ENV_VARS       ?= CGO_ENABLED=0 GO111MODULE=on
LDFLAGS         				?= -X main.Version=$(VERSION) \
													 -X main.Commit=$(GIT_COMMIT) \
													 -w -s

.PHONY: build clean test
.DEFAULT_GOAL := build

dist/frontend:
	mkdir -p dist/frontend
	cp -r ./frontend/* ./dist/frontend/

dist/maske-auf: test dist/frontend
	mkdir -p dist
	packr2
	$(GO_BUILD_ENV_VARS) go build -o dist/maske-auf -ldflags "$(LDFLAGS)"

build: dist/maske-auf

test:
	go test -v -race -covermode=atomic -coverprofile=single.coverprofile

clean:
	rm -rf dist/
	rm -f *.coverprofile
