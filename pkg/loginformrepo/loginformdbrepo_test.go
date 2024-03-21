package loginformrepo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoginFormsDBRepository(db)
	ctx := context.Background()

	form := &LoginForm{
		Login:    "user1",
		Password: "pass123",
	}

	mock.ExpectExec("insert into login_forms").
		WithArgs(form.Login, form.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SignUp(ctx, form)
	assert.NoError(t, err)
}

func TestSignIn(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoginFormsDBRepository(db)
	ctx := context.Background()

	form := &LoginForm{
		Login:    "user1",
		Password: "pass123",
	}
	exist := true

	mock.ExpectQuery("select is_exist").
		WithArgs(form.Login, form.Password).
		WillReturnRows(sqlmock.NewRows([]string{"exist"}).AddRow(exist))

	existRes, err := repo.SignIn(ctx, form)
	assert.NoError(t, err)
	assert.Equal(t, exist, existRes)
}
