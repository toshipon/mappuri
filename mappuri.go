package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/gorilla/schema"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"net/url"
)

const mongoDbUrl = "localhost"

var decoder = schema.NewDecoder()

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

func loadOutings() (*[]Outing, error) {
	session, err := mgo.Dial(mongoDbUrl)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	c := session.DB("mappuri").C("outings")
	outings := []Outing{}
	err = c.Find(nil).All(&outings)
	if err != nil {
		return nil, err
	}
	return &outings, nil
}

func GetOutingHandler(w http.ResponseWriter, r *http.Request) {
	outingId := r.URL.Query().Get(":outingId")
	outing, err := loadOuting(outingId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(outing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(j))
}

func GetOutingsHandler(w http.ResponseWriter, r *http.Request) {
	outings, err := loadOutings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(outings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(j))
}

func CreateOutingHandler(w http.ResponseWriter, r *http.Request) {
	outing := new(Outing)
	r.ParseForm()
	err := decoder.Decode(outing, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(outing)
}

func main() {
	r := pat.New()
	r.Get("/outings/{outingId}", http.HandlerFunc(GetOutingHandler))
	r.Get("/outings", http.HandlerFunc(GetOutingsHandler))
	r.Post("/outings", http.HandlerFunc(CreateOutingHandler))
	http.Handle("/", r)
	fmt.Println("Running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
