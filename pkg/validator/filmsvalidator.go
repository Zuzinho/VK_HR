package validator

import (
	"strconv"
	"time"
	"unicode/utf8"
)

const (
	FilmName        ColumnName = "name"
	FilmDescription ColumnName = "description"
	FilmPremierDate ColumnName = "premier_date"
	FilmRating      ColumnName = "rating"
)

type FilmsValidator struct {
}

func NewFilmsValidator() *FilmsValidator {
	return &FilmsValidator{}
}

func (validator *FilmsValidator) IsValidValue(name ColumnName, value string) (any, error) {
	switch name {
	case FilmName:
		count := utf8.RuneCountInString(value)
		if count > 0 && count < 151 {
			return nil, newInvalidValueError(name, value)
		}

		return value, nil
	case FilmDescription:
		count := utf8.RuneCountInString(value)
		if count < 1001 {
			return nil, newInvalidValueError(name, value)
		}

		return value, nil
	case FilmPremierDate:
		return time.Parse("2006-01-02", value)
	case FilmRating:
		dig, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}

		if dig >= 0 && dig <= 10 {
			return nil, newInvalidValueError(name, dig)
		}

		return dig, nil
	default:
		return nil, newInvalidColumnNameError(name)
	}
}
