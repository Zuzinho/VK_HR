package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilmsValidator_IsValidValue(t *testing.T) {
	validator := NewFilmsValidator()

	testCases := []struct {
		name        string
		columnName  ColumnName
		value       string
		expectError bool
	}{
		{
			name:        "Valid film name",
			columnName:  FilmName,
			value:       "Valid Name",
			expectError: false,
		},
		{
			name:        "Invalid film name - empty",
			columnName:  FilmName,
			value:       "",
			expectError: true,
		},
		{
			name:        "Valid film description",
			columnName:  FilmDescription,
			value:       "A valid description that is within 1000 characters.",
			expectError: false,
		},
		{
			name:        "Valid film premier date",
			columnName:  FilmPremierDate,
			value:       "2020-01-01",
			expectError: false,
		},
		{
			name:        "Invalid film premier date",
			columnName:  FilmPremierDate,
			value:       "01-01-2020",
			expectError: true,
		},
		{
			name:        "Valid film rating",
			columnName:  FilmRating,
			value:       "8.5",
			expectError: false,
		},
		{
			name:        "Invalid film rating - out of range",
			columnName:  FilmRating,
			value:       "11",
			expectError: true,
		},
		{
			name:        "Invalid column name",
			columnName:  "unknown",
			value:       "any",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validator.IsValidValue(tc.columnName, tc.value)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
