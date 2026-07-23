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
			} else if err != nil && filepath.Ext(path) == "" {
				if htmlInfo, htmlErr := os.Stat(path + ".html"); htmlErr == nil && !htmlInfo.IsDir() {
					path = path + ".html"
				} else if idxInfo, idxErr := os.Stat(filepath.Join(path, "index.html")); idxErr == nil && !idxInfo.IsDir() {
					path = filepath.Join(path, "index.html")
				}
			}

			if filepath.Ext(path) != ".html" {
				if _, err := os.Stat(path); err == nil {
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
			}

			data, err := os.ReadFile(path)

			if err != nil {
				custom404Path := filepath.Join(distRoot, "404.html")
				if data404, err404 := os.ReadFile(custom404Path); err404 == nil {
					html := InjectReload(string(data404))
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.WriteHeader(http.StatusNotFound)
					_, _ = io.WriteString(w, html)
					return
				}

				custom404DirPath := filepath.Join(distRoot, "404", "index.html")
				if data404, err404 := os.ReadFile(custom404DirPath); err404 == nil {
					html := InjectReload(string(data404))
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.WriteHeader(http.StatusNotFound)
					_, _ = io.WriteString(w, html)
					return
				}

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
