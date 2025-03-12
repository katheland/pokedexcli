package pokeapi

import (
	"net/http"
	"io"
	"encoding/json"
	"fmt"
	"internal/pokecache"
	"time"
)

// structure of the encounters data
type Encounters struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// initialize a cache
var exploreCache pokecache.Cache
func init() {
	exploreCache = pokecache.NewCache(5 * time.Second)
}

// lists the Pokemon found at a given location
func ExploreFunction(url string) error {
	fullUrl := "https://pokeapi.co/api/v2/location-area/" + url

	jsonData, ok := exploreCache.Get(fullUrl)
	if !ok { // it's not in the cache so we're calling the API
		res, err := http.Get(fullUrl)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		jsonData, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		mapCache.Add(url, jsonData)
	}
	
	var encounters Encounters
	if err := json.Unmarshal(jsonData, &encounters); err != nil {
		return err
	}

	for _, mon := range encounters.PokemonEncounters {
		fmt.Println(mon.Pokemon.Name)
	}
	
	return nil
}