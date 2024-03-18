package e

import "fmt"

// WrapErr wraps error in a template "msg: err".
func WrapErr(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
