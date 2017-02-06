package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type location struct {
	host     string
	creds    *credentials
	repo     string
	group    string
	artifact string
	version  string
	file     string
}

func (loc *location) URL() string {
	hostEscaped := strings.TrimRight(loc.host, "/")
	groupEscaped := strings.Replace(loc.group, ".", "/", -1)
	return fmt.Sprintf("%s/%s/%s/%s/%s/%s", hostEscaped, loc.repo, groupEscaped, loc.artifact, loc.version, filepath.Base(loc.file))
}

func forFile(base *location, file string) *location {
	return &location{
		host:     base.host,
		creds:    base.creds,
		repo:     base.repo,
		group:    base.group,
		artifact: base.artifact,
		version:  base.version,
		file:     file,
	}
}
