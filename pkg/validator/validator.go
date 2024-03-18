package validator

type ColumnName string

type ValueValidator interface {
	IsValidValue(name ColumnName, value string) (any, error)
}
