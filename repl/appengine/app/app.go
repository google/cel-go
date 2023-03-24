// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
