package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/s-lambert/pokedexcli/internal/pokecache"
)

type cliConfig struct {
	cache       *pokecache.Cache
	Next        string
	Prev        *string
	AreaInfo    string
	PokemonInfo string
	Pokedex     map[string]PokemonInfo
}

func GetLocationAreasUrl(offset int) string {
	return fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?take=20&offset=%v", offset)
}

type cliCommand struct {
	name        string
	description string
	callback    func(*cliConfig, string) error
}

func allowedCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(_ *cliConfig, _ string) error {
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
		callback: func(_ *cliConfig, _ string) error {
			return errors.New("exit")
		},
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "Display next 20 locations",
		callback: func(cfg *cliConfig, _ string) error {
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
		callback: func(cfg *cliConfig, _ string) error {
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

	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Visit a location",
		callback: func(cfg *cliConfig, locationName string) error {
			if locationName == "" {
				return errors.New("explore expects a location name")
			}

			la, err := GetLocationArea(cfg, locationName)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Exploring %s...\n", la.Name)
			fmt.Print("Found Pokemon:\n")
			for _, p := range la.PokemonEncounters {
				fmt.Printf("  - %s\n", p.Pokemon.Name)
			}

			return nil
		},
	}

	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Attempt to catch a pokemon",
		callback: func(cfg *cliConfig, pokemonName string) error {
			if pokemonName == "" {
				return errors.New("catch expects a pokemon name")
			}

			p, err := GetPokemonInfo(cfg, pokemonName)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Throwing a Pokeball at %s...\n", p.Name)

			catchCheck := 100 - rand.Intn(p.BaseExperience)
			if catchCheck > 0 {
				fmt.Printf("%s was caught!\n", p.Name)
				cfg.Pokedex[p.Name] = p
			} else {
				fmt.Printf("%s escaped!\n", p.Name)
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
		cache:       cache,
		Next:        GetLocationAreasUrl(0),
		Prev:        nil,
		AreaInfo:    "https://pokeapi.co/api/v2/location-area/",
		PokemonInfo: "https://pokeapi.co/api/v2/pokemon/",
		Pokedex:     make(map[string]PokemonInfo),
	}

	fmt.Printf("Pokedex > ")
	for scanner.Scan() {
		args := strings.Split(scanner.Text(), " ")
		if args[0] == "" {
			continue
		}
		command := args[0]
		cliCommand, ok := commands[command]
		if !ok {
			fmt.Printf("Unknown command: current commands are help, exit\n")
			fmt.Printf("Pokedex > ")
			continue
		}
		err := cliCommand.callback(&config, strings.Join(args[1:], ""))
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("Pokedex > ")
	}
}
