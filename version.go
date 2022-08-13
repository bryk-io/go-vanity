package main

import (
	"fmt"
	"runtime"
	"time"
)

var (
	coreVersion    string
	buildCode      string
	buildTimestamp string
)

func versionInfo() map[string]string {
	var components = map[string]string{
		"version":    coreVersion,
		"build_code": buildCode,
		"os":         fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"go":         runtime.Version(),
		"release":    releaseCode(),
	}
	if buildTimestamp != "" {
		rd, err := time.Parse(time.RFC3339, buildTimestamp)
		if err == nil {
			components["release_date"] = rd.Format(time.RFC822)
		}
	}
	return components
}

func releaseCode() string {
	release := "govanity"
	if coreVersion != "" {
		release += "@" + coreVersion
	}
	if buildCode != "" {
		release += "+" + buildCode
	}
	return release
}
