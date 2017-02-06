package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

//examples:
// - ARTIFACTORY_HOST=https://example.com ARTIFACTORY_USER=foo ARTIFACTORY_PASSWORD=bar part -g com.example -a goo -v 1.2 -r local-release goo-jdk1.2.zip
// - part -c example.com.json -g com.example -a goo -v 1.2-SNAPSHOT -r local-snapshot goo.zip

const hostEnvVariable = "ARTIFACTORY_HOST"

var (
	verbose = flag.Bool("verbose", false, "Show verbose output.")
	getFlag = flag.Bool("get", false, "Get the artifact instead of publishing it.")
	pomOnly = flag.Bool("pomOnly", false, "Do NOT publish.  Generate poms only")
	timeout = flag.String("t", "30s", "Client timeout")

	credentialsFile = flag.String("credentials", "", fmt.Sprintf("File with user, password.  If .json extension assumes json otherwise ini.  If not provided assumes %s, %s environment variables are provided.", userEnvVariable, passwordEnvVariable))
	host            = flag.String("h", os.Getenv(hostEnvVariable), fmt.Sprintf("Artifactory REST API endpoint (ie https://artifactory.example.com/artifactory/). If not provided looks at environment variable %s.", hostEnvVariable))
	repo            = flag.String("r", "", "Repository to publish to")
	group           = flag.String("g", "", "Maven group")
	artifact        = flag.String("a", "", "Maven artifact")
	version         = flag.String("v", "", "Maven version")
)

func parseLocations() ([]*location, error) {
	creds, err := getCredentials(*credentialsFile)
	if err != nil {
		return nil, err
	}

	if len(flag.Args()) < 1 {
		return nil, fmt.Errorf("Must provide something to publish")
	}

	if *host == "" ||
		*repo == "" ||
		*group == "" ||
		*version == "" {
		return nil, fmt.Errorf("Must provide all required fields")
	}

	artifacts := make(map[string][]string)
	if *artifact != "" {
		artifacts[*artifact] = flag.Args()
	} else {
		for _, colonDelimited := range flag.Args() {
			pair := strings.Split(colonDelimited, ":")
			if len(pair) != 2 {
				return nil, fmt.Errorf("Could not parse: %v", colonDelimited)
			}

			artifacts[pair[0]] = append(artifacts[pair[0]], pair[1])
		}
	}

	locations := []*location{}
	for artifact, files := range artifacts {
		for _, file := range files {

			loc := &location{}
			loc.creds = creds

			loc.host = *host
			loc.repo = *repo
			loc.group = *group
			loc.version = *version

			loc.artifact = artifact
			loc.file = file

			locations = append(locations, loc)
		}
	}

	return locations, nil
}

func main() {
	flag.Parse()
	locations, err := parseLocations()

	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	timeoutDuration, parseErr := time.ParseDuration(*timeout)
	if parseErr != nil {
		log.Printf("Cannout parse timeout, using 30s: %v", parseErr)
		timeoutDuration = 30 * time.Second
	}

	if *getFlag {
		for _, loc := range locations {
			getErr := getArtifact(loc)
			if getErr != nil {
				log.Fatal(getErr)
			}
		}
	} else {

		for _, loc := range locations {
			fileResponse, pomResponse, publishErr := publish(timeoutDuration, *pomOnly, loc)

			if publishErr != nil {
				log.Fatal(publishErr)
			}

			fmt.Println(fileResponse.AsString(*verbose))
			fmt.Println(pomResponse.AsString(*verbose))
		}
	}
}
