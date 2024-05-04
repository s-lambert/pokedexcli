package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func FetchWithCache(
	cfg *cliConfig,
	url string,
) []byte {
	cached, ok := cfg.cache.Get(url)
	if ok {
		return cached
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
		cfg.cache.Add(url, body)
		return body
	}
}
