APPNAME ?= parrot

export PORT= 8080
export DATABASE_URL= postgres://bird:docker@localhost:5433/parrot

build:
	mkdir -p build
	CONF_DIR=$(CONF_DIR) GOOS=$(GOOS) GOARCH=$(GOARCH) APPNAME=$(APPNAME) ./scripts/build.sh

run: build
	./build/${APPNAME}

test:
	./scripts/unit-test.sh

integration-test:
	./scripts/integration-test.sh

lint:
	./scripts/lint.sh

docker-build:
	./scripts/docker-build.sh

docker-run:
	docker run parrot -d -e

.PHONY: build run test travis-build lint
