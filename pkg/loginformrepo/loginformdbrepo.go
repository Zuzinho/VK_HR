package loginformrepo

import (
	"context"
	"database/sql"
)

type LoginFormsDBRepository struct {
	DB *sql.DB
}

func NewLoginFormsDBRepository(db *sql.DB) *LoginFormsDBRepository {
	return &LoginFormsDBRepository{
		DB: db,
	}
}

func (repo *LoginFormsDBRepository) SignUp(ctx context.Context, form *LoginForm) error {
	_, err := repo.DB.ExecContext(ctx, "insert into login_forms (login, password) values ($1, $2)", form.Login, form.Password)

	return err
}

func (repo *LoginFormsDBRepository) SignIn(ctx context.Context, form *LoginForm) (bool, error) {
	var exist bool

	err := repo.DB.QueryRowContext(ctx, "return exists(select * from login_forms "+
		"where login = $1 and password = $2)", form.Login, form.Password).Scan(&exist)

	return exist, err
}
