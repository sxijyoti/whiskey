package devserver

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func NewFileHandler(
	root string,
	dist string,
) http.Handler {

	distRoot := filepath.Join(
		root,
		dist,
	)

	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			path := filepath.Join(
				distRoot,
				r.URL.Path,
			)

			info, err := os.Stat(path)

			if err == nil && info.IsDir() {

				path = filepath.Join(
					path,
					"index.html",
				)
			}

			if filepath.Ext(path) != ".html" {

				if ctype := mime.TypeByExtension(
					filepath.Ext(path),
				); ctype != "" {

					w.Header().Set(
						"Content-Type",
						ctype,
					)
				}

				http.ServeFile(
					w,
					r,
					path,
				)

				return
			}

			data, err := os.ReadFile(path)

			if err != nil {

				http.NotFound(
					w,
					r,
				)

				return
			}

			html := InjectReload(
				string(data),
			)

			w.Header().Set(
				"Content-Type",
				"text/html; charset=utf-8",
			)

			_, _ = io.WriteString(
				w,
				html,
			)
		},
	)
}