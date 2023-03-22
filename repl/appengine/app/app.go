// Package app provides a simple JSON API for using the REPL in a web application.
package app

import (
	"log"
	"net/http"
)

// NewAppMux returns an http.ServeMux that handles requests for the REPL app.
//
// Optionally serves the static website content from directory staticDir (if
// non-empty).
func NewAppMux(staticDir string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler("/ng/", 301))
	mux.HandleFunc("/api", NewJSONHandler())
	if staticDir != "" {
		log.Printf("serving static from '%s'", staticDir)
		mux.Handle("/ng/",
			http.StripPrefix("/ng",
				http.FileServer(http.Dir(staticDir))))
	}
	return mux
}
