package entity

import "strings"

type ProperNoun string

func CreateProperNoun(newName string) ProperNoun {
	return ProperNoun(
		strings.Trim(newName, " "),
	)
}
