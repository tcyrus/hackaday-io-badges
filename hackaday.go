package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

var HACKADAY_IO_API_KEY = os.Getenv("HACKADAY_IO_API_KEY")

func getProject(id string) (data map[string]interface{}, err error) {
	r2, err := http.Get("https://api.hackaday.io/v1/projects/" + id + "?api_key=" + HACKADAY_IO_API_KEY)
	if err != nil {
		return nil, err
	}
	if r2 == nil {
		return nil, errors.New("Empty API Response")
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
