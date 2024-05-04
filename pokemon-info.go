package main

import (
	"encoding/json"
	"fmt"
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

func PrintPokemonInfo(p PokemonInfo) {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Println("Stats")
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
}
