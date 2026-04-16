package handlers

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(staticFS fs.FS) http.Handler {

	mux := chi.NewMux()
	mux.Use(middleware.CleanPath)

	apiRouter := chi.NewMux()
	apiRouter.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "ok\n")
	})

	mux.Mount("/api", apiRouter)

	mux.Handle("GET /*", makePwaHandler(staticFS, "webapp/dist/webapp/browser", "index.html"))

	return mux
}

func makePwaHandler(staticFS fs.FS, rootDir, indexFile string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appRoot, err := fs.Sub(staticFS, rootDir)
		if err != nil {
			http.Error(w, fmt.Sprintf("error in fs.Sub(): %v", err), http.StatusInternalServerError)
			return
		}

		filePath := strings.TrimPrefix(r.URL.Path, "/")

		if _, err = fs.Stat(appRoot, filePath); err != nil {
			http.ServeFileFS(w, r, appRoot, indexFile)
			return
		}
		http.FileServerFS(appRoot).ServeHTTP(w, r)
	})
}
