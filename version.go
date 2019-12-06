package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

var (
	coreVersion string
	buildCode string
	buildTimestamp string
)

func versionInfo() map[string]string {
	var components = map[string]string{
		"version":    coreVersion,
		"build_code": buildCode,
		"os":    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"go": runtime.Version(),
	}
	if buildTimestamp != "" {
		st, err := strconv.ParseInt(buildTimestamp, 10, 64)
		if err == nil {
			components["release_date"] = time.Unix(st, 0).Format(time.RFC822)
		}
	}
	return components
}
