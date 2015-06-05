package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import "testing"

func TestUrlEncodingDoesGroupPaths(t *testing.T) {
	gotten := url("foo.zip", "https://example.com", "repo1", "com.mediamath", "foo", "121-SNAPSHOT")
	if gotten != "https://example.com/repo1/com/mediamath/foo/121-SNAPSHOT/foo.zip" {
		t.Errorf(gotten)
	}
}

func TestUrlEncodingOnlyFileName(t *testing.T) {
	gotten := url("moo/bar/goo/foo.zip", "https://example.com", "repo1", "com.mediamath", "foo", "121-SNAPSHOT")
	if gotten != "https://example.com/repo1/com/mediamath/foo/121-SNAPSHOT/foo.zip" {
		t.Errorf(gotten)
	}
}

func TestUrlEncodingStripsHostIfNecessary(t *testing.T) {
	gotten := url("moo/bar/goo/foo.zip", "https://example.com/", "repo1", "com.mediamath", "foo", "121-SNAPSHOT")
	if gotten != "https://example.com/repo1/com/mediamath/foo/121-SNAPSHOT/foo.zip" {
		t.Errorf(gotten)
	}
}

func TestDeployNonExistantFileIsErrorNotPanic(t *testing.T) {
	resp := deploy("doesntexist", "https//artifactory.mediamath.com/artifactory", &credentials{}, "libs-release-global", "com.mediamath", "part", "failing-test")

	if resp.PublishError == nil {
		t.Errorf("Should have error: %v", resp)
	}
}
