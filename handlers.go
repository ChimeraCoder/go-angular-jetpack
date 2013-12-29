package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
	"time"
)

// serveAngularHome is an example of a handler function for serving a template
func serveAngularHome(w http.ResponseWriter, r *http.Request) {
	renderAngularTemplate(w, nil, "templates/index.tmpl")
}

func servePhones(w http.ResponseWriter, r *http.Request) {
	renderAngularTemplate(w, nil, "templates/phones.tmpl")
}

// serveUserJson is an example of a hander function that returns a JSON response
// instead of a response with Content-Type: text/html
func serveUserJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		// This shouldn't really be a runtime panic, but it domonstrates how
		// the HTTP response to the client will handled by handleError()
		panic(fmt.Errorf("Username not provided"))
	}

	// Create a struct that should be returned
	// This could be a database lookup
	user := struct {
		Username  string
		CreatedAt int64
	}{username, time.Now().Unix()}

	serveJson(w, r, user)
}

func serveUser(w http.ResponseWriter, r *http.Request) {
	// TODO fill this out
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
	var b *bytes.Buffer
	err = s1.ExecuteTemplate(b, "base", nil)
	if err != nil {
		panic(err)
	}
	w.Write(b.Bytes())
}

// handleError specifies the behavior when a handler function (controller)
// encounters a runtime panic
func handleError(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("ERROR: Recovern from panic: %v", r)
				http.Error(w, "An unexpected server error has occurred", http.StatusInternalServerError)
			}
		}()
		f(w, r)
	}
}

// serveJson serves the JSON representation of arbitrary data
// Useful for serving api.example.com/users/1
func serveJson(w http.ResponseWriter, r *http.Request, data interface{}) {
	bts, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(bts)
}
