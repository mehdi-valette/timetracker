package entity

import (
	"errors"
	"strings"
)

type ProperNoun string

var EmptyNounErr = errors.New("noun is empty")

func CreateProperNoun(newName string) (ProperNoun, error) {
	name := strings.Trim(newName, " ")

	if name == "" {
		return ProperNoun(""), EmptyNounErr
	}

	return ProperNoun(
		strings.Trim(newName, " "),
	), nil
}
