package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"pokedexcli/internal/pokecache"
)

type config struct {
	next     *string
	previous *string
	cache    *pokecache.Cache
	pokedex  map[string]Pokemon
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokémon",
			callback:    commandCatch(cfg),
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokémon",
			callback:    commandInspect(cfg),
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught Pokémon",
			callback:    commandPokedex(cfg),
		},
	}
}

func commandHelp(cfg *config) func(args []string) error {
	return func(args []string) error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("Usage:")
		fmt.Println()
		for _, cmd := range getCommands(cfg) {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}

func commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) func(args []string) error {
	return func(args []string) error {
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

func commandMapb(cfg *config) func(args []string) error {
	return func(args []string) error {
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

func commandExplore(cfg *config) func(args []string) error {
	return func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: explore <location-area>")
		}
		area, err := fetchLocationArea(args[0], cfg.cache)
		if err != nil {
			return err
		}
		fmt.Printf("Exploring %s...\n", area.Name)
		fmt.Println("Found Pokémon:")
		for _, enc := range area.PokemonEncounters {
			fmt.Printf(" - %s\n", enc.Pokemon.Name)
		}
		return nil
	}
}

func commandCatch(cfg *config) func(args []string) error {
	return func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: catch <pokemon>")
		}
		name := args[0]
		fmt.Printf("Throwing a Pokeball at %s...\n", name)

		pokemon, err := fetchPokemon(name, cfg.cache)
		if err != nil {
			return err
		}

		// Higher base experience → smaller catch window → harder to catch.
		// rand.Intn(baseExp + 50) < 50 gives ~58% for Caterpie (36), ~31% for
		// Pikachu (112), ~13% for Mewtwo (340).
		if rand.Intn(pokemon.BaseExperience+50) < 50 {
			fmt.Printf("%s was caught!\n", name)
			fmt.Println("You may now inspect it with the inspect command.")
			cfg.pokedex[name] = pokemon
		} else {
			fmt.Printf("%s escaped!\n", name)
		}
		return nil
	}
}

func commandPokedex(cfg *config) func(args []string) error {
	return func(args []string) error {
		fmt.Println("Your Pokedex:")
		for name := range cfg.pokedex {
			fmt.Printf(" - %s\n", name)
		}
		return nil
	}
}

func commandInspect(cfg *config) func(args []string) error {
	return func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: inspect <pokemon>")
		}
		pokemon, ok := cfg.pokedex[args[0]]
		if !ok {
			fmt.Println("you have not caught that pokemon")
			return nil
		}
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, s := range pokemon.Stats {
			fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}
		return nil
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
