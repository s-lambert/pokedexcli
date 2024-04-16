package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func allowedCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func() error {
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
		callback: func() error {
			return errors.New("Exit")
		},
	}
	return commands
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := allowedCommands()

	fmt.Printf("Pokedex > ")
	for scanner.Scan() {
		command := scanner.Text()
		cliCommand, ok := commands[command]
		if !ok {
			fmt.Printf("Unknown command: current commands are help, exit\n")
			fmt.Printf("Pokedex > ")
			continue
		}
		err := cliCommand.callback()
		if err != nil {
			break
		}
		fmt.Printf("Pokedex > ")
	}
}
