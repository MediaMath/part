package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

//examples:
// - ARTIFACTORY_HOST=https://example.com ARTIFACTORY_USER=foo ARTIFACTORY_PASSWORD=bar part -g com.example -a goo -v 1.2 -r local-release goo-jdk1.2.zip
// - part -c example.com.json -g com.example -a goo -v 1.2-SNAPSHOT -r local-snapshot goo.zip

const hostEnvVariable = "ARTIFACTORY_HOST"

func main() {
	verbose := flag.Bool("verbose", false, "Show verbose output.")
	pomOnly := flag.Bool("pomOnly", false, "Do NOT publish.  Generate poms only")
	credentialsFile := flag.String("credentials", "", fmt.Sprintf("File with user, password.  If .json extension assumes json otherwise ini.  If not provided assumes %s, %s environment variables are provided.", userEnvVariable, passwordEnvVariable))
	host := flag.String("h", os.Getenv(hostEnvVariable), fmt.Sprintf("Artifactory REST API endpoint (ie https://artifactory.example.com/artifactory/). If not provided looks at environment variable %s.", hostEnvVariable))
	repo := flag.String("r", "", "Repository to publish to")
	group := flag.String("g", "", "Maven group")
	artifact := flag.String("a", "", "Maven artifact")
	version := flag.String("v", "", "Maven version")
	flag.Parse()

	if *host == "" ||
		*repo == "" ||
		*group == "" ||
		*artifact == "" ||
		*version == "" {
		flag.PrintDefaults()
		log.Fatal("Must provide the host, repo, group, artifact and version")
	}

	if len(flag.Args()) != 1 {
		log.Fatal("Must provide a file to publish")
	}

	file := flag.Args()[0]

	creds, credErr := getCredentials(*credentialsFile)
	if credErr != nil {
		log.Fatal(credErr)
	}

	fileResponse, pomResponse, publishErr := publish(*pomOnly, file, *host, creds, *repo, *group, *artifact, *version)

	if publishErr != nil {
		log.Fatal(publishErr)
	}

	fmt.Print(fileResponse.AsString(*verbose))
	fmt.Print(pomResponse.AsString(*verbose))
}