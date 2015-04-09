package main

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

	created := newPom("com.mediamath", "part", "12-SNAPSHOT")

	if parsed != *created {
		t.Errorf("\n%v\n%v", parsed, *created)
	}
}
