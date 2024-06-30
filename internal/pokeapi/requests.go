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

const baseUrl = "https://pokeapi.co/api/v2/location-area/"

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
		return locations, errors.New(fmt.Sprintf("Response fail with status code: %d", response.StatusCode))
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

func GetPokemons(area string) ([]string, error) {
	url := baseUrl + area
	var body []byte

	body, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if res.StatusCode > 299 {
			if res.StatusCode == 404 {
				return nil, fmt.Errorf("location: %s not found", area)
			}
			return nil, fmt.Errorf("Response fail with status code: %d", res.StatusCode)
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, err
		}

		cache.Add(url, body)
	}

	data := struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string
			}
		} `json:"pokemon_encounters"`
	}{}

	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	pokemons := make([]string, len(data.PokemonEncounters))
	for i := 0; i < len(data.PokemonEncounters); i++ {
		pokemons[i] = data.PokemonEncounters[i].Pokemon.Name
	}
	return pokemons, nil
}
