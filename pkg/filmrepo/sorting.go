package filmrepo

import "VK_HR/pkg/validator"

type SortingDirection string

const (
	ASG  SortingDirection = "ASG"
	DESG SortingDirection = "DESG"
)

type SortingConfig struct {
	ColumnName validator.ColumnName
	Direction  SortingDirection
}

func NewSortingConfig(name validator.ColumnName, direction SortingDirection) (*SortingConfig, error) {
	if name == "" {
		name = "rating"
	} else {
		switch name {
		case validator.FilmName, validator.FilmRating, validator.FilmPremierDate:
		default:
			return nil, newNoAccessSortingColumnNameError(name)
		}
	}

	if direction == "" {
		direction = DESG
	} else {
		if direction != ASG && direction != DESG {
			return nil, newUnknownSortingDirectionError(direction)
		}
	}

	return &SortingConfig{
		ColumnName: name,
		Direction:  direction,
	}, nil
}
