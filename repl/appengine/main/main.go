package main

import (
	"flag"
	"net/http"
	"text/template"

	"log"
	"os"

	app "github.com/google/cel-go/repl/appengine"
)

var indexTmpl *template.Template
var port string
var serveStatic string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	flag.StringVar(&serveStatic, "serve_static", "", "serve static files from binary using argument as root")
	flag.Parse()
}

type context struct{}

func main() {
	mux := app.NewAppMux(serveStatic)
	log.Printf("Listening on port %s", port)
	http.Handle("/", mux)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
