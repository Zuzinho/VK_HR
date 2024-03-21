package userrepo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRole_IsValid(t *testing.T) {
	testCases := []struct {
		name        string
		role        Role
		expectError bool
	}{
		{
			name:        "Valid role: RegularUser",
			role:        RegularUser,
			expectError: false,
		},
		{
			name:        "Valid role: Admin",
			role:        Admin,
			expectError: false,
		},
		{
			name:        "Invalid role",
			role:        "Guest",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.role.IsValid()
			if tc.expectError {
				assert.Error(t, err)
				// Проверяем также, что возвращается ожидаемый тип ошибки
				assert.IsType(t, newUnknownUserRoleError(""), err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
