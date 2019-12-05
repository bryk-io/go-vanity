package main

import (
	"html/template"
)

var indexTpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
<h1>{{.Host}}</h1>
<ul>
{{range .Repos}}<li><a href="{{.Source}}">{{.Import}} ({{.VCS}})</a></li>{{end}}
</ul>
</html>
`))

var repoTpl = template.Must(template.New("repo").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.Import}} {{.VCS}} {{.Source}}">
<meta name="go-source" content="{{.Import}} {{.Display}}">
</head>
<body>Nothing to see here.</body>
</html>`))