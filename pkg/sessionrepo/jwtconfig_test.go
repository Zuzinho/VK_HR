package sessionrepo

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTConfig(t *testing.T) {
	tests := []struct {
		name         string
		methodName   string
		expectedAlgo jwt.SigningMethod
	}{
		{
			name:         "HS256",
			methodName:   "HS256",
			expectedAlgo: jwt.SigningMethodHS256,
		},
		{
			name:         "HS384",
			methodName:   "HS384",
			expectedAlgo: jwt.SigningMethodHS384,
		},
		{
			name:         "HS512",
			methodName:   "HS512",
			expectedAlgo: jwt.SigningMethodHS512,
		},
		// Добавьте дополнительные тесты для других методов подписи по необходимости
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenSecret := []byte("test_secret")
			config := NewJWTConfig(tt.methodName, tokenSecret)

			assert.NotNil(t, config)
			assert.Equal(t, tt.expectedAlgo, config.Method)
			assert.Equal(t, tokenSecret, config.TokenSecret)
		})
	}

	// Тестирование с неподдерживаемым методом подписи
	t.Run("Unsupported", func(t *testing.T) {
		config := NewJWTConfig("Unsupported", []byte("test_secret"))

		assert.Nil(t, config.Method) // Предполагается, что неподдерживаемый метод возвращает nil
	})
}
