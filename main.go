package main
import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"net/http"
	"io"
	"encoding/json"
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
//var mapConfig Config

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
			description: "Lists areas",
			callback: commandMap,
			config: &mapConfig,
		},
		"mapb": {
			name: "mapb",
			description: "Lists areas",
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

type Location struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// print a page of 20 locations
func mapFunction(url *string) error {
	if url == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	res, err := http.Get(*url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var locations Location
	if err = json.Unmarshal(jsonData, &locations); err != nil {
		return err
	}

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	validCommands["map"].config.next = locations.Next
	validCommands["map"].config.previous = locations.Previous
	
	return nil
}

// print the next page of locations
func commandMap() error {
	return mapFunction(&validCommands["map"].config.next)
}

// print the previous page of locations
func commandMapB() error {
	return mapFunction(validCommands["map"].config.previous)
}