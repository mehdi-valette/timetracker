package main

import (
	"fmt"

	"github.com/mehdi-valette/timetracker/internal/command"
)

func main() {
	_, err := command.Run()

	if err != nil {
		fmt.Print(err)
	}
}
