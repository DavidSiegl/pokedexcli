package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2/location-area/"

type locationAreaResponse struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func fetchLocationAreas(url string, cache *pokecache.Cache) (locationAreaResponse, error) {
	if data, ok := cache.Get(url); ok {
		fmt.Println("(cache hit)")
		var result locationAreaResponse
		if err := json.Unmarshal(data, &result); err != nil {
			return locationAreaResponse{}, err
		}
		return result, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return locationAreaResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return locationAreaResponse{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result locationAreaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return locationAreaResponse{}, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return locationAreaResponse{}, err
	}
	cache.Add(url, data)

	return result, nil
}
