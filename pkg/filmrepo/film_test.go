package filmrepo

import (
	"VK_HR/pkg/customtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFilmsAppend(t *testing.T) {
	films := Films{}

	filmToAdd := &Film{
		FilmID:      1,
		Name:        "Test Film",
		Description: "Test Description",
		PremierDate: customtime.CustomTime{Time: time.Now()},
		Rating:      8.0,
		ActorsID:    []int{1, 2, 3},
	}

	films.Append(filmToAdd)

	// Проверяем, что фильм был добавлен
	assert.Equal(t, 1, len(films), "Expected films length to be 1 after appending")

	// Проверяем, что добавленный фильм соответствует ожидаемому
	assert.Equal(t, filmToAdd, films[0], "The appended film does not match the expected one")
}
