package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type cleanup func()

func writeToTempFile(prefix string, name string, content string) (string, cleanup, error) {
	tempdir, dirErr := ioutil.TempDir("", "cred_test")
	if dirErr != nil {
		return "", func() {}, dirErr
	}
	clean := func() { os.RemoveAll(tempdir) }

	fullName := filepath.Join(tempdir, name)
	if writeErr := ioutil.WriteFile(fullName, []byte(content), 0777); writeErr != nil {
		return "", clean, writeErr
	}

	return fullName, clean, nil
}

func TestGetCredentialsFromIniFile(t *testing.T) {
	iniConfig := `
	user=iniUser
	password=iniPassword
`
	config, clean, cfgErr := writeToTempFile("inifile", "example.com", iniConfig)
	defer clean()
	if cfgErr != nil {
		t.Errorf("Error on credential file write: %v", cfgErr)
	}

	creds, err := getCredentials(config)
	if err != nil {
		t.Errorf("Error on credential file read: %v", err)
		t.FailNow()
	}

	if creds.User != "iniUser" || creds.Password != "iniPassword" {
		t.Errorf("Credentials were not read correctly: %v", creds)
	}
}

func TestGetCredentialsFromJsonFile(t *testing.T) {
	jsonConfig := `
{
	"host": "https://example.com",
	"user": "jsonUser",
	"password": "jsonPassword"
}
`
	config, clean, cfgErr := writeToTempFile("jsonfile", "credentials.json", jsonConfig)
	defer clean()
	if cfgErr != nil {
		t.Errorf("Error on credential file write: %v", cfgErr)
	}

	creds, err := getCredentials(config)
	if err != nil {
		t.Errorf("Error on credential file read: %v", err)
		t.FailNow()
	}

	if creds.User != "jsonUser" || creds.Password != "jsonPassword" {
		t.Errorf("Credentials were not read correctly: %v", creds)
	}
}

func TestGetCredentialsFromEnvironment(t *testing.T) {
	os.Setenv("ARTIFACTORY_USERNAME", "envUser")
	defer os.Setenv("ARTIFACTORY_USERNAME", "")
	os.Setenv("ARTIFACTORY_PASSWORD", "envPassword")
	defer os.Setenv("ARTIFACTORY_PASSWORD", "")

	creds, err := getCredentials("")
	if err != nil {
		t.Errorf("Error on credential environment read: %v", err)
		t.FailNow()
	}

	if creds.User != "envUser" || creds.Password != "envPassword" {
		t.Errorf("Environment was not read correctly: %v", creds)
	}
}

func TestCredentialsIniFileIsVerified(t *testing.T) {
	iniConfig := `
	password=bar
`
	config, clean, cfgErr := writeToTempFile("inimissing", "mising.com", iniConfig)
	defer clean()
	if cfgErr != nil {
		t.Errorf("Error on credential file write: %v", cfgErr)
	}

	_, err := getCredentials(config)
	if err == nil {
		t.Errorf("Credentials file read is not verified")
	}
}

func TestCredentialsJsonFileIsVerified(t *testing.T) {
	jsonConfig := `
{
	"host": "https://example.com",
	"password": "jsonPassword"
}
`
	config, clean, cfgErr := writeToTempFile("missing", "credentials.json", jsonConfig)
	defer clean()
	if cfgErr != nil {
		t.Errorf("Error on credential file write: %v", cfgErr)
	}

	_, err := getCredentials(config)
	if err == nil {
		t.Errorf("Credentials file read is not verified")
	}
}

func TestCredentialsNullWhenEnvironmentIsNotFullySet(t *testing.T) {
	os.Setenv("ARTIFACTORY_USERNAME", "envUser")
	defer os.Setenv("ARTIFACTORY_USERNAME", "")

	creds, err := getCredentials("")
	if err != nil || creds != nil {
		t.Errorf("%v:%v", err, creds)
	}
}

func TestMissingItemsInCredentials(t *testing.T) {
	missingUser := &credentials{Password: "foo"}
	if err := verifyCredentials(*missingUser, "bar"); err == nil {
		t.Errorf("Missing user allowed")
	}

	missingPassword := &credentials{User: "foo"}
	if err := verifyCredentials(*missingPassword, "boo"); err == nil {
		t.Errorf("Missing password allowed")
	}

	good := &credentials{User: "foo", Password: "bar"}
	if err := verifyCredentials(*good, "goo"); err != nil {
		t.Errorf("error:%v", err)
	}
}
