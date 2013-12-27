package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
	"text/template"
)

const PROFILE_SESSION = "profile"

var (
	httpAddr        = flag.String("addr", ":8000", "HTTP server address")
	baseTmpl string = "templates/base.tmpl"
	store           = sessions.NewCookieStore([]byte(COOKIE_SECRET)) //CookieStore uses secure cookies
	decoder         = schema.NewDecoder()                            //From github.com/gorilla/schema

	//The following three variables can be defined using environment variables
	//to avoid committing them by mistake

	COOKIE_SECRET = []byte(os.Getenv("COOKIE_SECRET"))
)

func renderTemplate(w http.ResponseWriter, name string, data interface{}, filenames ...string) {
	s1, _ := template.ParseFiles(filenames...)
	s1.ExecuteTemplate(w, name, data)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	serveAngularHome(w, r)
}

func serveBase(w http.ResponseWriter, r *http.Request) {
	s1, err := template.ParseFiles("templates/base.tmpl", "templates/index.tmpl")
	if err != nil {
		panic(err)
	}
	s1.ExecuteTemplate(w, "base", nil)
}

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
	r.HandleFunc("/profile", serveHome).Methods("GET")
	http.Handle("/static/", http.FileServer(http.Dir("public")))
	http.Handle("/", r)

	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatalf("Error listening, %v", err)
	}
}
