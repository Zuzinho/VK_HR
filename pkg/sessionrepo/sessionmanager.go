package sessionrepo

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type CustomClaims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

type SessionManager struct {
	Config *JWTConfig
}

func NewSessionManager(config *JWTConfig) *SessionManager {
	return &SessionManager{
		Config: config,
	}
}

func (manager *SessionManager) Pack(sub string) (*string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &CustomClaims{
		Login: sub,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(manager.Config.Method, claims)

	tokenString, err := token.SignedString(manager.Config.TokenSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (manager *SessionManager) Unpack(inToken string) (*Session, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(inToken, claims, func(token *jwt.Token) (interface{}, error) {
		return manager.Config.TokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return &Session{
			Sub: claims.Login,
			Exp: time.Unix(claims.ExpiresAt, 0),
		}, nil
	} else {
		return nil, newInvalidTokenError(inToken)
	}
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
