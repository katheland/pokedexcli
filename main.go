package main
import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"internal/pokeapi"
)

// The structure of a valid command
type cliCommand struct {
	name string
	description string
	callback func() error
	config *Config
}

type Config struct {
	next string
	previous *string
}

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
	}
}

// main
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		//fmt.Println("Your command was:", input[0])
		command, exists := validCommands[input[0]]
		if exists {
			err := command.callback()
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
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// print the help text
func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, val := range validCommands {
		fmt.Println(fmt.Sprintf("%s: %s", val.name, val.description))
	}
	return nil
}

// print the next page of locations
func commandMap() error {
	next, previous, err := pokeapi.MapFunction(&validCommands["map"].config.next)
	if err != nil {
		return err
	}

	validCommands["map"].config.next = next
	validCommands["map"].config.previous = previous
	return nil
}

// print the previous page of locations
func commandMapB() error {
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