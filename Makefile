APPNAME ?= parrot

# used by `lint` and `test` targets
export REPORTS_DIR=./reports

# used by `rpm` target
export RPM_VERSION=$(shell date +"%Y%m%d%H%M%S")

export CONF_DIR = ./conf

build:
	mkdir -p build
	CONF_DIR=$(CONF_DIR) GOOS=$(GOOS) GOARCH=$(GOARCH) APPNAME=$(APPNAME) ./scripts/build.sh

run: build
	./build/${APPNAME}

test:
	./scripts/go.test.sh

lint:
	./scripts/lint.sh

.PHONY: build run test travis-build lint
