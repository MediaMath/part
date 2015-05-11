package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"os"
	"testing"
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

	_, _, shouldWork := publish(false, "integration_test.go", host, creds, repo, "com.mediamath", "part_integration_test", "SNAPSHOT")

	if shouldWork != nil {
		t.Errorf("%v", shouldWork)
	}

	_, _, shouldFail := publish(false, "integration_test.go", host, creds, "NOTREAL-SHOULDFAIL-REPO", "com.mediamath", "part_integration_test", "SNAPSHOT")

	if shouldFail == nil {
		t.Errorf("Didn't fail")
	}
}
