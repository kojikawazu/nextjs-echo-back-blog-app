package supabase

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildConnString は接続URLに sslmode を付与した接続文字列を組み立てることを検証する。
func TestBuildConnString(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		sslmode string
		want    string
	}{
		{
			// 正常系: sslmode を指定した場合はその値を付与する
			name:    "with_sslmode",
			baseURL: "postgres://user:pass@localhost:5432/db",
			sslmode: "disable",
			want:    "postgres://user:pass@localhost:5432/db?sslmode=disable",
		},
		{
			// 準正常系: sslmode が空の場合は本番デフォルトの require を付与する
			name:    "empty_sslmode_defaults_to_require",
			baseURL: "postgres://user:pass@localhost:5432/db",
			sslmode: "",
			want:    "postgres://user:pass@localhost:5432/db?sslmode=require",
		},
		{
			// 準正常系: 既にクエリ文字列を持つ URL は & で連結する
			name:    "existing_query_uses_ampersand",
			baseURL: "postgres://user:pass@localhost:5432/db?application_name=blog",
			sslmode: "disable",
			want:    "postgres://user:pass@localhost:5432/db?application_name=blog&sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, buildConnString(tt.baseURL, tt.sslmode))
		})
	}
}

// TestNewPoolConfig はコネクションプールの設定が組み立てられることを検証する（正常系）。
// Simple Protocol の指定は pooler の無い IT / E2E では外れても検出できないため、ここで固定する。
func TestNewPoolConfig(t *testing.T) {
	config, err := newPoolConfig("postgres://user:pass@localhost:5432/db?sslmode=disable")

	require.NoError(t, err)
	assert.Equal(t, pgx.QueryExecModeSimpleProtocol, config.ConnConfig.DefaultQueryExecMode)
	assert.Equal(t, int32(10), config.MaxConns)
	assert.Equal(t, 30*time.Second, config.MaxConnIdleTime)
	assert.Equal(t, "localhost", config.ConnConfig.Host)
}

// TestNewPoolConfig_Invalid は解析できない接続文字列でエラーを返すことを検証する。
func TestNewPoolConfig_Invalid(t *testing.T) {
	tests := []struct {
		name       string
		connString string
	}{
		{
			// 準正常系: URL として不正
			name:       "malformed_url",
			connString: "postgres://user:pass@localhost:notaport/db",
		},
		{
			// 準正常系: 未知の sslmode
			name:       "unknown_sslmode",
			connString: "postgres://user:pass@localhost:5432/db?sslmode=bogus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := newPoolConfig(tt.connString)

			assert.Error(t, err)
			assert.Nil(t, config)
		})
	}
}
