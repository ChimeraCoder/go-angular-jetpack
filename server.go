package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const PROFILE_SESSION = "profile"

var (
	httpAddr = flag.String("addr", ":8000", "HTTP server address")
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleError(serveAngularHome))
	r.HandleFunc("/users/{username:[A-Za-z0-9]+}.json", handleError(serveUserJson))
	r.HandleFunc("/users/{username:[A-Za-z0-9]+}.json", handleError(serveUser))
	r.HandleFunc("/phones", handleError(servePhones))
	http.Handle("/static/", http.FileServer(http.Dir("public")))
	http.Handle("/", r)

	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatalf("Error listening, %v", err)
	}
}
