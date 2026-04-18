package entity

import "strings"

type Name string

func CreateName(newName string) Name {
	return Name(
		strings.Trim(newName, " "),
	)
}
