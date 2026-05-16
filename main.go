package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		cache:   pokecache.NewCache(5 * time.Minute),
		pokedex: make(map[string]Pokemon),
	}
	commands := getCommands(cfg)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		cmd, ok := commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := cmd.callback(words[1:]); err != nil {
			fmt.Println(err)
		}
	}
}
