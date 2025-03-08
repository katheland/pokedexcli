package pokeapi

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
)

type Location struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// print a page of 20 locations
func MapFunction(url *string) (string, *string, error) {
	res, err := http.Get(*url)
	if err != nil {
		return "", nil, err
	}
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil, err
	}

	var locations Location
	if err = json.Unmarshal(jsonData, &locations); err != nil {
		return "", nil, err
	}

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	
	return locations.Next, locations.Previous, nil
}