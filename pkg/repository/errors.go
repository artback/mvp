package repository

import (
	"fmt"
)

type EmptyError struct{}

func (e EmptyError) Error() string {
	return "empty response"
}

type DuplicateError struct {
	Constraint string
	Err        error
}

func (d DuplicateError) Error() string {
	return fmt.Sprintf("duplicate %s", d.Constraint)
}

type InvalidError struct {
	Title string
}

func (i InvalidError) Error() string {
	return fmt.Sprint(i.Title)
}
