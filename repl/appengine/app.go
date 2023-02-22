package appengine

import (
	"log"
	"net/http"
)

func NewAppMux(staticDir string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler("/ng/", 301))
	mux.HandleFunc("/api", NewJsonHandler())
	if staticDir != "" {
		log.Printf("serving static from '%s'", staticDir)
		mux.Handle("/ng/",
			http.StripPrefix("/ng",
				http.FileServer(http.Dir(staticDir))))
	}
	return mux
}
