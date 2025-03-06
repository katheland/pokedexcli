package main
import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

// The structure of a valid command
type cliCommand struct {
	name string
	description string
	callback func() error
}

// A map of valid commands that the program recognizes, and what to do with them
var validCommands map[string]cliCommand

// but it's important to initialize the map here, to avoid an initialization loop
func init() {
	validCommands = map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
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
			command.callback()
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