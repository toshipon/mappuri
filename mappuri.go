package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"net/url"
)

const mongoDbUrl = "localhost"

type Place struct {
	Id             bson.ObjectId `bson:"_id"`
	Name           string
	MapLink        *url.URL `bson:",omitempty"`
	FoursquareLink *url.URL `bson:",omitempty"`
	Website        *url.URL `bson:",omitempty"`
}

type Outing struct {
	Id     bson.ObjectId `bson:"_id"`
	Name   string
	Places []Place
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprint(w, string(b))
	return
}

func loadOuting(outingId string) (*Outing, error) {
	session, err := mgo.Dial(mongoDbUrl)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	c := session.DB("mappuri").C("outings")
	outing := Outing{}
	err = c.FindId(bson.ObjectIdHex(outingId)).One(&outing)
	if err != nil {
		return nil, err
	}
	return &outing, nil
}

func OutingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	outingId := vars["outingId"]
	outing, err := loadOuting(outingId)
	if err != nil {
		log.Fatal(err)
	}
	t, _ := template.ParseFiles("templates/outing_detail.html")
	t.Execute(w, outing)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/outings/{outingId}", OutingHandler)
	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":8080", nil)
}
