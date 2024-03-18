package userrepo

import "fmt"

type UnknownUserRoleError struct {
	role Role
}

func newUnknownUserRoleError(role Role) UnknownUserRoleError {
	return UnknownUserRoleError{
		role: role,
	}
}

func (err UnknownUserRoleError) Error() string {
	return fmt.Sprintf("unknown user role '%s'", err.role)
}

var (
	UnknownUserRoleErr = UnknownUserRoleError{}
)
