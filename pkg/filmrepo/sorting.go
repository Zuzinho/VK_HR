package filmrepo

import "VK_HR/pkg/validator"

type SortingDirection string

const (
	ASC  SortingDirection = "ASC"
	DESC SortingDirection = "DESC"
)

func (direction SortingDirection) IsValid() error {
	if direction != ASC && direction != DESC {
		return newUnknownSortingDirectionError(direction)
	}

	return nil
}

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
		direction = DESC
	} else {
		if err := direction.IsValid(); err != nil {
			return nil, err
		}
	}

	return &SortingConfig{
		ColumnName: name,
		Direction:  direction,
	}, nil
}
