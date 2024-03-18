package gender

import (
	"fmt"
)

type UnknownGenderError struct {
	gender Gender
}

func newUnknownGenderError(gender Gender) UnknownGenderError {
	return UnknownGenderError{
		gender: gender,
	}
}

func (err UnknownGenderError) Error() string {
	return fmt.Sprintf("unknown gender '%s'", err.gender)
}

var (
	UnknownGenderErr = UnknownGenderError{}
)
