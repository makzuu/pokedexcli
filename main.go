package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name string
	desc string
	cb   func() error
}

var commands = map[string]cliCommand{
	"help": {
		name: "help",
		desc: "Displays a help message",
		cb:   commandHelp,
	},
	"exit": {
		name: "exit",
		desc: "Exits the Pokedex",
		cb:   commandExit,
	},
}

func commandHelp() error {
	fmt.Println("")

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("\tUsage:")
	fmt.Println("\t- help: Displays a help message")
	fmt.Println("\t- exit: Exits the Pokedex")

	fmt.Println("")
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func main() {
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
		v.cb()
	}
}
