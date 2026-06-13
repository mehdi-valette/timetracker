package main

import (
	"fmt"
	"os"

	"github.com/mehdi-valette/timetracker/internal/command"
)

func main() {
	home, homeErr := os.UserHomeDir()

	if homeErr != nil {
		panic("cannot find path to home directory")
	}

	directory := home + "/.local/share/timetracker"

	dirErr := os.MkdirAll(directory, 0700)

	if dirErr != nil {
		panic("cannot create the directory " + directory)
	}

	databasePath := directory + "/db.sqlite"

	_, err := command.Run(databasePath)

	if err != nil {
		fmt.Print(err)
	}
}
