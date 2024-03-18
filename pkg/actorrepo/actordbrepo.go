package actorrepo

import (
	"VK_HR/pkg/filmrepo"
	"context"
	"database/sql"
	"fmt"
)

type ActorsDBRepository struct {
	DB *sql.DB
}

func NewActorsDBRepository(db *sql.DB) *ActorsDBRepository {
	return &ActorsDBRepository{
		DB: db,
	}
}

func (repo *ActorsDBRepository) Insert(ctx context.Context, actor *Actor) (int, error) {
	var insertedID int

	err := repo.DB.QueryRowContext(ctx, "insert into actors (first_name, second_name, gender, birthday) values "+
		"($1, $2, $3, $4) returning actor_id", actor.FirstName, actor.SecondName, actor.Gender, actor.Birthday).Scan(&insertedID)

	return insertedID, err
}

func (repo *ActorsDBRepository) Update(ctx context.Context, setClause *string, args *[]any, id int) error {
	queryString := fmt.Sprintf("update actors set %s where actor_id = $%d", *setClause, len(*args)+1)

	smtm, err := repo.DB.PrepareContext(ctx, queryString)
	defer smtm.Close()
	if err != nil {
		return err
	}

	_, err = smtm.ExecContext(ctx, append(*args, id)...)
	return err
}

func (repo *ActorsDBRepository) Delete(ctx context.Context, id int) error {
	_, err := repo.DB.ExecContext(ctx, "delete from actors where actor_id = $1", id)

	return err
}

func (repo *ActorsDBRepository) SelectAll(ctx context.Context) (*Actors, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT a.actor_id, a.first_name, a.second_name, a.gender, a.birthday, "+
		"f.film_id, f.name, f.description, f.premier_date, f.rating "+
		"FROM actors a JOIN actors_has_films ahf ON a.actor_id = ahf.actor_id "+
		"JOIN films f ON ahf.film_id = f.film_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actors := make(Actors, 0)
	prevActor := &Actor{}
	for rows.Next() {
		actor := &Actor{}

		err = rows.Scan(&actor.ActorID, &actor.FirstName, &actor.SecondName, &actor.Gender, &actor.Birthday)
		if err != nil {
			continue
		}

		if prevActor.ActorID != actor.ActorID {
			actors.Append(prevActor)

			prevActor = actor
		}

		film := &filmrepo.Film{}

		err = rows.Scan(&film.FilmID, &film.Name, &film.Description, &film.PremierDate, &film.Rating)
		if err != nil {
			continue
		}

		prevActor.Films.Append(film)
	}

	return &actors, rows.Err()
}
