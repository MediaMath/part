package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os"
)

func publish(pomOnly bool, file string, host string, creds *credentials, repo string, group string, artifact string, version string) (*artifactoryResponse, *artifactoryResponse, error) {

	pomName := fmt.Sprintf("%v.pom", artifact)
	if pomErr := createPom(pomName, group, artifact, version); pomErr != nil {
		return nil, nil, pomErr
	}

	if !pomOnly {
		defer os.RemoveAll(pomName)

		fileResponse := deploy(file, host, creds, repo, group, artifact, version)
		if deployErr := fileResponse.AsError(); deployErr != nil {
			return nil, nil, deployErr
		}

		pomResponse := deploy(pomName, host, creds, repo, group, artifact, version)
		if deployErr := pomResponse.AsError(); deployErr != nil {
			return fileResponse, nil, deployErr
		}

		return fileResponse, pomResponse, nil
	}

	return nil, nil, nil
}
