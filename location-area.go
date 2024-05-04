package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LocationAreaInfo struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"encounter_method,omitempty"`
		VersionDetails []struct {
			Rate    int `json:"rate,omitempty"`
			Version struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"version,omitempty"`
		} `json:"version_details,omitempty"`
	} `json:"encounter_method_rates,omitempty"`
	GameIndex int `json:"game_index,omitempty"`
	ID        int `json:"id,omitempty"`
	Location  struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"location,omitempty"`
	Name  string `json:"name,omitempty"`
	Names []struct {
		Language struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"language,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"names,omitempty"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"pokemon,omitempty"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance,omitempty"`
				ConditionValues []any `json:"condition_values,omitempty"`
				MaxLevel        int   `json:"max_level,omitempty"`
				Method          struct {
					Name string `json:"name,omitempty"`
					URL  string `json:"url,omitempty"`
				} `json:"method,omitempty"`
				MinLevel int `json:"min_level,omitempty"`
			} `json:"encounter_details,omitempty"`
			MaxChance int `json:"max_chance,omitempty"`
			Version   struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"version,omitempty"`
		} `json:"version_details,omitempty"`
	} `json:"pokemon_encounters,omitempty"`
}

func GetLocationArea(cfg *cliConfig, locationName string) (LocationAreaInfo, error) {
	url := cfg.AreaInfo + locationName
	la := LocationAreaInfo{}
	cached, ok := cfg.cache.Get(url)
	if ok {
		fmt.Printf("Getting from cache: %s\n", url)
		err := json.Unmarshal(cached, &la)
		if err != nil {
			log.Fatal(err)
		}
	} else {
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
		err = json.Unmarshal(body, &la)
		if err != nil {
			log.Fatal(err)
		}
		cfg.cache.Add(url, body)

	}

	return la, nil
}
