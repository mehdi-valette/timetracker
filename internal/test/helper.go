package test

import "fmt"

func NoError(err error) error {
	return fmt.Errorf("should not return an error, got: %w", err)
}
