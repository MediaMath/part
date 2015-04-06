package main

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
