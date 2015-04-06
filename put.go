package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func deploy(fileName string, host string, creds *credentials, repo string, group string, artifact string, version string) *artifactoryResponse {

	location := url(fileName, host, repo, group, artifact, version)
	resp, putErr := put(location, fileName, creds)
	defer resp.Body.Close()

	if putErr != nil {
		return &artifactoryResponse{Location: location, PublishError: putErr}
	}

	return parseResponse(location, resp)
}

func isOk(status int) bool {
	return status > 199 && status < 300
}

func url(fileName string, host string, repo string, group string, artifact string, version string) string {
	hostEscaped := strings.TrimRight(host, "/")
	groupEscaped := strings.Replace(group, ".", "/", -1)
	ur := fmt.Sprintf("%s/%s/%s/%s/%s/%s", hostEscaped, repo, groupEscaped, artifact, version, filepath.Base(fileName))
	return ur
}

func put(url string, fileName string, creds *credentials) (*http.Response, error) {
	file, openErr := os.Open(fileName)
	if openErr != nil {
		return nil, openErr
	}

	req, reqErr := http.NewRequest("PUT", url, file)
	if reqErr != nil {
		return nil, reqErr
	}

	if creds != nil {
		req.SetBasicAuth(creds.User, creds.Password)
	}

	client := &http.Client{}
	return client.Do(req)
}

func bodyError(resp *http.Response) error {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("COULD NOT READ BODY:%v", err)
	}

	return fmt.Errorf(string(bodyBytes))
}

func parseResponse(location string, resp *http.Response) *artifactoryResponse {
	artResp := &artifactoryResponse{Location: location, Header: resp.Header, StatusCode: resp.StatusCode, StatusMessage: resp.Status}

	if isOk(resp.StatusCode) {
		var pubResult artifactoryPublishResult
		if parseError := json.NewDecoder(resp.Body).Decode(&pubResult); parseError != nil {
			artResp.PublishError = fmt.Errorf("%v:%v", parseError, bodyError(resp))
		}

		artResp.PublishResult = pubResult
	} else {
		artResp.PublishError = bodyError(resp)
	}

	return artResp
}

type artifactoryResponse struct {
	Location      string
	Header        http.Header
	StatusCode    int
	StatusMessage string
	PublishResult artifactoryPublishResult
	PublishError  error
}

type artifactoryPublishResult struct {
	Repo              string               `json:repo`
	Path              string               `json:path`
	Created           string               `json:created`
	CreatedBy         string               `json:createdBy`
	DownloadURI       string               `json:downloadUri`
	MimeType          string               `json:mimeType`
	Size              string               `json:size`
	Checksums         artifactoryChecksums `json:checksums`
	OriginalChecksums artifactoryChecksums `json:originalChecksums`
	URI               string               `json:uri`
}

type artifactoryChecksums struct {
	Sha1 string `json:sha1`
	Md5  string `json:md5`
}

func (resp *artifactoryResponse) AsString(verbose bool) string {
	var lines []string

	if verbose {
		lines = append(lines, resp.Location)
		lines = append(lines, resp.StatusMessage)
		for key, value := range resp.Header {
			lines = append(lines, fmt.Sprintf("%v=%v", key, value))
		}
	}

	if resp.PublishError != nil {
		lines = append(lines, resp.PublishError.Error())
	} else {
		jsonString, marshalErr := json.MarshalIndent(resp.PublishResult, " ", "   ")
		if marshalErr != nil {
			lines = append(lines, marshalErr.Error())
		} else {
			lines = append(lines, string(jsonString))
		}
	}

	return strings.Join(lines, "\n")
}

func (resp *artifactoryResponse) AsError() error {
	if resp.PublishError != nil {
		return fmt.Errorf(resp.AsString(true))
	}

	return nil
}
