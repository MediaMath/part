package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

func deploy(timeout time.Duration, loc *location) *artifactoryResponse {
	timing, req, resp, putErr := put(timeout, loc)
	return parseResponse(loc.URL(), putErr, resp, req, timing)
}

func isOk(status int) bool {
	return status > 199 && status < 300
}

func put(timeout time.Duration, loc *location) (*artifactoryTiming, *http.Request, *http.Response, error) {
	file, openErr := os.Open(loc.file)
	if openErr != nil {
		return nil, nil, nil, fmt.Errorf("error opening %v: %v", loc.file, openErr)
	}

	req, reqErr := http.NewRequest("PUT", loc.URL(), file)
	if reqErr != nil {
		return nil, nil, nil, fmt.Errorf("error creating PUT request for %v: %v", loc.URL(), reqErr)
	}

	if loc.creds != nil {
		req.SetBasicAuth(loc.creds.User, loc.creds.Password)
	}

	req.Close = true

	r, rerr := httputil.DumpRequest(req, false)
	log.Printf("%v:%s", rerr, r)

	client := &http.Client{Timeout: timeout}
	timing := &artifactoryTiming{Start: time.Now()}
	resp, err := client.Do(req)
	timing.End = time.Now()
	timing.Duration = timing.End.Sub(timing.Start)
	if err != nil {
		err = fmt.Errorf("client error: %v", err)
	}
	return timing, req, resp, err
}

func bodyError(resp *http.Response) error {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("COULD NOT READ BODY:%v", err)
	}

	return fmt.Errorf("body: %q", bodyBytes)
}

func parseResponse(location string, err error, resp *http.Response, req *http.Request, timing *artifactoryTiming) *artifactoryResponse {
	artResp := &artifactoryResponse{Location: location, PublishError: err, Timing: timing}

	if resp == nil {
		return artResp
	}

	defer resp.Body.Close()

	artResp.Header = resp.Header
	artResp.StatusCode = resp.StatusCode
	artResp.StatusMessage = resp.Status

	if isOk(resp.StatusCode) {
		var pubResult artifactoryPublishResult
		if parseError := json.NewDecoder(resp.Body).Decode(&pubResult); parseError != nil {
			artResp.PublishError = fmt.Errorf("could not parse result %v:%v", parseError, bodyError(resp))
		}

		artResp.PublishResult = pubResult
	} else {
		artResp.PublishError = fmt.Errorf("invalid response %v: %v", resp.StatusCode, bodyError(resp))
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
	Timing        *artifactoryTiming
}

type artifactoryPublishResult struct {
	Repo              string               `json:"repo"`
	Path              string               `json:"path"`
	Created           string               `json:"created"`
	CreatedBy         string               `json:"createdBy"`
	DownloadURI       string               `json:"downloadUri"`
	MimeType          string               `json:"mimeType"`
	Size              string               `json:"size"`
	Checksums         artifactoryChecksums `json:"checksums"`
	OriginalChecksums artifactoryChecksums `json:"originalChecksums"`
	URI               string               `json:"uri"`
}

type artifactoryTiming struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

type artifactoryChecksums struct {
	Sha1 string `json:"sha1"`
	Md5  string `json:"md5"`
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

	if verbose && resp.Timing != nil {
		lines = append(lines, fmt.Sprintf("start: %v end: %v duration: %v", resp.Timing.Start, resp.Timing.End, resp.Timing.Duration))
	}

	return strings.Join(lines, "\n")
}

func (resp *artifactoryResponse) AsError() error {
	if resp.PublishError != nil {
		return fmt.Errorf(resp.AsString(true))
	}

	return nil
}
