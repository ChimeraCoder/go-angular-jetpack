package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

const PROFILE_SESSION = "profile"

var (
	httpAddr = flag.String("addr", ":8000", "HTTP server address")
)

// renderAngularTemplate sets the delimiters for the specificed template(s) to be "[[" and "]]"
// and then parses and renders all templates specified by the filenames
func renderAngularTemplate(w http.ResponseWriter, data interface{}, filenames ...string) {
	t := template.New("base")
	t.Delims("[[", "]]")
	s1, err := t.ParseFiles(append(filenames, "templates/angularBase.tmpl")...)
	if err != nil {
		panic(err)
	}
	err = s1.ExecuteTemplate(w, "base", nil)
	if err != nil {
		panic(err)
	}
}

func serveAngularHome(w http.ResponseWriter, r *http.Request) {
	renderAngularTemplate(w, nil, "templates/index.tmpl")
}

func servePhones(w http.ResponseWriter, r *http.Request) {
	renderAngularTemplate(w, nil, "templates/phones.tmpl")
}

// serveJSON serves the JSON representation of arbitrary data
// Useful for serving api.example.com/users/1
func serveJson(w http.ResponseWriter, r *http.Request, data interface{}) {
	bts, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, string(bts))
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", serveAngularHome)
	r.HandleFunc("/phones", servePhones)
	http.Handle("/static/", http.FileServer(http.Dir("public")))
	http.Handle("/", r)

	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatalf("Error listening, %v", err)
	}
}
