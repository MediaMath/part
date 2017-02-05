package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"os"
	"testing"
	"time"
)

const testRepoEnvVariable = "PART_TEST_REPO"
const testCredentialsFileEnvVar = "PART_TEST_CREDENTIALS"

func TestPublishIntegration(t *testing.T) {
	host := os.Getenv(hostEnvVariable)
	repo := os.Getenv(testRepoEnvVariable)

	if testing.Short() {
		t.Skipf("Skipping integration test publish in short mode.")
	}

	if host == "" {
		t.Skipf("Skipping integration test publish as we don't have a host available via %s.", hostEnvVariable)

	}

	if repo == "" {
		t.Skipf("Skipping integration test publish as we don't have a repo available via %s.", testRepoEnvVariable)
	}

	creds, credsErr := getCredentials(os.Getenv(testCredentialsFileEnvVar))
	if credsErr != nil {
		t.Errorf("%v", credsErr)
	}

	loc := &location{
		host:     host,
		creds:    creds,
		repo:     repo,
		group:    "com.mediamath",
		artifact: "part_integration_test",
		version:  "SNAPSHOT",
		file:     "integration_test.go",
	}

	_, _, shouldWork := publish(30*time.Second, false, loc)

	if shouldWork != nil {
		t.Errorf("%v", shouldWork)
	}

	loc = &location{
		host:     host,
		creds:    creds,
		repo:     "NOTREAL-SHOULDFAIL-REPO",
		group:    "com.mediamath",
		artifact: "part_integration_test",
		version:  "SNAPSHOT",
		file:     "integration_test.go",
	}
	_, _, shouldFail := publish(30*time.Second, false, loc)

	if shouldFail == nil {
		t.Errorf("Didn't fail")
	}
}
