package main

import (
	"fmt"
	"strings"
)

// VCS represents a specific version control system to match
// a path against.
type VCS string

const (
	// GIT source control management (https://git-scm.com/).
	GIT VCS = "git"
	// HG stands for Mercurial (https://www.mercurial-scm.org/wiki/HgSubversion).
	HG VCS = "hg"
	// SVN stands for Subversion (https://subversion.apache.org/).
	SVN VCS = "svn"
	// BZR stands for Bazaar (https://bazaar.canonical.com/en/).
	BZR VCS = "bzr"
)

// Configuration provides the required server parameters.
// Compatible with https://github.com/GoogleCloudPlatform/govanityurls
type Configuration struct {
	// Host name to use in meta tags.
	Host string `json:"host,omitempty" yaml:"host,omitempty"`

	// The amount of time to cache package pages in seconds. Controls the max-age
	// directive sent in the Cache-Control HTTP header. Defaults to 3600.
	CacheMaxAge uint `json:"cache_max_age,omitempty" yaml:"cache_max_age,omitempty"`

	// Path configurations.
	Paths map[string]Path `json:"paths" yaml:"paths"`
}

// Path matches a specific vanity URL to an external version control system.
type Path struct {
	// Root URL of the repository as it would appear in go-import meta tag.
	Repo string `json:"repo" yaml:"repo"`

	// Version control system used.
	VCS VCS `json:"vcs" yaml:"vcs"`

	// The last three fields of the go-source meta tag.
	// home:      URL of the repository's home page
	// directory: URL template for a page listing the files in the package
	// file:      URL template for a link to a line in a source file
	Display string `json:"display,omitempty" yaml:"display,omitempty"`
}

// SCL return a properly formatted "Source Code Link" for the path.
// https://github.com/golang/gddo/wiki/Source-Code-Links
func (p Path) SCL() string {
	if p.Display != "" {
		return fmt.Sprintf("%s %s", p.Repo, p.Display)
	}
	if strings.HasSuffix(p.Repo, "https://github.com/") {
		return fmt.Sprintf("%s %s/tree/master{/dir} %s/blob/master{/dir}/{file}#L{line}", p.Repo, p.Repo, p.Repo)
	}
	if strings.HasSuffix(p.Repo, "https://bitbucket.org/") {
		return fmt.Sprintf("%s %s/src/default{/dir} %s/src/default{/dir}/{file}#{file}-{line}", p.Repo, p.Repo, p.Repo)
	}
	return fmt.Sprintf("%s %s/tree{/dir} %s/blob{/dir}/{file}#L{line}", p.Repo, p.Repo, p.Repo)
}

// NewServerConfig returns an empty server configuration instance.
func NewServerConfig() *Configuration {
	return &Configuration{
		CacheMaxAge: 3600,
		Paths:       make(map[string]Path),
	}
}
