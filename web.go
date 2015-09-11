package main

import (
	"encoding/json"
	"github.com/tcyrus/hackaday-io-badges/Godeps/_workspace/src/github.com/flosch/pongo2"
	"github.com/tcyrus/hackaday-io-badges/Godeps/_workspace/src/github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Pre-compiling the templates at application startup using the
// little Must()-helper function (Must() will panic if FromFile()
// or FromString() will return with an error - that's it).
// It's faster to pre-compile it anywhere at startup and only
// execute the template later.

var tplBadge = pongo2.Must(pongo2.FromFile("views/badge.svg"))

func RedirectHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently)
	}
}

func FileHandler(str string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, str)
	}
}

func BadgeHandler(w http.ResponseWriter, r *http.Request) {
	var dat map[string]interface{}
	id := mux.Vars(r)["id"]

	// Execute the template per HTTP request
	req_get, _ := http.Get("https://api.hackaday.io/v1/projects/" + id + "?api_key=" + os.Getenv("HACKADAY_IO_API_KEY"))
	defer req_get.Body.Close()

	body, _ := ioutil.ReadAll(req_get.Body)

	if err := json.Unmarshal(body, &dat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := dat["project"]; ok {
		http.Error(w, "Invalid", http.StatusInternalServerError)
		return
	}

	if val, ok := dat["message"]; ok {
		http.Error(w, val.(string), http.StatusInternalServerError)
		return
	}

	skulls := int(dat["skulls"].(float64))
	name := dat["name"]

	w.Header().Set("Content-Type", "image/svg+xml")
	if err := tplBadge.ExecuteWriter(pongo2.Context{"skulls": skulls, "name": name}, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", RedirectHandler("/hackaday"))
	r.HandleFunc("/hackaday", FileHandler("views/index.html"))
	r.HandleFunc("/hackaday/favicon.ico", FileHandler("views/favicon.ico"))
	r.HandleFunc("/hackaday/{id}.svg", BadgeHandler)
	port := os.Getenv("PORT")
  if port == "" {
    log.Fatal("$PORT must be set")
  }
	http.ListenAndServe(":" + port, r)
}
