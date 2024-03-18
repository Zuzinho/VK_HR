package validator

import "fmt"

type InvalidColumnNameError struct {
	columnName ColumnName
}

func newInvalidColumnNameError(name ColumnName) InvalidColumnNameError {
	return InvalidColumnNameError{
		columnName: name,
	}
}

func (err InvalidColumnNameError) Error() string {
	return fmt.Sprintf("no column by name '%s'", err.columnName)
}

type InvalidValueError struct {
	columnName ColumnName
	value      any
}

func newInvalidValueError(name ColumnName, value any) InvalidValueError {
	return InvalidValueError{
		columnName: name,
		value:      value,
	}
}

func (err InvalidValueError) Error() string {
	return fmt.Sprintf("column by name '%s' can not be equal '%v'", err.columnName, err.value)
}

var (
	InvalidColumnNameErr = InvalidColumnNameError{}
	InvalidValueErr      = InvalidValueError{}
)
