package main

import (
	"encoding/json"
	"log"
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

	ls := LocationAreasResponse{}

	body := FetchWithCache(cfg, url)
	err := json.Unmarshal(body, &ls)
	if err != nil {
		log.Fatal(err)
	}

	cfg.Next = ls.Next
	if !fetchPrevious {
		cfg.Prev = &url
	} else {
		cfg.Prev = ls.Previous
	}
	return ls.Results, nil
}
