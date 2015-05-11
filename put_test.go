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
