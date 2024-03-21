package actorrepo

import (
	"VK_HR/pkg/customtime"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewActorsDBRepository(db)

	ctx := context.TODO()
	actor := &Actor{
		FirstName:  "John",
		SecondName: "Doe",
		Gender:     "Male",
		Birthday:   customtime.CustomTime{Time: time.Now()},
	}

	mock.ExpectQuery("insert into actors").
		WithArgs(actor.FirstName, actor.SecondName, actor.Gender, actor.Birthday.Format("2006-01-02")).
		WillReturnRows(sqlmock.NewRows([]string{"actor_id"}).AddRow(1))

	insertedID, err := repo.Insert(ctx, actor)
	require.NoError(t, err)
	require.Equal(t, 1, insertedID)
}

func TestUpdateActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActorsDBRepository(db)
	ctx := context.Background()

	setClause := "first_name = $1, second_name = $2"
	args := []any{"Jane", "Doe"}

	mock.ExpectPrepare("update actors set").
		ExpectExec().
		WithArgs(args[0], args[1], 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(ctx, &setClause, &args, 1)
	assert.NoError(t, err)
}

func TestDeleteActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActorsDBRepository(db)
	ctx := context.Background()

	mock.ExpectExec("delete from actors where actor_id =").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestSelectAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActorsDBRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"actor_id", "first_name", "second_name", "gender", "birthday", "film_id", "name", "description", "premier_date", "rating"}).
		AddRow(1, "John", "Doe", "male", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			1, "Film 1", "Description 1", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 8.5).
		AddRow(1, "John", "Doe", "male", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			2, "Film 2", "Description 2", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 9.0)

	mock.ExpectQuery("SELECT a.actor_id, a.first_name, a.second_name, a.gender, a.birthday, f.film_id, f.name, f.description, f.premier_date, f.rating FROM actors a LEFT JOIN actors_has_films ahf ON a.actor_id = ahf.actor_id LEFT JOIN films f ON ahf.film_id = f.film_id order by a.actor_id").
		WillReturnRows(rows)

	actors, err := repo.SelectAll(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, actors)
	assert.Len(t, *actors, 1)
	assert.Len(t, *(*actors)[0].Films, 2)

	firstFilm := (*(*actors)[0].Films)[0]
	assert.Equal(t, int32(1), firstFilm.FilmID)
	assert.Equal(t, "Film 1", firstFilm.Name)
}
