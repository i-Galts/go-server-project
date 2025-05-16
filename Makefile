LB_BIN=lb
BE_BIN=backend
CLIENT_BIN=client

BUILDDIR=build

LB_SRC=./cmd/loadbalancer
BE_SRC=./cmd/backend
CLIENT_SRC=./cmd/client

.PHONY: build clean

.DEFAULT_GOAL := build

build:
	@go build -v -o ${BUILDDIR}/${LB_BIN} ${LB_SRC}
	cp ./configs/lb_conf.json ${BUILDDIR}

clean:
	@rm -rf ${BUILDDIR}
	@mkdir -p ${BUILDDIR}