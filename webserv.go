package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func FileHandler(str string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, str)
	}
}

func BadgeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getProject(mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmp_data := &BadgeData{Skulls: int(data["skulls"].(float64)), Name: data["name"].(string)}

	w.Header().Set("Content-Type", "image/svg+xml")

	// Execute the template per HTTP request
	Badge.Execute(w, tmp_data)
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
