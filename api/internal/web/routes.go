package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
)

var (
	//go:embed static
	staticFS embed.FS
)

func AddRoutes(mux *http.ServeMux) {
	web, _ := fs.Sub(staticFS, "static")
	webServer := http.FileServerFS(web)
	//mux.Handle("/", webServer)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := http.FS(web).Open(path.Clean(r.URL.Path)); err != nil {
			// if asset isn't recognized, defer to client side router
			http.ServeFileFS(w, r, web, "index.html")
		} else {
			webServer.ServeHTTP(w, r)
		}
	})
}
