#!/usr/bin/env bash

set -o errexit
set -o nounset

if [ -z "${APPNAME}" ]; then
    echo "APPNAME must be set"
    exit 1
fi

if [ -z "${GOARCH}" ]; then
    echo "Using default GOARCH"
else
    echo "Using GOARCH=${GOARCH}"
fi

if [ -z "${GOOS}" ]; then
    echo "Using default GOOS"
else
    echo "Using GOOS=${GOOS}"
fi

if [ -z "${CONF_DIR}" ]; then
    echo "Using default CONF_DIR"
else
    echo "Using CONF_DIR=${CONF_DIR}"
fi

export CGO_ENABLED=0

echo "Go building app"
go build -o build/${APPNAME} *.go
echo "Successfully built, exiting build script"