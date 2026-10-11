//go:build integration || e2e

package testsupport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCheckLocalDB は接続先がローカルの DB のみ許可されることを検証する。
func TestCheckLocalDB(t *testing.T) {
	tests := []struct {
		name       string
		connString string
		wantErr    error
	}{
		// 正常系: ローカルの DB（ローカル Supabase・コンテナ等）は許可する
		{name: "localhost_url", connString: "postgresql://postgres:postgres@localhost:54322/postgres", wantErr: nil},
		{name: "localhost_uppercase", connString: "postgresql://postgres:postgres@LOCALHOST:5432/postgres", wantErr: nil},
		{name: "ipv4_loopback", connString: "postgresql://postgres:postgres@127.0.0.1:54322/postgres", wantErr: nil},
		{name: "ipv4_loopback_range", connString: "postgresql://postgres:postgres@127.0.0.2:5432/postgres", wantErr: nil},
		{name: "ipv6_loopback", connString: "postgresql://postgres:postgres@[::1]:5432/postgres", wantErr: nil},
		{name: "keyword_value_dsn", connString: "host=localhost port=5432 user=postgres password=postgres dbname=testdb", wantErr: nil},
		{name: "unix_socket", connString: "host=/var/run/postgresql user=postgres dbname=testdb", wantErr: nil},
		{name: "multiple_local_hosts", connString: "postgresql://postgres:postgres@localhost:5432,127.0.0.1:5433/postgres", wantErr: nil},

		// 準正常系: リモートの DB（本番 Supabase 等）は拒否する
		{name: "supabase_direct", connString: "postgresql://postgres:secret-pw@db.abcdefghijklmnop.supabase.co:5432/postgres", wantErr: errRemoteDB},
		{name: "supabase_pooler", connString: "postgresql://postgres.abcdefghijklmnop:secret-pw@aws-0-ap-northeast-1.pooler.supabase.com:6543/postgres", wantErr: errRemoteDB},
		{name: "private_ip", connString: "postgresql://postgres:secret-pw@10.0.0.5:5432/postgres", wantErr: errRemoteDB},
		{name: "localhost_lookalike", connString: "postgresql://postgres:secret-pw@localhost.example.com:5432/postgres", wantErr: errRemoteDB},
		{name: "remote_keyword_value_dsn", connString: "host=db.example.com user=postgres password=secret-pw dbname=postgres", wantErr: errRemoteDB},
		{name: "remote_in_fallback_hosts", connString: "postgresql://postgres:secret-pw@localhost:5432,db.example.com:5432/postgres", wantErr: errRemoteDB},

		// 異常系: 解析できない接続文字列は拒否する
		{name: "unparsable", connString: "postgresql://postgres:secret-pw@%zz:5432/postgres", wantErr: errUnparsableDB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkLocalDB(tt.connString)

			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				// 公開される CI ログに接続情報を残さない
				assert.NotContains(t, err.Error(), "secret-pw")
				assert.NotContains(t, err.Error(), "supabase.co")
				assert.NotContains(t, err.Error(), "example.com")
			}
		})
	}
}

// TestSetup_RejectsRemoteDB は SUPABASE_URL がリモートを指すとき、接続せずに Setup がエラーを返すことを検証する（準正常系）。
func TestSetup_RejectsRemoteDB(t *testing.T) {
	t.Setenv("SUPABASE_URL", "postgresql://postgres:secret-pw@db.abcdefghijklmnop.supabase.co:5432/postgres")

	teardown, err := Setup()

	assert.ErrorIs(t, err, errRemoteDB)
	assert.Nil(t, teardown)
}
