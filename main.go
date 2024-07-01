package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/makzuu/pokedexcli/internal/pokeapi"
	"math/rand"
	"os"
	"strings"
)

type cliCommand struct {
	name     string
	desc     string
	callback func(*config, string) error
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
	"explore": {
		name:     "explore",
		desc:     "Displays pokemons in the expecified area",
		callback: commandExplore,
	},
	"catch": {
		name:     "catch",
		desc:     "...",
		callback: commandCatch,
	},
}

var pokedex = map[string]pokeapi.Pokemon{}

type config struct {
	next, previous string
}

func commandHelp(c *config, param string) error {
	fmt.Println("")

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("\tUsage:")
	fmt.Println("\t- help: Displays a help message")
	fmt.Println("\t- exit: Exits the Pokedex")
	fmt.Println("\t- map: Displays the next 20 locations")
	fmt.Println("\t- mapb: Displays the previous 20 locations")
	fmt.Println("\t- explore: Displays pokemons in the expecified area")
	fmt.Println("\t- catch: Try to catch a pokemon")

	fmt.Println("")
	return nil
}

func commandExit(c *config, param string) error {
	os.Exit(0)
	return nil
}

func commandMap(c *config, param string) error {
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

func commandMapb(c *config, param string) error {
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

func commandExplore(c *config, areaName string) error {
	if areaName == "" {
		return fmt.Errorf("explore command expects the area name")
	}
	pokemons, err := pokeapi.GetPokemons(areaName)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", areaName)
	fmt.Println("Found Pokemon:")
	for _, name := range pokemons {
		fmt.Println(" - ", name)
	}
	return nil
}

func commandCatch(c *config, pokemonName string) error {
	if pokemonName == "" {
		return fmt.Errorf("catch command expects pokemon name")
	}
	pokemon, err := pokeapi.PokemonInfo(pokemonName)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if float64(rand.Intn(100)) > (float64(pokemon.BaseExperience) * .74) {
		pokedex[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
		return nil
	}
	fmt.Printf("%s escaped!\n", pokemon.Name)
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
		commandFields := strings.Fields(command)
		if len(commandFields) < 1 {
			fmt.Println("invalid command")
			continue
		}
		commandName := commandFields[0]
		param := ""
		if len(commandFields) == 2 {
			param = commandFields[1]
		}
		v, ok := commands[commandName]
		if !ok {
			fmt.Println("invalid command")
			continue
		}
		err := v.callback(&conf, param)
		if err != nil {
			fmt.Println(err)
		}
	}
}
