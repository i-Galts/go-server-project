BINARY=server
BUILDDIR=build
MAINPKG=./cmd/server

.PHONY: build clean

.DEFAULT_GOAL := build

build:
	@go build -v -o ${BUILDDIR}/${BINARY} ${MAINPKG}

clean:
	@rm -rf ${BUILDDIR}
	@mkdir -p ${BUILDDIR}