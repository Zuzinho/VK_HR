package actorrepo

import (
	"VK_HR/pkg/customtime"
	"VK_HR/pkg/filmrepo"
	"context"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
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
		"($1, $2, $3, $4) returning actor_id", actor.FirstName, actor.SecondName, actor.Gender,
		actor.Birthday.Format("2006-01-02")).Scan(&insertedID)

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
		"FROM actors a LEFT JOIN actors_has_films ahf ON a.actor_id = ahf.actor_id "+
		"LEFT JOIN films f ON ahf.film_id = f.film_id order by a.actor_id")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	actors := make(Actors, 0)
	prevActor := &Actor{}
	for rows.Next() {
		actor := &Actor{}

		var filmID sql.NullInt32
		var name, description sql.NullString
		var premierDate sql.NullTime
		var rating sql.NullFloat64

		var birthday sql.NullTime

		err = rows.Scan(&actor.ActorID, &actor.FirstName, &actor.SecondName, &actor.Gender, &birthday,
			&filmID, &name, &description, &premierDate, &rating)
		if err != nil {
			log.Errorf("Skipped err when getting actors or films: %s", err.Error())
			continue
		}
		actor.Birthday = customtime.CustomTime{
			Time: birthday.Time,
		}

		if prevActor.ActorID != actor.ActorID {
			if prevActor.ActorID > 0 {
				actors.Append(prevActor)
			}

			films := make(filmrepo.Films, 0)
			actor.Films = &films

			prevActor = actor
		}

		film := &filmrepo.Film{
			FilmID:      filmID.Int32,
			Name:        name.String,
			Description: description.String,
			PremierDate: customtime.CustomTime{Time: premierDate.Time},
			Rating:      float32(rating.Float64),
		}

		prevActor.Films.Append(film)
	}

	actors.Append(prevActor)

	return &actors, rows.Err()
}
