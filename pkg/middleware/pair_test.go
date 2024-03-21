package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPair(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		method   string
		expected Pair
	}{
		{
			name:   "GET request to root",
			path:   "/",
			method: "GET",
			expected: Pair{
				path:   "/",
				method: "GET",
			},
		},
		{
			name:   "POST request to /login",
			path:   "/login",
			method: "POST",
			expected: Pair{
				path:   "/login",
				method: "POST",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NewPair(test.path, test.method)

			assert.Equal(t, test.expected, result)
		})
	}
}
