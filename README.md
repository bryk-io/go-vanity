# go-vanity
[![Build Status](https://drone.bryk.io/api/badges/bryk-io/go-vanity/status.svg)](https://drone.bryk.io/bryk-io/go-vanity)
[![Version](https://img.shields.io/github/tag/bryk-io/go-vanity.svg)](https://github.com/bryk-io/go-vanity/releases)
[![Software License](https://img.shields.io/badge/license-BSD3-red.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bryk-io/go-vanity?style=flat)](https://goreportcard.com/report/github.com/bryk-io/go-vanity)

This `go-vanity` tool provides a basic server implementation capable of providing
custom URLs to be used by the standard go tools.

```bash
# To install latest version
curl -sfl https://raw.githubusercontent.com/bryk-io/go-vanity/master/install.sh | sh -s -- -b /usr/local/bin
```

## Background

When using third-party packages in Go, they are imported by a path that represents
how to download that package source. For example, to use the popular structured 
logging library, Logrus, it would be imported at the top of the Go program like so:

```go
import (
  "github.com/sirupsen/logrus"
)
```

In practice, import paths are as simple as the URL to the source code repository
that holds the package’s code. When `go get` is then executed, it fetches the Logrus
source code from GitHub and places the code in the $GOPATH/src directory.

```
$ tree $GOPATH/src
...
├── github.com
│   ├── sirupsen
│   │   └── logrus
...
```

The upside of the approach is that there is no need for a centrally managed package
server (like RubyGems, NPM, etc) to host distributable packages distinct from
source code. One limitation is the direct coupling between import statements
that depend on a certain package, and the location where that package’s author has
decided to host the code. Introducing a potential breaking change for users of the
package if the repository is ever moved to a different location.

## Remote Import Paths

The solution introduced by the Go team is the use of [Remote Import Paths](https://golang.org/cmd/go/#hdr-Remote_import_paths). Basically, any URL used as an 
import path could potentially return custom metadata specifying the location and
mechanisms to acquire the source code for the package.

The metadata tag must be of the form:

```html
<meta name="go-import" content="import-prefix vcs repo-root">
```

## Usage

A server instance uses a simple YAML or JSON configuration file to register paths
with their corresponding source code management system. For example, using the
following configuration file:

```yaml
host: custom.company.com
cache_max_age: 3600
paths:
  sample:
    repo: https://github.com/company/sample
    vcs: git
  another:
    repo: https://bitbucket.org/company/another
    vcs: git
```

- A user of the package `sample` could get the package's source code from Github
  by running `go get custom.company.com/sample`.
- A user of the package `another` could get the package's source code from Bitbucket
  by running `go get custom.company.com/another`.
- The import path are no longer tied to the location where the code is hosted.
- The package's author can move the code to a different location without introducing
  breaking changes in the projects using the package.