package main

import (
	"fmt"
	"github.com/gorilla/pat"
	"github.com/gorilla/schema"
	"html/template"
	"io/ioutil"
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

func showOuting(w http.ResponseWriter, r *http.Request) {
	outingId := r.URL.Query().Get(":outingId")
	outing, err := loadOuting(outingId)
	if err != nil {
		log.Fatal(err)
	}
	t, _ := template.ParseFiles("templates/outing_detail.html")
	t.Execute(w, outing)
}

func listOutings(w http.ResponseWriter, r *http.Request) {
	outing, err := loadOutings()
	if err != nil {
		log.Fatal(err)
	}
	t, _ := template.ParseFiles("templates/outing_detail.html")
	t.Execute(w, outing)
}

func createOuting(w http.ResponseWriter, r *http.Request) {
    outing := Outing{}
	r.ParseForm()
	decoder.Decode(outing, r.PostForm)
	fmt.Println(outing)
}

func main() {
	r := pat.New()
	r.Get("/", HomeHandler)
	r.Get("/outings/{outingId}", http.HandlerFunc(showOuting))
	r.Get("/outings", http.HandlerFunc(listOutings))
	r.Post("/outings", http.HandlerFunc(createOuting))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
