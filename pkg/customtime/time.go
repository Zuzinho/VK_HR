package customtime

import "time"

type CustomTime struct {
	time.Time
}

func (date *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)

	// Убираем кавычки вокруг JSON строки
	s = s[1 : len(s)-1]

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	date.Time = t
	return nil
}

func (date *CustomTime) MarshalJSON() ([]byte, error) {
	str := date.Format("2006-01-02")

	return []byte("\"" + str + "\""), nil
}
