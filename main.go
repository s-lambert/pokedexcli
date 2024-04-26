package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/s-lambert/pokedexcli/internal/pokecache"
)

type cliConfig struct {
	cache *pokecache.Cache
	Next  string
	Prev  *string
}

func GetLocationAreasUrl(offset int) string {
	return fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?take=20&offset=%v", offset)
}

type cliCommand struct {
	name        string
	description string
	callback    func(*cliConfig) error
}

func allowedCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(_ *cliConfig) error {
			fmt.Printf("\nWelcome to the Pokedex!\n")
			fmt.Printf("Usage:\n\n")
			for key := range commands {
				command := commands[key]
				fmt.Printf("%s: %s\n", command.name, command.description)
			}
			fmt.Printf("\n")
			return nil
		},
	}

	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback: func(_ *cliConfig) error {
			return errors.New("exit")
		},
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "Display next 20 locations",
		callback: func(cfg *cliConfig) error {
			ls, err := GetLocationAreas(cfg, false)
			if err != nil {
				log.Fatal(err)
			}
			for _, l := range ls {
				fmt.Printf("%s\n", l.Name)
			}
			return nil
		},
	}

	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Display the previous 20 locations",
		callback: func(cfg *cliConfig) error {
			if cfg.Prev == nil {
				return errors.New("cannot go back to previous loations, at the start of the list")
			}
			ls, err := GetLocationAreas(cfg, true)
			if err != nil {
				log.Fatal(err)
			}
			for _, l := range ls {
				fmt.Printf("%s\n", l.Name)
			}
			return nil
		},
	}

	return commands
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := allowedCommands()
	cache := pokecache.NewCache()
	config := cliConfig{
		cache: cache,
		Next:  GetLocationAreasUrl(0),
		Prev:  nil,
	}

	fmt.Printf("Pokedex > ")
	for scanner.Scan() {
		command := scanner.Text()
		cliCommand, ok := commands[command]
		if !ok {
			fmt.Printf("Unknown command: current commands are help, exit\n")
			fmt.Printf("Pokedex > ")
			continue
		}
		err := cliCommand.callback(&config)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("Pokedex > ")
	}
}
