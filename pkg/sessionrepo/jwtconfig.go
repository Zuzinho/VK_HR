package sessionrepo

import "github.com/golang-jwt/jwt"

type JWTConfig struct {
	Method      *jwt.SigningMethodHMAC
	TokenSecret []byte
}

func NewJWTConfig(methodName string, tokenSecret []byte) *JWTConfig {
	return &JWTConfig{
		Method: &jwt.SigningMethodHMAC{
			Name: methodName,
		},
		TokenSecret: tokenSecret,
	}
}
