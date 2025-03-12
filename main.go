package main
import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"internal/pokeapi"
	"errors"
)

// The structure of a valid command
type cliCommand struct {
	name string
	description string
	callback func([]string) error
	config *Config
}
type Config struct {
	next string
	previous *string
}

// A map of Pokemon the user has caught (to be expanded upon)
var pokedex map[string]bool

// A map of valid commands that the program recognizes, and what to do with them
var validCommands map[string]cliCommand

// but it's important to initialize the map here, to avoid an initialization loop
func init() {
	mapConfig := Config{
		next: "https://pokeapi.co/api/v2/location-area/",
		previous: nil,
	}

	validCommands = map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
			config: nil,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
			config: nil,
		},
		"map": {
			name: "map",
			description: "Lists the next page of areas",
			callback: commandMap,
			config: &mapConfig,
		},
		"mapb": {
			name: "mapb",
			description: "Lists the previous page of areas",
			callback: commandMapB,
			config: &mapConfig,
		},
		"explore": {
			name: "explore <area_name>",
			description: "Lists the Pokemon found at the given location",
			callback: commandExplore,
			config: nil,
		},
		"catch": {
			name: "catch <pokemon>",
			description: "Attempts to catch a Pokemon of the given species",
			callback: commandCatch,
			config: nil,
		},
	}

	pokedex = map[string]bool{}
}

// main
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		command, exists := validCommands[input[0]]
		if exists {
			params := []string{}
			if len(input) > 1 {
				params = input[1:]
			}
			err := command.callback(params)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

// removes whitespace from the input, splits it into words, and turns it all lowercase to clean it for use
func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

// exit the program
func commandExit(params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// print the help text
func commandHelp(params []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, val := range validCommands {
		fmt.Println(fmt.Sprintf("%s: %s", val.name, val.description))
	}
	return nil
}

// print the next page of locations
func commandMap(params []string) error {
	next, previous, err := pokeapi.MapFunction(&validCommands["map"].config.next)
	if err != nil {
		return err
	}

	validCommands["map"].config.next = next
	validCommands["map"].config.previous = previous
	return nil
}

// print the previous page of locations
func commandMapB(params []string) error {
	if validCommands["map"].config.previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	next, previous, err := pokeapi.MapFunction(validCommands["map"].config.previous)
	if err != nil {
		return err
	}

	validCommands["map"].config.next = next
	validCommands["map"].config.previous = previous
	return nil
}

// print the pokemon found at a given location
func commandExplore(params []string) error {
	if len(params) == 0 {
		return errors.New("explore requires a location parameter")
	}
	return pokeapi.ExploreFunction(params[0])
}

// attempt to catch the given species of pokemon
func commandCatch(params []string) error {
	if len(params) == 0 {
		return errors.New("catch requires a pokemon parameter")
	}
	p, b, err := pokeapi.CatchFunction(params[0])
	if b {
		fmt.Println(p + " was caught!")
		pokedex[p] = true
	} else {
		fmt.Println(p + " escaped!")
	}
	return err
}