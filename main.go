package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/makzuu/pokedexcli/pokeapi"
	"os"
)

type cliCommand struct {
	name     string
	desc     string
	callback func(*config) error
}

var commands = map[string]cliCommand{
	"help": {
		name:     "help",
		desc:     "Displays a help message",
		callback: commandHelp,
	},
	"exit": {
		name:     "exit",
		desc:     "Exits the Pokedex",
		callback: commandExit,
	},
	"map": {
		name:     "map",
		desc:     "Displays the next 20 locations",
		callback: commandMap,
	},
	"mapb": {
		name:     "mapb",
		desc:     "Displays the previous 20 locations",
		callback: commandMapb,
	},
}

type config struct {
	next, previous string
}

func commandHelp(c *config) error {
	fmt.Println("")

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("\tUsage:")
	fmt.Println("\t- help: Displays a help message")
	fmt.Println("\t- exit: Exits the Pokedex")
	fmt.Println("\t- map: Displays the next 20 locations")
	fmt.Println("\t- mapb: Displays the previous 20 locations")

	fmt.Println("")
	return nil
}

func commandExit(c *config) error {
	os.Exit(0)
	return nil
}

func commandMap(c *config) error {
	l, err := pokeapi.GetLocations(c.next)
	if err != nil {
		return err
	}
	c.previous = l.Previous
	c.next = l.Next
	for _, locations := range l.Locations {
		fmt.Println(locations.Name)
	}
	return nil
}

func commandMapb(c *config) error {
	if c.previous == "" {
		return errors.New("you are on the first page")
	}
	l, err := pokeapi.GetLocations(c.previous)
	if err != nil {
		return err
	}
	c.previous = l.Previous
	c.next = l.Next
	for _, locations := range l.Locations {
		fmt.Println(locations.Name)
	}
	return nil
}

func main() {
	var conf config = config{
		previous: "",
		next:     "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
	}
	sc := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if ok := sc.Scan(); !ok {
			if err := sc.Err(); err != nil {
				fmt.Println(err)
			}
			break
		}
		command := sc.Text()
		v, ok := commands[command]
		if !ok {
			fmt.Println("invalid command")
			continue
		}
		if err := v.callback(&conf); err != nil {
			fmt.Println(err)
		}
	}
}
