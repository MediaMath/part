package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/xml"
	"io/ioutil"
)

type pom struct {
	XMLName        xml.Name `xml:"http://maven.apache.org/POM/4.0.0 project"`
	SchemaLocation string   `xml:"http://www.w3.org/2001/XMLSchema-instance schemaLocation,attr"`
	ModelVersion   string   `xml:"modelVersion"`
	GroupID        string   `xml:"groupId"`
	ArtifactID     string   `xml:"artifactId"`
	Version        string   `xml:"version"`
}

func newPom(group string, artifact string, version string) *pom {
	return &pom{
		//dont like this duplication but can't figure out how to get around it
		XMLName:        xml.Name{"http://maven.apache.org/POM/4.0.0", "project"},
		SchemaLocation: "http://maven.apace.org/POM/4.0.0 http://mave.apache.org/xsd/maven-4.0.0.xsd",
		ModelVersion:   "4.0.0",
		GroupID:        group,
		ArtifactID:     artifact,
		Version:        version}
}

func createPom(filename string, group string, artifact string, version string) error {
	p := newPom(group, artifact, version)

	pomString, marshalErr := xml.MarshalIndent(p, " ", "   ")
	if marshalErr != nil {
		return marshalErr
	}

	if writeErr := ioutil.WriteFile(filename, []byte(pomString), 0644); writeErr != nil {
		return writeErr
	}

	return nil
}
