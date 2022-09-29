package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	api "github.com/google/cel-go/repl/appengine"
)

const indexTmplSrc = `
<!DOCTYPE html>
<meta charset="utf-8">
<title>CEL REPL</title>

<h1>REPL</h1>
<p>Add statements to evaluate.</p>
<p>Do not include sensitive data.</p>
<p>See <a href="https://github.com/google/cel-spec">github.com/google/cel-spec</a> for CEL language
overview. See
<a href="https://github.com/google/cel-go/tree/master/repl/main">github.com/google/cel-go/tree/master/repl/main</a>
for REPL syntax guide.</p>

<div class="input-block">
  <span>loading...</span><!--insertion point-->
</div>
<div class="controls">
  <button id="add-statement">Add Statement</button>
  <button id="evaluate">Evaluate</button>
</div>
<div class="output">
  <pre id="error-box"></pre>
</div>
<div class="output">
  <pre><code id="result">
  </code></pre>
</div>

<script src="/static/app.js"></script>
`

var indexTmpl *template.Template
var port string
var serveStatic string

func init() {
	indexTmpl = template.Must(template.New("index.html").Parse(indexTmplSrc))
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	flag.StringVar(&serveStatic, "serve_static", "", "serve static files from binary using argument as root")
	flag.Parse()
}

type context struct{}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("unhandled request: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	err := indexTmpl.Execute(w, &context{})
	if err != nil {
		log.Printf("tmpl error -- %s", err)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api", api.NewJsonHandler())
	if serveStatic != "" {
		log.Printf("serving static '%s'", serveStatic)
		http.Handle("/static/",
			http.StripPrefix("/static",
				http.FileServer(http.Dir(serveStatic))))
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
