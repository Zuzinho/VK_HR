package actorrepo

import (
	"VK_HR/pkg/filmrepo"
	"VK_HR/pkg/gender"
	"context"
	"time"
)

type Actor struct {
	ActorID    int             `json:"actor_id,omitempty"`
	FirstName  string          `json:"first_name"`
	SecondName string          `json:"second_name"`
	Gender     gender.Gender   `json:"gender"`
	Birthday   time.Time       `json:"birthday"`
	Films      *filmrepo.Films `json:"films,omitempty"`
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
