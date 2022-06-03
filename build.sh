#!/bin/bash -eu
set -o pipefail

OSES="linux darwin windows"
ARCHS="amd64 arm64"

rm -f build/*
for GOOS in $OSES; do
    for GOARCH in $ARCHS; do
        if [ "${GOOS}-${GOARCH}" == "windows-arm64" ]; then
            continue
        fi
        EXTENSION=""
        if [ "${GOOS}" == "windows" ]; then
            EXTENSION=".exe"
        fi
        echo ">> Building ${GOOS}/${GOARCH}"
        export GOOS
        export GOARCH
        time go build -o "build/Concierge_${GOOS}_${GOARCH}${EXTENSION}"
        echo ""
    done
done