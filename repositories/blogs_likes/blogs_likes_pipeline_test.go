package repositories_blogs_likes

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository_BlogLike_Test(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := NewBlogLikeRepository()

	// UUIDを生成
	blogID := os.Getenv("TEST_BLOG_ID")
	visitorID := uuid.New().String()

	// ---------------------------------------------------------
	// 1. 「いいね」を作成
	// ---------------------------------------------------------
	like, err := repo.CreateBlogLike(blogID, visitorID)
	if err != nil {
		t.Fatalf("Failed to create blog like: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, like)

	// ---------------------------------------------------------
	// 2. 「いいね」が存在するか確認
	// ---------------------------------------------------------
	liked, err := repo.IsBlogLiked(blogID, visitorID)
	if err != nil {
		t.Fatalf("Failed to check if blog is liked: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.True(t, liked)

	// ---------------------------------------------------------
	// 3. 「いいね」を削除
	// ---------------------------------------------------------
	err = repo.DeleteBlogLike(blogID, visitorID)
	if err != nil {
		t.Fatalf("Failed to delete blog like: %v", err)
	}

	// エラーチェック
	assert.NoError(t, err)

	// ---------------------------------------------------------
	// 4. 「いいね」が削除されたことを確認
	// ---------------------------------------------------------
	liked, err = repo.IsBlogLiked(blogID, visitorID)

	// エラーチェックとデータ確認（削除後なので false を期待）
	assert.Error(t, err)
	assert.False(t, liked)
}
