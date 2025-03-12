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
var pokedex map[string]pokeapi.Pokemon

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
		"inspect": {
			name: "inspect <pokemon>",
			description: "Prints stats about a given caught Pokemon species",
			callback: commandInspect,
			config: nil,
		},
	}

	pokedex = map[string]pokeapi.Pokemon{}
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
		// if it's not registered to the Pokedex yet, register it
		if _, ok := pokedex[p]; ok != true {
			fmt.Println("New pokedex data is being registered for " + p)
			data, err := pokeapi.GetSpeciesData(p)
			if err != nil {
				return err
			}
			pokedex[p] = data
		}
	} else {
		fmt.Println(p + " escaped!")
	}
	return err
}

// print the stats of the registered species
func commandInspect(params []string) error {
	if len(params) == 0 {
		return errors.New("inspect requires a pokemon parameter")
	}
	data, ok := pokedex[params[0]]
	if !ok {
		return errors.New("you have not caught that pokemon")
	}
	fmt.Println(fmt.Sprintf("Name: %s", data.Name))
	fmt.Println(fmt.Sprintf("Height: %d", data.Height))
	fmt.Println(fmt.Sprintf("Weight: %d", data.Weight))
	fmt.Println("Stats:")
	for _, stat := range data.Stats {
		fmt.Println(fmt.Sprintf("  -%s: %d", stat.Stat.Name, stat.BaseStat))
	}
	fmt.Println("Types:")
	for _, pokeType := range data.Types {
		fmt.Println(fmt.Sprintf("  -%s", pokeType.Type.Name))
	}

	return nil
}