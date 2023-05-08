// https://github.com/ignite/cli/blob/main/ignite/pkg/openapiconsole/console.go

package docs

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed index.tpl
var index embed.FS

// Handler returns a http handler that servers OpenAPI console for an OpenAPI spec at specURL.
func Handler(title, specURL string) http.HandlerFunc {
	t, _ := template.ParseFS(index, "index.tpl") // nolint:errcheck non-critical for docs

	return func(w http.ResponseWriter, req *http.Request) {
		_ = t.Execute(w, struct {
			Title string
			URL   string
		}{
			title,
			specURL,
		}) //nolint:errcheck non-critical for docs
	}
}
