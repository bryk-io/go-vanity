package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type VCS int

const (
	GIT VCS = iota
	HG
	SVN
	BZR
)

func (v *VCS) fromString(s string) error {
	switch s {
	case "git":
		*v = GIT
	case "hg":
		*v = HG
	case "svn":
		*v = SVN
	case "bzr":
		*v = BZR
	default:
		return errors.New("unknown")
	}
	return nil
}

func (v VCS) String() string {
	names := [...]string{
		"git",
		"hg",
		"svn",
		"bzr",
	}
	if int(v) < 0 || int(v) >= len(names) {
		return "unknown"
	}
	return names[v]
}

func (v VCS) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v VCS) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return v.fromString(s)
}

func (v VCS) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

func (v VCS) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	return v.fromString(s)
}

// Configuration is compatible with https://github.com/GoogleCloudPlatform/govanityurls
type Server struct {
	// Host name to use in meta tags.
	Host string `json:"host,omitempty" yaml:"host,omitempty"`

	// The amount of time to cache package pages in seconds. Controls the max-age
	// directive sent in the Cache-Control HTTP header. Defaults to 3600.
	CacheMaxAge uint `json:"cache_max_age,omitempty" yaml:"cache_max_age,omitempty"`

	// Path configurations.
	Paths map[string]Path `json:"paths" yaml:"paths"`
}

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

// SCL return a properly formatted "Source Code Link" for the path
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

func New() *Server {
	return &Server{
		CacheMaxAge: 3600,
		Paths:       make(map[string]Path),
	}
}
