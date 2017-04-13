package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type BadgeData struct {
	Skulls int
	Name string
}

var HACKADAY_IO_API_KEY = os.Getenv("HACKADAY_IO_API_KEY")

var Badge, _ = template.ParseFiles("views/badge.svg")

func FileHandler(str string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, str)
	}
}

func getProject(id string) (data map[string]interface{}, err error) {
	r2, err := http.Get("https://api.hackaday.io/v1/projects/" + id + "?api_key=" + HACKADAY_IO_API_KEY)
	if err != nil {
		return nil, err
	}

	defer r2.Body.Close()

	body, err := ioutil.ReadAll(r2.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if _, invalid := data["project"]; invalid {
		return nil, errors.New("Invalid Project ID")
	}

	if message, ok := data["message"]; ok {
		return nil, errors.New(message.(string))
	}

	return data, nil
}

func BadgeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getProject(mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	skulls := int(data["skulls"].(float64))
	name := data["name"].(string)

	w.Header().Set("Content-Type", "image/svg+xml")

	// Execute the template per HTTP request
	Badge.Execute(w, &BadgeData{Skulls: skulls, Name: name})
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", FileHandler("views/index.html"))
	r.HandleFunc("/{id:[0-9]+}.svg", BadgeHandler)
	r.HandleFunc("/{id:[0-9]+}", BadgeHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Fatal(http.ListenAndServe(":" + port, r))
}
