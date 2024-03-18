package filmrepo

import (
	"context"
	"time"
)

type Film struct {
	FilmID      int       `json:"film_id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PremierDate time.Time `json:"premier_date"`
	Rating      float32   `json:"rating"`
	ActorsID    []int     `json:"actors_id,omitempty"`
}

type Films []*Film

func (films *Films) Append(film *Film) {
	*films = append(*films, film)
}

type FilmsRepository interface {
	InsertFilm(ctx context.Context, film *Film) (int, error)
	InsertActorsByFilm(ctx context.Context, filmID int, actorsID []int) error
	Update(ctx context.Context, setClause *string, args *[]any, id int) error
	Delete(ctx context.Context, id int) error
	SelectBySorting(ctx context.Context, config *SortingConfig) (*Films, error)
	SelectByFragment(ctx context.Context, filmFragment, nameFragment string) (*Films, error)
}
