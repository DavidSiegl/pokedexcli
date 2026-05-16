package main

import (
	"fmt"
	"os"
	"strings"

	"pokedexcli/internal/pokecache"
)

type config struct {
	next     *string
	previous *string
	cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

func getCommands(cfg *config) map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp(cfg),
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas",
			callback:    commandMap(cfg),
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb(cfg),
		},
		"explore": {
			name:        "explore",
			description: "Lists all Pokémon in a location area",
			callback:    commandExplore(cfg),
		},
	}
}

func commandHelp(cfg *config) func() error {
	return func() error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("Usage:")
		fmt.Println()
		for _, cmd := range getCommands(cfg) {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) func() error {
	return func() error {
		url := baseURL
		if cfg.next != nil {
			url = *cfg.next
		}

		result, err := fetchLocationAreas(url, cfg.cache)
		if err != nil {
			return err
		}

		cfg.next = result.Next
		cfg.previous = result.Previous

		for _, area := range result.Results {
			fmt.Println(area.Name)
		}
		return nil
	}
}

func commandMapb(cfg *config) func() error {
	return func() error {
		if cfg.previous == nil {
			return fmt.Errorf("you're on the first page")
		}

		result, err := fetchLocationAreas(*cfg.previous, cfg.cache)
		if err != nil {
			return err
		}

		cfg.next = result.Next
		cfg.previous = result.Previous

		for _, area := range result.Results {
			fmt.Println(area.Name)
		}
		return nil
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
