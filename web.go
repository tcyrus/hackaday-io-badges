package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type BadgeData struct {
	Skulls int
	Name string
}

var HACKADAY_IO_API_KEY = os.Getenv("HACKADAY_IO_API_KEY")

var Badge, _ = template.ParseFiles("views/badge.svg")

func RedirectHandler(path string) http.Handler {
	return http.RedirectHandler(path, http.StatusMovedPermanently)
}

func FileHandler(str string) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func BadgeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data, err := getProject(strings.Replace(ps.ByName("id"), ".svg", "", 1))

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
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/", FileHandler("views/index.html"))
	router.GET("/:id", BadgeHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Fatal(http.ListenAndServe(":" + port, router))
}
