package appengine

import (
	"html/template"
	"log"
	"net/http"
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

type templContext struct{}

func init() {
	indexTmpl = template.Must(template.New("index.html").Parse(indexTmplSrc))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("unhandled request: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	err := indexTmpl.Execute(w, &templContext{})
	if err != nil {
		log.Printf("tmpl error -- %s", err)
	}
}

func NewAppMux(staticDir string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	http.HandleFunc("/api", NewJsonHandler())
	if staticDir != "" {
		log.Printf("serving static from '%s'", staticDir)
		mux.Handle("/static/",
			http.StripPrefix("/static",
				http.FileServer(http.Dir(staticDir))))
	}
	return mux
}
