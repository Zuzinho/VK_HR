package filmrepo

import (
	"VK_HR/pkg/validator"
	"fmt"
)

type UnknownSortingDirectionError struct {
	direction SortingDirection
}

func newUnknownSortingDirectionError(direction SortingDirection) UnknownSortingDirectionError {
	return UnknownSortingDirectionError{
		direction: direction,
	}
}

func (err UnknownSortingDirectionError) Error() string {
	return fmt.Sprintf("unknown sortring direction '%s'", err.direction)
}

type NoAccessSortingColumnNameError struct {
	name validator.ColumnName
}

func newNoAccessSortingColumnNameError(name validator.ColumnName) NoAccessSortingColumnNameError {
	return NoAccessSortingColumnNameError{
		name: name,
	}
}

func (err NoAccessSortingColumnNameError) Error() string {
	return fmt.Sprintf("can`t sorting by column '%s'", err.name)
}

var (
	UnknownSortingDirectionErr   = UnknownSortingDirectionError{}
	NoAccessSortingColumnNameErr = NoAccessSortingColumnNameError{}
)
