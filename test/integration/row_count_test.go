package integration_test

import (
	"testing"
)

func TestUserSettingsUpsert_RowCount(t *testing.T) {
	db := setupDB(t)

	// raw SQLを直接実行してRowsAffectedを取得することで、MySQLのROW_COUNT()の挙動を検証する。
	// ON DUPLICATE KEY UPDATEにおけるROW_COUNT()の仕様:
	//   INSERT実行時（新規挿入）: 1
	//   UPDATE実行時（値が変わった場合）: 2
	//   変更なし（同じ値でのUpsert）: 0
	upsertSQL := `
		INSERT INTO user_settings (user_id, theme, language, updated_at)
		VALUES (?, ?, ?, NOW()) AS new
		ON DUPLICATE KEY UPDATE
		  theme      = new.theme,
		  language   = new.language,
		  updated_at = NOW()
	`

	t.Run("INSERT実行時_ROW_COUNTが1を返す", func(t *testing.T) {
		cleanTable(t, db, "user_settings")

		// 存在しないユーザーIDでINSERTを実行する
		result, err := db.Exec(upsertSQL, 10, "dark", "ja")
		if err != nil {
			t.Fatalf("Execに失敗しました: %v", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			t.Fatalf("RowsAffectedの取得に失敗しました: %v", err)
		}

		// 新規INSERTなのでROW_COUNTは1になることを確認する
		if rows != 1 {
			t.Errorf("RowsAffected: got %d, want 1", rows)
		}
	})

	t.Run("UPDATE実行時_ROW_COUNTが2を返す", func(t *testing.T) {
		cleanTable(t, db, "user_settings")

		// 事前に初期データを挿入しておく
		_, err := db.Exec(upsertSQL, 11, "light", "en")
		if err != nil {
			t.Fatalf("初回Execに失敗しました: %v", err)
		}

		// 別の値でUpsertを実行する（ON DUPLICATE KEY UPDATEが発動する）
		result, err := db.Exec(upsertSQL, 11, "dark", "ja")
		if err != nil {
			t.Fatalf("2回目のExecに失敗しました: %v", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			t.Fatalf("RowsAffectedの取得に失敗しました: %v", err)
		}

		// 値が変わってUPDATEされたのでROW_COUNTは2になることを確認する
		if rows != 2 {
			t.Errorf("RowsAffected: got %d, want 2", rows)
		}
	})

	t.Run("変更なし_ROW_COUNTが0を返す", func(t *testing.T) {
		cleanTable(t, db, "user_settings")

		// ROW_COUNT=0を確実に再現するため、固定のupdated_atを使うSQLで検証する。
		// NOW()を使うとUPDATEの度に時刻が変わってしまいROW_COUNTが2になるため、ここでは固定値を使う。
		fixedUpsertSQL := `
			INSERT INTO user_settings (user_id, theme, language, updated_at)
			VALUES (?, ?, ?, '2024-01-01 00:00:00') AS new
			ON DUPLICATE KEY UPDATE
			  theme      = new.theme,
			  language   = new.language,
			  updated_at = new.updated_at
		`

		// 初回は固定値で挿入する
		_, err := db.Exec(fixedUpsertSQL, 12, "dark", "ja")
		if err != nil {
			t.Fatalf("初回Execに失敗しました: %v", err)
		}

		// 完全に同じ値で再度Upsertを実行する
		result, err := db.Exec(fixedUpsertSQL, 12, "dark", "ja")
		if err != nil {
			t.Fatalf("2回目のExecに失敗しました: %v", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			t.Fatalf("RowsAffectedの取得に失敗しました: %v", err)
		}

		// MySQLはUPDATEが発動しても値が変わらない場合、ROW_COUNTを0にする
		if rows != 0 {
			t.Errorf("RowsAffected: got %d, want 0", rows)
		}
	})
}
