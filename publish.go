package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"os"
	"time"
)

func publish(timeout time.Duration, pomOnly bool, loc *location) (*artifactoryResponse, *artifactoryResponse, error) {

	pomName, pomErr := createPom(loc)
	if pomErr != nil {
		return nil, nil, pomErr
	}

	if !pomOnly {
		defer os.RemoveAll(pomName)

		fileResponse := deploy(timeout, loc)
		if deployErr := fileResponse.AsError(); deployErr != nil {
			return nil, nil, deployErr
		}

		pomResponse := deploy(timeout, forFile(loc, pomName))
		if deployErr := pomResponse.AsError(); deployErr != nil {
			return fileResponse, nil, deployErr
		}

		return fileResponse, pomResponse, nil
	}

	return nil, nil, nil
}
