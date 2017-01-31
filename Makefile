.PHONY:	part publish test clean 

# Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

TIMESTAMP := $(shell date +"%s")
BUILD_TIME := $(shell date +"%Y%m%d.%H%M%S")
ARTIFACTORY_HOST = artifactory.mediamath.com

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
	part -verbose -credentials=$(HOME)/.ivy2/credentials/$(ARTIFACTORY_HOST) -h="https://$(ARTIFACTORY_HOST)/artifactory" -r=$(REPOSITORY) -g=com.mediamath -a=part -v=$(VERSION) tmp/part-$(VERSION).zip
	part -verbose -credentials=$(HOME)/.ivy2/credentials/$(ARTIFACTORY_HOST) -h="https://$(ARTIFACTORY_HOST)/artifactory" -r=$(REPOSITORY) -g=com.mediamath.changes -a=part1 -v=$(VERSION) tmp/part-$(VERSION).zip
	part -verbose -credentials=$(HOME)/.ivy2/credentials/$(ARTIFACTORY_HOST) -h="https://$(ARTIFACTORY_HOST)/artifactory" -r=$(REPOSITORY) -g=com.mediamath.changes -a=part2 -v=$(VERSION) tmp/part-$(VERSION).zip
	part -verbose -credentials=$(HOME)/.ivy2/credentials/$(ARTIFACTORY_HOST) -h="https://$(ARTIFACTORY_HOST)/artifactory" -r=$(REPOSITORY) -g=com.mediamath.changes -a=part3 -v=$(VERSION) tmp/part-$(VERSION).zip
	part -verbose -credentials=$(HOME)/.ivy2/credentials/$(ARTIFACTORY_HOST) -h="https://$(ARTIFACTORY_HOST)/artifactory" -r=$(REPOSITORY) -g=com.mediamath.changes -a=part4 -v=$(VERSION) tmp/part-$(VERSION).zip

clean:
	go clean ./...
	rm -rf tmp/*
