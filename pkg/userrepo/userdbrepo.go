package userrepo

import (
	"context"
	"database/sql"
)

type UserDBRepository struct {
	DB *sql.DB
}

func NewUserDBRepository(db *sql.DB) *UserDBRepository {
	return &UserDBRepository{
		DB: db,
	}
}

func (repo *UserDBRepository) Insert(ctx context.Context, login string, role Role) error {
	_, err := repo.DB.ExecContext(ctx, "insert into users (login, role) values ($1, $2)", login, role)

	return err
}

func (repo *UserDBRepository) SelectRole(ctx context.Context, login string) (Role, error) {
	var role Role
	err := repo.DB.QueryRowContext(ctx, "select role from users where login = $1", login).Scan(&role)

	return role, err
}
