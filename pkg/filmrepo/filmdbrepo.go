package filmrepo

import (
	"context"
	"database/sql"
	"fmt"
)

type FilmsDBRepository struct {
	DB *sql.DB
}

func NewFilmsDBRepository(db *sql.DB) *FilmsDBRepository {
	return &FilmsDBRepository{
		DB: db,
	}
}

func (repo *FilmsDBRepository) InsertFilm(ctx context.Context, film *Film) (int, error) {
	var insertedID int

	err := repo.DB.QueryRowContext(ctx, "insert into films (name, description, premier_date, rating) values "+
		"($1, $2, $3, $4) returning film_id", film.Name, film.Description, film.PremierDate, film.Rating).Scan(&insertedID)

	return insertedID, err
}

func (repo *FilmsDBRepository) InsertActorsByFilm(ctx context.Context, filmID int, actorsID []int) error {
	_, err := repo.DB.ExecContext(ctx, "call procedure insert_all_actors($1, $2)", filmID, actorsID)

	return err
}

func (repo *FilmsDBRepository) Update(ctx context.Context, setClause *string, args *[]any, id int) error {
	queryString := fmt.Sprintf("update films set %s where film_id = $%d", *setClause, len(*args)+1)

	smtm, err := repo.DB.PrepareContext(ctx, queryString)
	defer smtm.Close()
	if err != nil {
		return err
	}

	_, err = smtm.ExecContext(ctx, append(*args, id)...)
	return err
}

func (repo *FilmsDBRepository) Delete(ctx context.Context, id int) error {
	_, err := repo.DB.ExecContext(ctx, "delete from films where film_id = $1", id)

	return err
}

func (repo *FilmsDBRepository) SelectBySorting(ctx context.Context, config *SortingConfig) (*Films, error) {
	return repo.selectByQuery(ctx, "select film_id, name, description, premier_date, rating from films "+
		"order by $1 $2", config.ColumnName, config.Direction)
}

func (repo *FilmsDBRepository) SelectByFragment(ctx context.Context, filmFragment, nameFragment string) (*Films, error) {
	return repo.selectByQuery(ctx, "SELECT f.film_id, f.name, f.description, f.premier_date, f.rating "+
		"FROM films f JOIN actors_has_films af ON f.film_id = af.film_id "+
		"JOIN actors a ON a.actor_id = af.actor_id WHERE f.name LIKE '%$1%' "+
		"AND a.name LIKE '%$2%'", filmFragment, nameFragment)
}

func (repo *FilmsDBRepository) selectByQuery(ctx context.Context, query string, args ...any) (*Films, error) {
	rows, err := repo.DB.QueryContext(ctx, query, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	films := make(Films, 0)
	for rows.Next() {
		film := &Film{}
		err = rows.Scan(&film.FilmID, &film.Name, &film.Description, &film.PremierDate, &film.Rating)
		if err != nil {
			continue
		}

		films.Append(film)
	}

	return &films, nil
}
