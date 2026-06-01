package template

import "html/template"

var Page = template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>{{ .Title }}</title>
	<meta
		name="description"
		content="{{ .Description }}"
	>
</head>
<body>
	{{ .Content }}
</body>
</html>
`),
)
