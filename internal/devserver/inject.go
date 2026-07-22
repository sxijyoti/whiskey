package devserver

import "strings"

func InjectReload(
	html string,
) string {

	script := `
<script>
const es = new EventSource("/__reload");

es.onmessage = () => {
	location.reload();
};
</script>
`

	return strings.Replace(
		html,
		"</body>",
		script+"</body>",
		1,
	)
}
