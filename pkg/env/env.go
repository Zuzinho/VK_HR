package env

import (
	"VK_HR/pkg/sessionrepo"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func InitEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("No .env file")
	}
}

func MustDBConnString() string {
	val, exist := os.LookupEnv("DB_CONNECTION_STRING")
	if !exist {
		log.Fatal("No db connection string")
	}

	return val
}

func MustPort() string {
	val, exist := os.LookupEnv("PORT")
	if !exist {
		log.Fatal("No db connection string")
	}

	return val
}

func MustJWTConfig() *sessionrepo.JWTConfig {
	methodName, exist := os.LookupEnv("JWT_METHOD_NAME")
	if !exist {
		log.Fatal("No jwt method name")
	}

	tokenSecret, exist := os.LookupEnv("JWT_TOKEN_SECRET")
	if !exist {
		log.Fatal("No jwt secret token")
	}

	return sessionrepo.NewJWTConfig(methodName, []byte(tokenSecret))
}
