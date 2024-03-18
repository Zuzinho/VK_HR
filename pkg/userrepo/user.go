package userrepo

import "context"

type Role string

const (
	RegularUser Role = "Regular User"
	Admin       Role = "Admin"
)

func (role Role) IsValid() error {
	switch role {
	case RegularUser, Admin:
		return nil
	default:
		return newUnknownUserRoleError(role)
	}
}

type User struct {
	Login string `json:"login"`
	Role  Role   `json:"role"`
}

type UsersRepository interface {
	Insert(ctx context.Context, login string, role Role) error
	SelectRole(ctx context.Context, login string) (Role, error)
}
