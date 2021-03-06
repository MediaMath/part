package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/xml"
	"testing"
)

func TestPomCreation(t *testing.T) {
	expected := `
	<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apace.org/POM/4.0.0 http://mave.apache.org/xsd/maven-4.0.0.xsd">
		<modelVersion>4.0.0</modelVersion>
		<groupId>com.mediamath</groupId>
		<artifactId>part</artifactId>
		<version>12-SNAPSHOT</version>
	</project>`

	var parsed pom
	xml.Unmarshal([]byte(expected), &parsed)

	loc := &location{
		host:     "https://example.com",
		repo:     "fake",
		group:    "com.mediamath",
		artifact: "part",
		version:  "12-SNAPSHOT",
		file:     "foo.zip",
	}

	created := newPom(loc)

	if parsed != *created {
		t.Errorf("\n%v\n%v", parsed, *created)
	}
}
