package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const userEnvVariable = "ARTIFACTORY_USER"
const passwordEnvVariable = "ARTIFACTORY_PASSWORD"

type credentials struct {
	User     string
	Password string
}

func getCredentials(fileName string) (*credentials, error) {
	user := os.Getenv(userEnvVariable)
	pass := os.Getenv(passwordEnvVariable)

	if fileName != "" && strings.ToLower(filepath.Ext(fileName)) == ".json" {
		return getCredentialsFromJSONFile(fileName)
	}

	if fileName != "" {
		return getCredentialsFromIniFile(fileName)
	}

	if user != "" && pass != "" {
		return &credentials{user, pass}, nil
	}

	return nil, nil
}

//not a really well defined ini file, doesnt have sections or comment parsing etc.  but will work for our purposes.
func getCredentialsFromIniFile(fileName string) (*credentials, error) {
	fileBytes, readErr := ioutil.ReadFile(fileName)
	if readErr != nil {
		return nil, readErr
	}

	fileString := string(fileBytes)

	lines := strings.Split(fileString, "\n")
	var creds credentials
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		key, value, parseErr := parseIniPair(line)
		if parseErr != nil {
			return nil, parseErr
		}

		if key == "user" {
			creds.User = value
		} else if key == "password" {
			creds.Password = value
		}
	}

	if verifyErr := verifyCredentials(creds, fmt.Sprintf("Config file must contain user and password: %v", fileString)); verifyErr != nil {
		return nil, verifyErr
	}

	return &creds, nil
}

func parseIniPair(line string) (string, string, error) {
	items := strings.Split(line, "=")
	if len(items) != 2 {
		return "", "", fmt.Errorf("Parse error: %v", line)
	}

	return strings.TrimSpace(items[0]), strings.TrimSpace(items[1]), nil
}

func getCredentialsFromJSONFile(fileName string) (*credentials, error) {
	configFile, openErr := os.Open(fileName)
	if openErr != nil {
		return nil, openErr
	}
	defer configFile.Close()

	var creds credentials
	if parseErr := json.NewDecoder(configFile).Decode(&creds); parseErr != nil {
		return nil, parseErr
	}

	if verifyErr := verifyCredentials(creds, `Config file must be { "user": "USERVAL", "password": "PASSVAL"}`); verifyErr != nil {
		return nil, verifyErr
	}

	return &creds, nil
}

func verifyCredentials(creds credentials, errorText string) error {
	if creds.Password == "" || creds.User == "" {
		return fmt.Errorf(errorText)
	}

	return nil
}
