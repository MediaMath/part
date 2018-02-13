.PHONY:	part publish test clean publish-integrationtest

# Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

TIMESTAMP := $(shell date +"%s")
BUILD_TIME := $(shell date +"%Y%m%d.%H%M%S")
ARTIFACTORY_HOST = artifactory-origin.mediamath.com

VERSION = $(strip $(TIMESTAMP))
ifndef REPOSITORY 
	REPOSITORY = libs-staging-local
endif

ifdef VERBOSE
	TEST_VERBOSITY=-v
else
	TEST_VERBOSITY=
endif

part:
	go get ./...
	go build ./...

tmp:
	mkdir tmp

tmp/part-$(VERSION).zip: part | tmp 
	zip -r -j $@ part 

test:
	go test ./...

publish: part tmp/part-$(VERSION).zip
	part -verbose -credentials=$(HOME)/.ivy2/credentials/artifactory.mediamath.com -h="https://$(ARTIFACTORY_HOST)/artifactory" -r=$(REPOSITORY) -g=com.mediamath -v=$(VERSION) part:tmp/part-$(VERSION).zip

publish-integrationtest: part tmp/part-$(VERSION).zip
	part -verbose -credentials=$(HOME)/.ivy2/credentials/artifactory.mediamath.com -h="https://$(ARTIFACTORY_HOST)/artifactory" -r=$(REPOSITORY) -g=com.mediamath -v=$(VERSION) part:tmp/part-$(VERSION).zip part-integrationtest-1:tmp/part-$(VERSION).zip part-integrationtest-2:tmp/part-$(VERSION).zip part-integrationtest-3:tmp/part-$(VERSION).zip


clean:
	go clean ./...
	rm -rf tmp/*
