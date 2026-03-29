package integration_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// setupDB はテスト用のDB接続を返す。
// 環境変数 TEST_DSN が設定されていればその値を使い、未設定の場合はデフォルトDSNを使う。
func setupDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3308)/upsert_test?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("DB接続のオープンに失敗しました: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("DBへの接続確認に失敗しました: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// cleanTable は指定したテーブルのレコードをすべて削除する。
// テスト間の状態汚染を防ぐために各テストの冒頭で呼ぶ。
func cleanTable(t *testing.T, db *sql.DB, tables ...string) {
	t.Helper()

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Fatalf("テーブル %s のクリーンアップに失敗しました: %v", table, err)
		}
	}
}
