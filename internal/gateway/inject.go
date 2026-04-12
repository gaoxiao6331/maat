package gateway

import "strings"

func InjectEnv(html, env string) string {
	inject := "<script>window.__MAAT_ENV__='" + env + "'</script>"
	lower := strings.ToLower(html)
	idx := strings.LastIndex(lower, "</body>")
	if idx < 0 {
		return html + inject
	}
	return html[:idx] + inject + html[idx:]
}
