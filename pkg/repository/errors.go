package repository

import (
	"fmt"
)

type EmptyErr struct {
}

func (e EmptyErr) Error() string {
	return fmt.Sprint("empty response")
}

type DuplicateErr struct {
	Pk  string
	Err error
}

func (d DuplicateErr) Error() string {
	return fmt.Sprintf("duplicate %s", d.Pk)
}

type InvalidErr struct {
	Title string
}

func (i InvalidErr) Error() string {
	return fmt.Sprint(i.Title)
}
