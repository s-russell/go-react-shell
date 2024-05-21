package web

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"path"
)

type FrontendConfig struct {
	Title string `json:"title"`
}

var (
	//go:embed static
	staticFS embed.FS
	config   FrontendConfig
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

	config = FrontendConfig{
		"[title goes here]",
	}

	mux.HandleFunc("/config.json", func(w http.ResponseWriter, r *http.Request) {
		configJSON, err := json.Marshal(config)
		if err != nil {
			log.Default().Printf("Failed to marshal config for frontend: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(configJSON)
	})
}
