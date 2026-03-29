package integration_test

import (
	"testing"

	"github.com/example/upsert-example/internal/repository"
)

func TestUserTagRepository_InsertIgnore(t *testing.T) {
	db := setupDB(t)
	repo := repository.NewUserTagRepository(db)

	t.Run("存在しないタグの追加_1件挿入される", func(t *testing.T) {
		cleanTable(t, db, "user_tags")

		// 存在しないuser_id・tag_idの組み合わせでINSERT IGNOREを実行する
		rows, err := repo.InsertIgnore(1, 10)
		if err != nil {
			t.Fatalf("InsertIgnoreに失敗しました: %v", err)
		}

		// 新規挿入時はRowsAffectedが1になることを確認する
		if rows != 1 {
			t.Errorf("RowsAffected: got %d, want 1", rows)
		}

		// レコードが存在することを確認する
		exists, err := repo.Exists(1, 10)
		if err != nil {
			t.Fatalf("Existsに失敗しました: %v", err)
		}
		if !exists {
			t.Error("レコードが存在しません")
		}
	})

	t.Run("既存タグのINSERT_IGNORE_スキップされる", func(t *testing.T) {
		cleanTable(t, db, "user_tags")

		// 先にレコードを挿入しておく
		if _, err := repo.InsertIgnore(2, 20); err != nil {
			t.Fatalf("初回InsertIgnoreに失敗しました: %v", err)
		}

		// 同じ主キーで再度INSERT IGNOREを実行する
		rows, err := repo.InsertIgnore(2, 20)
		if err != nil {
			// INSERT IGNOREはエラーを返さずスキップするため、ここには来ないはず
			t.Fatalf("INSERT IGNOREがエラーを返しました: %v", err)
		}

		// スキップ時はRowsAffectedが0になることを確認する
		if rows != 0 {
			t.Errorf("RowsAffected: got %d, want 0", rows)
		}
	})
}
