# part
Publish to artifactory as a cli

## What it does

Part is a simple cli that take some command line parameters (part -h for documentation) and an artifact (zip, jar, etc).  It then creates a maven pom and publishes it and the artifact to artifactory.  This allows for easy publishing of artifacts that can be found via Maven GAVC (Group, Artifact, Version, Classifier) style searches via the artifactory rest api.  See the [maven pom reference](https://maven.apache.org/pom.html) for more details about gavc coordinates generally and the [artifactory rest documentation](https://www.jfrog.com/confluence/display/RTF/Artifactory+REST+API#ArtifactoryRESTAPI-GAVCSearch) for artifactory GAVC searches.

## Artifactory REST API endpoint

The artifactory rest api endpoint must be provided either on the command line via the -h flag or via the ARTIFACTORY_HOST parameter.  This should be the fully formed REST API endpoint, not just the host domain. (e.g. https://artifactory.mediamath.com/artifactory).

## Artifactory Authentication

Part allows for the publishing to artifactory via basic auth via one of the following 3 mechanisms.

### JSON File

If a filename is passed in via the -credentials flag with the .json file extension, it will be parsed as a json object expected to be of the form: 

```json
{ 
	"user":"username", 
	"password":"password"
}
```

### INI File

If a filename is passed in via the -credentials flag without the .json file extension, it will be parsed as a limited ini file of the form: 

```ini
user=username
password=password
```

### Environment Variables

If no filename is passed in via the -credentials flag part will look for the ARTIFACTORY_USER and ARTIFACTORY_PASSWORD environment variables for authentication.

### No Auth

If no filename is passed and both environment variables are not provided, part will attempt to publish without authentication.

## Integration Tests

There is an integration test that actually attempts to publish to a real repository.  It is not run in short mode, or if the host and repo are not defined.  To run: 

```bash
ARTIFACTORY_HOST=https://artifactory.mediamath.com/artifactory PART_TEST_REPO=libs-snapshot-local PART_TEST_CREDENTIALS=$HOME/.ivy2/credentials/artifactory.mediamath.com go test -v
```
