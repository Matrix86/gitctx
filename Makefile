PROJECT_NAME := "gitctx"
TARGET=gitctx
LDFLAGS="-s -w"
PKG_LIST := $(shell go list ./... | grep -v /vendor/)


all: build

test:
	@go test -short ${PKG_LIST}

test-coverage:
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt

lint:
	@golint -set_exit_status ${PKG_LIST}

build: clean
	@mkdir -p bin
	go build -o bin/${TARGET} -v -ldflags=${LDFLAGS} cmd/${TARGET}/*.go

install: build
	go install -ldflags=${LDFLAGS} ./cmd/${TARGET}

clean:
	@rm -rf bin