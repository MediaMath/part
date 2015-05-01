#!/bin/bash

set -eu

. /opt/golang/go1.4.2/bin/go_env.sh

export GOPATH="$(pwd)/go"
export PATH="/usr/lib/postgresql/9.1/bin:$GOPATH/bin:$PATH"

cd "./$CLONE_PATH"

go get ./...
make clean check integration-test VERBOSE=yes

if [ "$GH_EVENT_NAME" == "push" -a "$GH_TARGET" == "master" ]; then
    REPOSITORY=libs-release-local make publish
fi
