package sessionrepo

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionManager_Pack(t *testing.T) {
	config := NewJWTConfig("HS256", []byte("secret"))
	manager := NewSessionManager(config)

	tokenString, err := manager.Pack("user1")
	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.ParseWithClaims(*tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(*CustomClaims)
	require.True(t, ok)
	assert.Equal(t, "user1", claims.Login)
	assert.True(t, claims.ExpiresAt > time.Now().Unix())
}

func TestSessionManager_Unpack(t *testing.T) {
	config := NewJWTConfig("HS256", []byte("secret"))
	manager := NewSessionManager(config)

	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &CustomClaims{
		Login: "user1",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	session, err := manager.Unpack(tokenString)
	require.NoError(t, err)
	assert.Equal(t, "user1", session.Sub)
	assert.WithinDuration(t, expirationTime, session.Exp, time.Second)
}
