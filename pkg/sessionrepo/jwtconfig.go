package sessionrepo

import "github.com/dgrijalva/jwt-go"

type JWTConfig struct {
	Method      jwt.SigningMethod
	TokenSecret []byte
}

func NewJWTConfig(methodName string, tokenSecret []byte) *JWTConfig {
	return &JWTConfig{
		Method:      jwt.GetSigningMethod(methodName),
		TokenSecret: tokenSecret,
	}
}
