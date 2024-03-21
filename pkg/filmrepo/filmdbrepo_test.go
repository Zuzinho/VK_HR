package filmrepo

import (
	"VK_HR/pkg/customtime"
	"context"
	"github.com/lib/pq"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertFilm(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewFilmsDBRepository(db)
	ctx := context.Background()

	testFilm := &Film{
		Name:        "Test Film",
		Description: "A test film",
		PremierDate: customtime.CustomTime{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		Rating:      5.5,
	}

	mock.ExpectQuery("insert into films").
		WithArgs(testFilm.Name, testFilm.Description, testFilm.PremierDate.Format("2006-01-02"), testFilm.Rating).
		WillReturnRows(sqlmock.NewRows([]string{"film_id"}).AddRow(1))

	insertedID, err := repo.InsertFilm(ctx, testFilm)

	assert.NoError(t, err)
	assert.Equal(t, 1, insertedID)
}

func TestInsertActorsByFilm(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewFilmsDBRepository(db)
	ctx := context.Background()
	filmID := 1
	actorsID := []int{2, 3, 4}

	mock.ExpectExec("call insert_all_actors").
		WithArgs(filmID, pq.Array(actorsID)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.InsertActorsByFilm(ctx, filmID, actorsID)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewFilmsDBRepository(db)
	ctx := context.Background()
	setClause := "name = $1, rating = $2"
	args := []any{"New Name", 8.5}
	filmID := 1

	mock.ExpectPrepare("update films set").
		ExpectExec().
		WithArgs(args[0], args[1], filmID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(ctx, &setClause, &args, filmID)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewFilmsDBRepository(db)
	ctx := context.Background()
	filmID := 1

	mock.ExpectExec("delete from films where film_id =").
		WithArgs(filmID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(ctx, filmID)
	assert.NoError(t, err)
}

func TestSelectByFragment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewFilmsDBRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"film_id", "name", "description", "premier_date", "rating"}).
		AddRow(1, "Test Film", "Description", "2020-01-01", 5.5)

	mock.ExpectQuery("SELECT DISTINCT").
		WithArgs("%fragment%", "%fragment%").
		WillReturnRows(rows)

	films, err := repo.SelectByFragment(ctx, "fragment", "fragment")
	assert.NoError(t, err)
	assert.NotNil(t, films)
	assert.Len(t, *films, 1)
}

func TestSelectBySorting(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewFilmsDBRepository(db)
	ctx := context.Background()

	config := &SortingConfig{
		ColumnName: "rating",
		Direction:  "DESC",
	}

	rows := sqlmock.NewRows([]string{"film_id", "name", "description", "premier_date", "rating"}).
		AddRow(1, "Film A", "Description A", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), 9.5).
		AddRow(2, "Film B", "Description B", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), 8.5)

	mock.ExpectQuery("SELECT film_id, name, description, premier_date, rating FROM films ORDER BY (.+) DESC").
		WithArgs("rating").
		WillReturnRows(rows)

	films, err := repo.SelectBySorting(ctx, config)

	assert.NoError(t, err)
	assert.NotNil(t, films)
	assert.Len(t, *films, 2)
	assert.Equal(t, "Film A", (*films)[0].Name)
	assert.Equal(t, float32(9.5), (*films)[0].Rating)
}
