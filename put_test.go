package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"testing"
	"time"
)

func TestUrlEncodingDoesGroupPaths(t *testing.T) {
	loc := &location{
		host:     "https://example.com",
		repo:     "repo1",
		group:    "com.mediamath",
		artifact: "foo",
		version:  "121-SNAPSHOT",
		file:     "foo.zip",
	}

	if loc.URL() != "https://example.com/repo1/com/mediamath/foo/121-SNAPSHOT/foo.zip" {
		t.Errorf(loc.URL())
	}
}

func TestUrlEncodingOnlyFileName(t *testing.T) {
	loc := &location{
		host:     "https://example.com",
		repo:     "repo1",
		group:    "com.mediamath",
		artifact: "foo",
		version:  "121-SNAPSHOT",
		file:     "moo/bar/goo/foo.zip",
	}

	if loc.URL() != "https://example.com/repo1/com/mediamath/foo/121-SNAPSHOT/foo.zip" {
		t.Errorf(loc.URL())
	}
}

func TestUrlEncodingStripsHostIfNecessary(t *testing.T) {
	loc := &location{
		host:     "https://example.com/",
		repo:     "repo1",
		group:    "com.mediamath",
		artifact: "foo",
		version:  "121-SNAPSHOT",
		file:     "moo/bar/goo/foo.zip",
	}

	if loc.URL() != "https://example.com/repo1/com/mediamath/foo/121-SNAPSHOT/foo.zip" {
		t.Errorf(loc.URL())
	}
}

func TestDeployNonExistantFileIsErrorNotPanic(t *testing.T) {
	loc := &location{
		host:     "https://artifactory.mediamath.com/artifactory",
		repo:     "libs-release-global",
		group:    "com.mediamath",
		artifact: "part",
		version:  "failing-test",
		file:     "doesntexist",
	}

	resp := deploy(30*time.Second, loc)

	if resp.PublishError == nil {
		t.Errorf("Should have error: %v", resp)
	}
}
