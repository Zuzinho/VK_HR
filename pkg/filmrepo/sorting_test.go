package filmrepo

import (
	"testing"

	"VK_HR/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestNewSortingConfig(t *testing.T) {
	tests := []struct {
		name           string
		columnName     validator.ColumnName
		direction      SortingDirection
		expectedError  error
		expectedConfig *SortingConfig
	}{
		{
			name:          "Valid config with DESC direction",
			columnName:    validator.FilmRating,
			direction:     DESC,
			expectedError: nil,
			expectedConfig: &SortingConfig{
				ColumnName: validator.FilmRating,
				Direction:  DESC,
			},
		},
		{
			name:           "Invalid direction",
			columnName:     validator.FilmRating,
			direction:      "INVALID",
			expectedError:  newUnknownSortingDirectionError("INVALID"),
			expectedConfig: nil,
		},
		{
			name:           "Invalid column name",
			columnName:     "invalid_column",
			direction:      ASC,
			expectedError:  newNoAccessSortingColumnNameError("invalid_column"),
			expectedConfig: nil,
		},
		{
			name:          "Empty direction defaults to DESC",
			columnName:    validator.FilmName,
			direction:     "",
			expectedError: nil,
			expectedConfig: &SortingConfig{
				ColumnName: validator.FilmName,
				Direction:  DESC,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewSortingConfig(tt.columnName, tt.direction)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedConfig, config)
			}
		})
	}
}

func TestSortingDirectionIsValid(t *testing.T) {
	tests := []struct {
		direction     SortingDirection
		expectedError error
	}{
		{ASC, nil},
		{DESC, nil},
		{"INVALID", newUnknownSortingDirectionError("INVALID")},
	}

	for _, tt := range tests {
		err := tt.direction.IsValid()
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError.Error(), err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
