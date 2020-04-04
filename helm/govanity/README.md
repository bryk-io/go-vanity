# Go Vanity

This package allows to provide custom import paths for go dependencies hosted
on various source code repositories by specifying a simple configuration file.

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