package entity

import (
	"errors"
	"strings"
)

var ShortNameWithSpaceErr = errors.New("a short name should contain no space")

type ShortName string

func CreateShortName(name string) ShortName {
	return ShortName(
		strings.ToLower(
			strings.ReplaceAll(name, " ", ""),
		),
	)
}

func (s ShortName) String() string {
	return string(s)
}
