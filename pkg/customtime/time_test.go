package customtime

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnmarshalJSON(t *testing.T) {
	jsonStr := `"2023-01-01"`
	var ct CustomTime

	err := json.Unmarshal([]byte(jsonStr), &ct)
	assert.NoError(t, err)

	expectedTime, _ := time.Parse("2006-01-02", "2023-01-01")
	assert.True(t, expectedTime.Equal(ct.Time), "Unmarshaled time is incorrect")
}

func TestMarshalJSON(t *testing.T) {
	expectedJsonStr := `"2023-01-01"`
	ct := CustomTime{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}

	marshaled, err := ct.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, expectedJsonStr, string(marshaled), "Marshaled JSON is incorrect")
}
