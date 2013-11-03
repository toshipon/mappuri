package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/gorilla/schema"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"net/url"
	"html/template"
)

const mongoDbUrl = "localhost"

var (
	decoder = schema.NewDecoder()
	port    = flag.String("p", "8000", "Port number (default 8000)")
)

type Place struct {
	Name           string
	MapLink        *url.URL
	FoursquareLink *url.URL
	Website        *url.URL
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
	sess, err := mgo.Dial(mongoDbUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer sess.Close()
	outing := new(Outing)
	r.ParseForm()
	err = decoder.Decode(outing, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	outing.Id = bson.NewObjectId()
	collection := sess.DB("mappuri").C("outings")
	err = collection.Insert(outing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CreatePlaceHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := mgo.Dial(mongoDbUrl)
	place := new(Place)
	r.ParseForm()
	outingId := r.FormValue("outingId")
	outing, err := loadOuting(outingId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Name := r.FormValue("Name")
	MapLink, err := url.Parse(r.FormValue("MapLink"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	FoursquareLink, err := url.Parse(r.FormValue("FoursquareLink"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	place.Name = Name
	place.MapLink = MapLink
	place.FoursquareLink = FoursquareLink
	collection := sess.DB("mappuri").C("outings")
	outing.Places = append(outing.Places, *place)
	err = collection.UpdateId(bson.ObjectIdHex(outingId), &outing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Create and parse Template
	t, err := template.New("index.html").ParseFiles("frontend/public/index.html")
	if err != nil {
		log.Panic(err)
	}
	// Render the template
	err = t.Execute(w, map[string]interface{}{})
	if err != nil {
		log.Panic(err)
	}
}

type ResponseWriterProxy struct {
	realResponseWriter	*http.ResponseWriter
	code			int
}

func (rwp *ResponseWriterProxy) Header() http.Header {
	return (*rwp.realResponseWriter).Header()
}

func (rwp *ResponseWriterProxy) Write(buf []byte) (int, error) {
	return (*rwp.realResponseWriter).Write(buf)
}

func (rwp *ResponseWriterProxy) WriteHeader(code int) {
	(*rwp.realResponseWriter).WriteHeader(code)
	rwp.code = code
}

func maxAgeHandler(seconds int, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseWriterProxy := &ResponseWriterProxy{realResponseWriter: &w}
		h.ServeHTTP(responseWriterProxy, r)
		if responseWriterProxy.code == 200 {	//only cache the response if it's successful
			w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))
		}

	})
}

func main() {
	r := pat.New()
	r.Get("/outings/{outingId}", http.HandlerFunc(GetOutingHandler))
	r.Get("/outings", http.HandlerFunc(GetOutingsHandler))
	r.Post("/outings", http.HandlerFunc(CreateOutingHandler))
	r.Post("/places", http.HandlerFunc(CreatePlaceHandler))
	r.HandleFunc("/", HomeHandler)
	http.Handle("/assets/", maxAgeHandler(10*365*24*60*60, http.StripPrefix("/assets/", http.FileServer(http.Dir("frontend/public/")))))
	http.Handle("/", r)
	fmt.Println("Running on localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
