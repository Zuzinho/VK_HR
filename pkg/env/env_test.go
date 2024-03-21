package env

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitEnv(t *testing.T) {
	assert.NoError(t, godotenv.Load("../env/.env.test"))
}

func TestMustDBConnString(t *testing.T) {
	expected := "test_connection_string"
	os.Setenv("DB_CONNECTION_STRING", expected)
	defer os.Unsetenv("DB_CONNECTION_STRING")

	result := MustDBConnString()
	assert.Equal(t, expected, result)
}

func TestMustPort(t *testing.T) {
	expected := "8080"
	os.Setenv("PORT", expected)
	defer os.Unsetenv("PORT")

	result := MustPort()
	assert.Equal(t, expected, result)
}

func TestMustJWTConfig(t *testing.T) {
	os.Setenv("JWT_METHOD_NAME", "HS256")
	os.Setenv("JWT_TOKEN_SECRET", "secret")
	defer func() {
		os.Unsetenv("JWT_METHOD_NAME")
		os.Unsetenv("JWT_TOKEN_SECRET")
	}()

	result := MustJWTConfig()
	assert.NotNil(t, result)
	assert.Equal(t, jwt.SigningMethodHS256, result.Method)
}

func TestMustMaxConnCount(t *testing.T) {
	expected := 10
	os.Setenv("MAX_CONNECTION_COUNT", "10")
	defer os.Unsetenv("MAX_CONNECTION_COUNT")

	result := MustMaxConnCount()
	assert.Equal(t, expected, result)
}
