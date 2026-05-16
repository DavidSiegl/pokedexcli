package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2/location-area/"

type locationAreaResponse struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func fetchLocationAreas(url string) (locationAreaResponse, error) {
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
	return result, nil
}
