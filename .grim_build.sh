#!/bin/bash

set -eu

. /opt/golang/go1.5.1/bin/go_env.sh

export GOPATH="$(pwd)/go"
export PATH="$GOPATH/bin:$PATH"

cd "./$CLONE_PATH"

if [ "$GH_EVENT_NAME" == "push" -a "$GH_TARGET" == "master" ]; then
	#on merge of master publish part to release artifactory repo
	REPOSITORY=libs-release-global make clean test publish
elif [ "$GH_EVENT_NAME" == "pull_request" -a "$GH_TARGET" == "master" ]; then
	#on any other event publish to the staging repo as this acts as an integration test
	#boostrapping ftw!
	REPOSITORY=libs-staging-local make clean test publish
else 
	make clean test
fi
