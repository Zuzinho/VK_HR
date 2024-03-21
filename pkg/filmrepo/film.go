package filmrepo

import (
	"VK_HR/pkg/customtime"
	"context"
)

type Film struct {
	FilmID      int32                 `json:"film_id,omitempty"`
	Name        string                `json:"name" validate:"required,min=1,max=150"`
	Description string                `json:"description" validate:"max=1000"`
	PremierDate customtime.CustomTime `json:"premier_date"`
	Rating      float32               `json:"rating"  validate:"required,gte=0,lte=10"`
	ActorsID    []int                 `json:"actors_id,omitempty"`
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
