package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Place struct {
	Name    string
	MapLink url.URL
}

type Outing struct {
	Places []Place
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprint(w, string(b))
	return
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
