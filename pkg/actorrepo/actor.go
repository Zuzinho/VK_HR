package actorrepo

import (
	"VK_HR/pkg/customtime"
	"VK_HR/pkg/filmrepo"
	"VK_HR/pkg/gender"
	"context"
)

type Actor struct {
	ActorID    int32                 `json:"actor_id,omitempty"`
	FirstName  string                `json:"first_name" validate:"required,min=1"`
	SecondName string                `json:"second_name" validate:"required,min=1"`
	Gender     gender.Gender         `json:"gender" validate:"oneof=Male Female"`
	Birthday   customtime.CustomTime `json:"birthday"`
	Films      *filmrepo.Films       `json:"films,omitempty"`
}

type Actors []*Actor

func (actors *Actors) Append(actor *Actor) {
	*actors = append(*actors, actor)
}

type ActorsRepository interface {
	Insert(ctx context.Context, actor *Actor) (int, error)
	Update(ctx context.Context, setClause *string, args *[]any, id int) error
	Delete(ctx context.Context, id int) error
	SelectAll(ctx context.Context) (*Actors, error)
}
