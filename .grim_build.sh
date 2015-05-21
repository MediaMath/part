#!/bin/bash

set -eu

. /opt/golang/go1.4.2/bin/go_env.sh

export GOPATH="$(pwd)/go"
export PATH="$GOPATH/bin:$PATH"

cd "./$CLONE_PATH"

make clean test

if [ "$GH_EVENT_NAME" == "push" -a "$GH_TARGET" == "master" ]; then
	REPOSITORY=libs-release-local make publish
fi
