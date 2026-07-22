package config

import (
	"log"
	"os"
)

// JwtKey は環境変数 JWT_SECRET_KEY から読み込む JWT 署名鍵。
var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// IsProduction は環境変数 ENV が "production" のとき true となる本番判定フラグ。
var IsProduction = os.Getenv("ENV") == "production"

func init() {
	if len(JwtKey) == 0 {
		log.Fatal("JWT_SECRET_KEY is not set in the environment")
	}
}
