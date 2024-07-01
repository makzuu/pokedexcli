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

const baseUrl = "https://pokeapi.co/api/v2/"

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
	url := baseUrl + "location-area/" + area
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

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Stats          []struct {
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

func PokemonInfo(name string) (pokemon Pokemon, err error) {
	url := baseUrl + "pokemon/" + name
	body, err := request(url)
	if err != nil {
		return
	}

	pokemon = Pokemon{}

	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return
	}

	return
}

func request(url string) (body []byte, err error) {
	body, ok := cache.Get(url)
	if ok {
		return
	}
	res, err := http.Get(url)
	if err != nil {
		return
	}
	body, err = io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return
	}
	if res.StatusCode > 299 {
		err = fmt.Errorf("response fail with status code: %d\n", res.StatusCode)
		return
	}
	cache.Add(url, body)
	return
}
