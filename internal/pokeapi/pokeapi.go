package pokeapi

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"internal/pokecache"
	"time"
)

// the structure of the location data
type Location struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// initialize a cache
var mapCache pokecache.Cache
func init() {
	mapCache = pokecache.NewCache(5 * time.Second)
}

// print a page of 20 locations
func MapFunction(url *string) (string, *string, error) {
	jsonData, ok := mapCache.Get(*url)
	if !ok { // it's not in the cache so we're calling the API
		res, err := http.Get(*url)
		if err != nil {
			return "", nil, err
		}
		defer res.Body.Close()

		jsonData, err = io.ReadAll(res.Body)
		if err != nil {
			return "", nil, err
		}

		mapCache.Add(*url, jsonData)
	}
	
	var locations Location
	if err := json.Unmarshal(jsonData, &locations); err != nil {
		return "", nil, err
	}

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	
	return locations.Next, locations.Previous, nil
}