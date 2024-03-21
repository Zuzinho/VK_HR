package userrepo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserDBRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserDBRepository(db)
	ctx := context.Background()

	// Мокируем ожидаемое действие
	mock.ExpectExec("insert into users").
		WithArgs("user1", RegularUser).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Insert(ctx, "user1", RegularUser)
	assert.NoError(t, err)
}

func TestUserDBRepository_SelectRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserDBRepository(db)
	ctx := context.Background()

	// Подготовка мокированного ответа
	rows := sqlmock.NewRows([]string{"role"}).AddRow(RegularUser)
	mock.ExpectQuery("select role from users where login =").
		WithArgs("user1").
		WillReturnRows(rows)

	role, err := repo.SelectRole(ctx, "user1")
	assert.NoError(t, err)
	assert.Equal(t, RegularUser, role)
}
