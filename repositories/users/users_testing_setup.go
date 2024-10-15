package repositories_users

import (
	"backend/supabase"
	"log"

	"github.com/joho/godotenv"
)

func setupSupabase() {
	// 環境変数の読み込み
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Println("No ../../.env.test file found")
	}

	// テストの前にSupabaseクライアントの初期化
	err = supabase.InitSupabase()
	if err != nil {
		log.Fatalf("Supabase initialization failed: %v", err)
	}
}
