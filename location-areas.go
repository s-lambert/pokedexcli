package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreasResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func GetLocationAreas(cfg *cliConfig, fetchPrevious bool) ([]LocationArea, error) {
	var url string
	if fetchPrevious && cfg.Prev != nil {
		url = *cfg.Prev
	} else {
		url = cfg.Next
	}
	fmt.Printf("Fetching: %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	ls := LocationAreasResponse{}
	err = json.Unmarshal(body, &ls)
	if err != nil {
		log.Fatal(err)
	}
	cfg.Next = ls.Next
	cfg.Prev = ls.Previous
	return ls.Results, nil
}
