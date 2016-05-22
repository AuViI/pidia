package diaweb

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
)

var (
	// GoPath environtment variable $GOPATH
	GoPath     = os.Getenv("GOPATH")
	tmplFile   = fmt.Sprintf("%s/%s", GoPath, "src/github.com/auvii/pidia/diaweb/template.html")
	tmpl, terr = template.ParseFiles(tmplFile)
)

// Execute executes template for server
func (s *Server) Execute(w http.ResponseWriter, r *http.Request) {
	if throw(w, terr) {
		return
	}
	err := tmpl.Execute(w, s)
	throw(w, err)
	s.SetRefresh(w)
}

func throw(w http.ResponseWriter, e error) (r bool) {
	r = e != nil
	if r {
		fmt.Fprintf(w, "Error: %s", e)
		fmt.Println("Error:", e)
	}
	return
}
