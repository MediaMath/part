package main

import (
	"io"
	"net/http"
	"os"
)

func getArtifact(fileName string, host string, creds *credentials, repo string, group string, artifact string, version string) error {
	location := url(fileName, host, repo, group, artifact, version)
	return get(location, fileName, creds)
}

func get(url string, fileName string, creds *credentials) error {
	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		return reqErr
	}

	if creds != nil {
		req.SetBasicAuth(creds.User, creds.Password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
