package config

import (
	"log"
	"os"
)

// 環境変数から読み込む
var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var IsProduction = os.Getenv("ENV") == "production"

func init() {
	if len(JwtKey) == 0 {
		log.Fatal("JWT_SECRET_KEY is not set in the environment")
	}
}
