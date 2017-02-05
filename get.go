package main

import (
	"io"
	"net/http"
	"os"
)

func getArtifact(loc *location) error {
	req, reqErr := http.NewRequest("GET", loc.URL(), nil)
	if reqErr != nil {
		return reqErr
	}

	if loc.creds != nil {
		req.SetBasicAuth(loc.creds.User, loc.creds.Password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(loc.file)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
