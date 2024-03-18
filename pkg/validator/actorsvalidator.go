package validator

import (
	"VK_HR/pkg/gender"
	"time"
)

const (
	ActorFirstName  ColumnName = "first_name"
	ActorSecondName ColumnName = "second_name"
	ActorGender     ColumnName = "gender"
	ActorBirthday   ColumnName = "birthday"
)

type ActorsValidator struct {
}

func NewActorsValidator() *ActorsValidator {
	return &ActorsValidator{}
}

func (validator *ActorsValidator) IsValidValue(name ColumnName, value string) (any, error) {
	switch name {
	case ActorFirstName, ActorSecondName:
		return value, nil
	case ActorGender:
		gender := gender.Gender(value)
		return gender, gender.IsValid()
	case ActorBirthday:
		return time.Parse("2006-01-02", value)
	default:
		return nil, newInvalidColumnNameError(name)
	}
}
