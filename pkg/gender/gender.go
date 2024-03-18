package gender

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

func (gender Gender) IsValid() error {
	switch gender {
	case Male, Female:
		return nil
	default:
		return newUnknownGenderError(gender)
	}
}
