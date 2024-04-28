package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoUrl    string
	Port        string
	SwaggerPath string

	Secret          string
	AesKey          string
	AccessTokenTTL  string
	RefreshTokenTTL string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load(".env")
	cfg := Config{
		MongoUrl:    os.Getenv("MONGO_URL"),
		Port:        os.Getenv("PORT"),
		SwaggerPath: os.Getenv("SWAGGER_PATH"),

		Secret:          os.Getenv("SECRET"),
		AesKey:          os.Getenv("AES_KEY"),
		AccessTokenTTL:  os.Getenv("ACCESS_TOKEN_TTL"),
		RefreshTokenTTL: os.Getenv("REFRESH_TOKEN_TTL"),
	}
	if err != nil {
		return cfg, err
	}

	return cfg, err
}
