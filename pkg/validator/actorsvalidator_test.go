package validator

import (
	"VK_HR/pkg/gender"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActorsValidator_IsValidValue(t *testing.T) {
	validator := NewActorsValidator()

	testCases := []struct {
		name          string
		columnName    ColumnName
		value         string
		expectError   bool
		expectedValue any
	}{
		{
			name:          "Valid first name",
			columnName:    ActorFirstName,
			value:         "John",
			expectError:   false,
			expectedValue: "John",
		},
		{
			name:          "Valid second name",
			columnName:    ActorSecondName,
			value:         "Doe",
			expectError:   false,
			expectedValue: "Doe",
		},
		{
			name:          "Valid gender (Male)",
			columnName:    ActorGender,
			value:         "Male",
			expectError:   false,
			expectedValue: gender.Male,
		},
		{
			name:        "Invalid gender",
			columnName:  ActorGender,
			value:       "Unknown",
			expectError: true,
		},
		{
			name:          "Valid birthday",
			columnName:    ActorBirthday,
			value:         "1990-01-01",
			expectError:   false,
			expectedValue: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "Invalid birthday format",
			columnName:  ActorBirthday,
			value:       "01-01-1990",
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
			value, err := validator.IsValidValue(tc.columnName, tc.value)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, value)
			}
		})
	}
}
