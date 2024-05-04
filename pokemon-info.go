package main

import (
	"encoding/json"
	"log"
)

type PokemonInfo struct {
	Height int    `json:"height"`
	Name   string `json:"name"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Weight         int `json:"weight"`
	BaseExperience int `json:"base_experience"`
}

func GetPokemonInfo(cfg *cliConfig, pokemonName string) (PokemonInfo, error) {
	url := cfg.PokemonInfo + pokemonName
	p := PokemonInfo{}

	body := FetchWithCache(cfg, url)
	err := json.Unmarshal(body, &p)
	if err != nil {
		log.Fatal(err)
	}

	return p, nil
}
