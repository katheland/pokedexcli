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

// initialize a cache
var catchCache pokecache.Cache
func init() {
	catchCache = pokecache.NewCache(60 * time.Second)
}

// gets the capture rate of a given species and attempts to catch it
func CatchFunction(species string) (string, bool, error) {
	fullUrl := "https://pokeapi.co/api/v2/pokemon-species/" + species

	jsonData, ok := exploreCache.Get(fullUrl)
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

		mapCache.Add(fullUrl, jsonData)
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