# Go Vanity
[![Build Status](https://drone.bryk.io/api/badges/bryk-io/go-vanity/status.svg)](https://drone.bryk.io/bryk-io/go-vanity)
[![Version](https://img.shields.io/github/tag/bryk-io/go-vanity.svg)](https://github.com/bryk-io/go-vanity/releases)
[![Software License](https://img.shields.io/badge/license-BSD3-red.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bryk-io/go-vanity?style=flat)](https://goreportcard.com/report/github.com/bryk-io/go-vanity)

Resolve custom import paths for go dependencies hosted on various source code repositories
by specifying a simple configuration file.

The following setup permits to import a package hosted on `https://github.com/example-com/my-repo`
as `import example.com/my-lib` with the following configuration on the values file.

```yaml
configuration:
  host: example.com
  cache_max_age: 3600
  paths:
    my-lib:
      repo: https://github.com/example-com/my-repo
      vcs: git
```

More information on the official go documentation: [Remote Import Paths](https://golang.org/cmd/go/#hdr-Remote_import_paths)