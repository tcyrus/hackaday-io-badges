package main

import (
	"encoding/json"
	"errors"
	"github.com/tcyrus/hackaday-io-badges/Godeps/_workspace/src/github.com/flosch/pongo2"
	"github.com/tcyrus/hackaday-io-badges/Godeps/_workspace/src/github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
)

const HACKADAY_IO_API_KEY string = os.Getenv("HACKADAY_IO_API_KEY")

const Badge *pongo2.Template = pongo2.Must(pongo2.FromFile("views/badge.svg"))

func RedirectHandler(path string) http.Handler {
	return http.RedirectHandler(path, http.StatusMovedPermanently)
}

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

	if _, ok := data["project"]; ok {
		return nil, errors.New("Invalid Project ID")
	}

	if val, ok := data["message"]; ok {
		return nil, errors.New(val.(string))
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
	name := data["name"]

	w.Header().Set("Content-Type", "image/svg+xml")

	// Execute the template per HTTP request
	if err := Badge.ExecuteWriter(pongo2.Context{"skulls": skulls, "name": name}, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.Handle("/", RedirectHandler("/hackaday"))
	r.Handle("/hackaday/", RedirectHandler("/hackaday"))
	r.HandleFunc("/hackaday", FileHandler("views/index.html"))
	r.HandleFunc("/hackaday/{id}.svg", BadgeHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	http.ListenAndServe(":" + port, r)
}
