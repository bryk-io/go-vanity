// nolint: gocognit
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"gopkg.in/yaml.v2"
)

var sampleConf string = `
host: go.bryk.io
cache_max_age: 3600
paths:
  sample:
    repo: https://github.com/bryk-io/sample
    vcs: git
`

func TestNewServerConfig(t *testing.T) {
	// Empty configuration
	conf := NewServerConfig()
	if len(conf.Paths) != 0 {
		t.Error("configuration should be empty by default")
	}

	// Load YAML content
	if err := yaml.Unmarshal([]byte(sampleConf), conf); err != nil {
		t.Error(err)
	}
	if conf.Paths["sample"].VCS.String() != "git" {
		t.Error("failed to decode VCS value")
	}

	// Re-encode in JSON format
	_, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		t.Error(err)
	}
}

func TestServer(t *testing.T) {
	// Load sample configuration
	conf := NewServerConfig()
	_ = yaml.Unmarshal([]byte(sampleConf), conf)

	// Dummy build values
	coreVersion = "0.1.0"
	buildCode = "foo-bar"

	// Start test server
	h := newHandler(conf)
	server := http.Server{Addr: ":9091", Handler: getServerMux(h)}
	go func() {
		_ = server.ListenAndServe()
	}()

	t.Run("api/version", func(t *testing.T) {
		res, err := http.Get("http://localhost:9091/api/version")
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			_ = res.Body.Close()
		}()

		if res.Header.Get("content-type") != "application/json" {
			t.Error("invalid content type")
		}
		if res.StatusCode != http.StatusOK {
			t.Error("invalid status code")
		}
		if res.Header.Get("X-Content-Type-Options") != "nosniff" {
			t.Error("invalid content options")
		}
		if res.Header.Get("X-Go-Vanity-Server-Version") != "0.1.0" {
			t.Error("missing version header")
		}
		if res.Header.Get("X-Go-Vanity-Server-Build") != "foo-bar" {
			t.Error("missing build code header")
		}
	})

	t.Run("api/conf", func(t *testing.T) {
		res, err := http.Get("http://localhost:9091/api/conf")
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			_ = res.Body.Close()
		}()
		if res.Header.Get("content-type") != "application/json" {
			t.Error("invalid content type")
		}
		if res.StatusCode != http.StatusOK {
			t.Error("invalid status code")
		}
	})

	t.Run("api/ping", func(t *testing.T) {
		res, err := http.Get("http://localhost:9091/api/ping")
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			_ = res.Body.Close()
		}()
		if res.StatusCode != http.StatusOK {
			t.Error("invalid status code")
		}
		response, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		if fmt.Sprintf("%s", response) != "pong" {
			t.Error("invalid response")
		}
	})

	t.Run("index.html", func(t *testing.T) {
		res, err := http.Get("http://localhost:9091/index.html")
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			_ = res.Body.Close()
		}()
		if res.StatusCode != http.StatusOK {
			t.Error("invalid status code")
		}
	})

	t.Run("invalid-path", func(t *testing.T) {
		res, err := http.Get("http://localhost:9091/invalid-path")
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			_ = res.Body.Close()
		}()
		if res.StatusCode != http.StatusNotFound {
			t.Error("invalid status code")
		}
		response, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		if fmt.Sprintf("%s", response) != "unknown path" {
			t.Error("invalid response")
		}
	})

	t.Run("valid-path", func(t *testing.T) {
		res, err := http.Get("http://localhost:9091/sample")
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			_ = res.Body.Close()
		}()
		if res.StatusCode != http.StatusOK {
			t.Error("invalid status code")
		}
	})

	defer func() {
		_ = server.Close()
	}()
}
