package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/makzuu/pokedexcli/internal/pokecache"
	"io"
	"net/http"
	"time"
)

var cache = pokecache.NewCache(5 * time.Minute)

type locationAreas struct {
	Next      string `json:"next"`
	Previous  string `json:"previous"`
	Locations []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func GetLocations(url string) (locationAreas, error) {
	var locations locationAreas
	val, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(val, &locations)
		return locations, err
	}
	response, err := http.Get(url)
	if err != nil {
		return locations, err
	}
	if response.StatusCode > 299 {
		return locations, errors.New(fmt.Sprintf("Response fail with status code: %d\n", response.StatusCode))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return locations, err
	}
	response.Body.Close()

	cache.Add(url, body)

	err = json.Unmarshal(body, &locations)
	return locations, err
}
