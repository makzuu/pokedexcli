package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type locationAreas struct {
	Next      string `json:"next"`
	Previous  string `json:"previous"`
	Locations []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func GetLocations(url string) (locationAreas, error) {
	response, err := http.Get(url)
	var locations locationAreas
	if err != nil {
		return locations, err
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return locations, err
	}
	if response.StatusCode > 299 {
		return locations, errors.New("Response fail with status code: " + string(response.StatusCode))
	}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return locations, err
	}

	return locations, err
}
