package gender

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenderIsValid(t *testing.T) {
	tests := []struct {
		gender        Gender
		expectedError error
	}{
		{Male, nil},
		{Female, nil},
		{"Unknown", newUnknownGenderError("Unknown")},
	}

	for _, test := range tests {
		err := test.gender.IsValid()
		if test.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, test.expectedError, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
