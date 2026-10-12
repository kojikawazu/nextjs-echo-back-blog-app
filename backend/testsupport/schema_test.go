//go:build integration || e2e

package testsupport

import (
	"backend/supabase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// queryStrings は 1 列の文字列を返すクエリを実行し、結果を返す。
func queryStrings(t *testing.T, query string) []string {
	t.Helper()
	rows, err := supabase.Pool.Query(supabase.Ctx, query)
	require.NoError(t, err)
	defer rows.Close()

	var got []string
	for rows.Next() {
		var s string
		require.NoError(t, rows.Scan(&s))
		got = append(got, s)
	}
	require.NoError(t, rows.Err())
	return got
}

// TestSchema_UpdatedAtTriggers は updated_at 列を持つ全テーブルに、
// set_updated_at() を呼ぶ BEFORE UPDATE（行単位）のトリガがあることを検証する（#117）。
// テーブルを追加したときのトリガの付け忘れを検知する。
func TestSchema_UpdatedAtTriggers(t *testing.T) {
	// 正常系: updated_at を持つテーブルは想定どおり 3 つ（検査が空振りしていないことの確認）
	tables := queryStrings(t, `
		SELECT table_name::text FROM information_schema.columns
		WHERE table_schema = 'public' AND column_name = 'updated_at'
		ORDER BY table_name`)
	assert.Equal(t, []string{"blog_likes", "blog_users", "blogs"}, tables)

	// 準正常系: トリガが無い（または BEFORE UPDATE 行単位でない）テーブルが無いこと。
	// tgtype のビット: 1 = 行単位（FOR EACH ROW）、2 = BEFORE、16 = UPDATE
	missing := queryStrings(t, `
		SELECT c.table_name::text FROM information_schema.columns c
		WHERE c.table_schema = 'public' AND c.column_name = 'updated_at'
		  AND NOT EXISTS (
		    SELECT 1 FROM pg_trigger tg
		    JOIN pg_class t ON t.oid = tg.tgrelid
		    JOIN pg_namespace n ON n.oid = t.relnamespace
		    JOIN pg_proc p ON p.oid = tg.tgfoid
		    WHERE n.nspname = 'public' AND t.relname = c.table_name
		      AND p.proname = 'set_updated_at' AND NOT tg.tgisinternal
		      AND tg.tgtype & 1 = 1 AND tg.tgtype & 2 = 2 AND tg.tgtype & 16 = 16
		  )
		ORDER BY c.table_name`)
	assert.Empty(t, missing, "set_updated_at の BEFORE UPDATE トリガが無いテーブル")
}

// TestSchema_AuditColumns は監査列（created_at / updated_at）が
// すべて timestamptz・NOT NULL・DEFAULT now() であることを検証する（#117。準正常系）。
func TestSchema_AuditColumns(t *testing.T) {
	invalid := queryStrings(t, `
		SELECT table_name || '.' || column_name || ' (' || data_type || ', nullable=' || is_nullable
		       || ', default=' || coalesce(column_default, '') || ')'
		FROM information_schema.columns
		WHERE table_schema = 'public' AND column_name IN ('created_at', 'updated_at')
		  AND (data_type <> 'timestamp with time zone' OR is_nullable <> 'NO'
		       OR column_default IS DISTINCT FROM 'now()')
		ORDER BY 1`)
	assert.Empty(t, invalid, "timestamptz NOT NULL DEFAULT now() になっていない監査列")
}

// TestSchema_UpdatedAtOverridesExplicitValue は UPDATE で updated_at を明示的に指定しても、
// トリガが現在時刻で上書きすることを検証する（準正常系。#117）。
// アプリから UPDATE する経路の無い blog_likes で、DB 側で自動設定されることを確認する。
func TestSchema_UpdatedAtOverridesExplicitValue(t *testing.T) {
	var id string
	var createdAt time.Time
	err := supabase.Pool.QueryRow(supabase.Ctx, `
		INSERT INTO blog_likes (blog_id, visit_id)
		VALUES (gen_random_uuid(), gen_random_uuid())
		RETURNING id, created_at`).Scan(&id, &createdAt)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := supabase.Pool.Exec(supabase.Ctx, `DELETE FROM blog_likes WHERE id = $1`, id)
		assert.NoError(t, err)
	})

	before := time.Now()
	var gotCreatedAt, gotUpdatedAt time.Time
	err = supabase.Pool.QueryRow(supabase.Ctx, `
		UPDATE blog_likes SET updated_at = '2000-01-01T00:00:00Z'
		WHERE id = $1
		RETURNING created_at, updated_at`, id).Scan(&gotCreatedAt, &gotUpdatedAt)
	require.NoError(t, err)

	assert.WithinDuration(t, before, gotUpdatedAt, time.Minute, "明示した古い値ではなく現在時刻になること")
	assert.True(t, gotCreatedAt.Equal(createdAt), "created_at は変わらないこと")
}
