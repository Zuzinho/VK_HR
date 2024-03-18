package sessionrepo

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type SessionManager struct {
	Config *JWTConfig
}

func NewSessionManager(config *JWTConfig) *SessionManager {
	return &SessionManager{
		Config: config,
	}
}

func (manager *SessionManager) Pack(sub string) (*string, error) {
	token := jwt.NewWithClaims(manager.Config.Method, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(24 * time.Hour),
	})

	tokenString, err := token.SignedString(manager.Config.TokenSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (manager *SessionManager) Unpack(inToken string) (*Session, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method != manager.Config.Method {
			return nil, fmt.Errorf("bad sign method")
		}
		return manager.Config.TokenSecret, nil
	}

	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err != nil || !token.Valid {
		return nil, newInvalidTokenError(inToken)
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, newInvalidTokenMapClaimsError(inToken)
	}

	var sub string
	var exp time.Time

	err = unpackValue[string](&sub, payload, "sub")
	if err != nil {
		return nil, err
	}

	err = unpackValue[time.Time](&exp, payload, "exp")
	if err != nil {
		return nil, err
	}

	return &Session{
		Sub: sub,
		Exp: exp,
	}, nil
}

func unpackValue[V string | time.Time](target any, payload jwt.MapClaims, key string) error {
	val, exist := payload[key]
	if !exist {
		return newNoPayloadKeyError(key)
	}

	target, ok := val.(V)
	if !ok {
		return newInvalidPayloadValueError(val)
	}

	return nil
}
