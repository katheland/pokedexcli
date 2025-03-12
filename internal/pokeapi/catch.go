package pokeapi

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"internal/pokecache"
	"time"
	"math/rand"
)

// structure of the encounter rate data
type CaptureRate struct {
	CaptureRate   int `json:"capture_rate"`
	Name                 string `json:"name"`
}

// structure of the pokemon species data
// (why is this separated, and why isn't catch rate here...)
type Pokemon struct {
	Height    int `json:"height"`
	ID                     int    `json:"id"`
	Name          string `json:"name"`
	Stats []struct {
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
	Weight int `json:"weight"`
}

// initialize a cache
var catchCache pokecache.Cache
func init() {
	catchCache = pokecache.NewCache(60 * time.Second)
}

// gets the capture rate of a given species and attempts to catch it
func CatchFunction(species string) (string, bool, error) {
	fullUrl := "https://pokeapi.co/api/v2/pokemon-species/" + species

	jsonData, ok := catchCache.Get(fullUrl)
	if !ok { // it's not in the cache so we're calling the API
		res, err := http.Get(fullUrl)
		if err != nil {
			return "", false, err
		}
		defer res.Body.Close()

		jsonData, err = io.ReadAll(res.Body)
		if err != nil {
			return "", false, err
		}

		catchCache.Add(fullUrl, jsonData)
	}

	// the guided project says to use the pokemon's base experience
	// but I know how pokemon works and that's not it
	// I'm not bothering with the entire equation for now
	// but I'm using the actual catch rate if nothing else >_<
	var catchRate CaptureRate
	if err := json.Unmarshal(jsonData, &catchRate); err != nil {
		return "", false, err
	}

	fmt.Println("Throwing a Pokeball at " + catchRate.Name + "...")
	attempt := rand.Intn(256)
	caught := attempt <= catchRate.CaptureRate
	
	return catchRate.Name, caught, nil
}

// get species data when we're registering a new one to the pokedex
func GetSpeciesData(species string) (Pokemon, error) {
	fullUrl := "https://pokeapi.co/api/v2/pokemon/" + species

	// not caching this because the pokedex itself will act as a cache
	res, err := http.Get(fullUrl)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	var pokemon Pokemon
	if err := json.Unmarshal(jsonData, &pokemon); err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}