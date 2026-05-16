package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2/location-area/"
const pokemonURL = "https://pokeapi.co/api/v2/pokemon/"

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

type locationAreaDetail struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
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
}

func fetchPokemon(name string, cache *pokecache.Cache) (Pokemon, error) {
	url := pokemonURL + name
	if data, ok := cache.Get(url); ok {
		var result Pokemon
		if err := json.Unmarshal(data, &result); err != nil {
			return Pokemon{}, err
		}
		return result, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Pokemon{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Pokemon{}, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return Pokemon{}, err
	}
	cache.Add(url, data)

	return result, nil
}

func fetchLocationArea(name string, cache *pokecache.Cache) (locationAreaDetail, error) {
	url := baseURL + name
	if data, ok := cache.Get(url); ok {
		fmt.Println("(cache hit)")
		var result locationAreaDetail
		if err := json.Unmarshal(data, &result); err != nil {
			return locationAreaDetail{}, err
		}
		return result, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return locationAreaDetail{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return locationAreaDetail{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result locationAreaDetail
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return locationAreaDetail{}, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return locationAreaDetail{}, err
	}
	cache.Add(url, data)

	return result, nil
}
