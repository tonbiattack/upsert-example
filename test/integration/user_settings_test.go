package integration_test

import (
	"testing"

	"github.com/example/upsert-example/internal/repository"
)

func TestUserSettingsRepository_Upsert(t *testing.T) {
	db := setupDB(t)
	repo := repository.NewUserSettingsRepository(db)

	t.Run("新規挿入_存在しないユーザー", func(t *testing.T) {
		cleanTable(t, db, "user_settings")

		// 存在しないユーザーIDでUpsertを実行する
		err := repo.Upsert(1, "dark", "ja")
		if err != nil {
			t.Fatalf("Upsertに失敗しました: %v", err)
		}

		// レコードが作成されていることを確認する
		got, err := repo.FindByUserID(1)
		if err != nil {
			t.Fatalf("FindByUserIDに失敗しました: %v", err)
		}
		if got == nil {
			t.Fatal("レコードが存在しません")
		}
		if got.Theme != "dark" {
			t.Errorf("Theme: got %q, want %q", got.Theme, "dark")
		}
		if got.Language != "ja" {
			t.Errorf("Language: got %q, want %q", got.Language, "ja")
		}
	})

	t.Run("更新_既存ユーザーの設定が上書きされる", func(t *testing.T) {
		cleanTable(t, db, "user_settings")

		// 1回目のUpsertで初期データを作成する
		if err := repo.Upsert(2, "light", "en"); err != nil {
			t.Fatalf("1回目のUpsertに失敗しました: %v", err)
		}

		// 2回目のUpsertで別の値に更新する
		if err := repo.Upsert(2, "dark", "ja"); err != nil {
			t.Fatalf("2回目のUpsertに失敗しました: %v", err)
		}

		// Themeが更新されていることを確認する
		got, err := repo.FindByUserID(2)
		if err != nil {
			t.Fatalf("FindByUserIDに失敗しました: %v", err)
		}
		if got == nil {
			t.Fatal("レコードが存在しません")
		}
		if got.Theme != "dark" {
			t.Errorf("Theme: got %q, want %q", got.Theme, "dark")
		}
		if got.Language != "ja" {
			t.Errorf("Language: got %q, want %q", got.Language, "ja")
		}
	})

	t.Run("更新_同じ値でもエラーにならない", func(t *testing.T) {
		cleanTable(t, db, "user_settings")

		// 1回目のUpsertで初期データを作成する
		if err := repo.Upsert(3, "dark", "ja"); err != nil {
			t.Fatalf("1回目のUpsertに失敗しました: %v", err)
		}

		// 同じ値で再度Upsertを呼んでもエラーにならないことを確認する
		if err := repo.Upsert(3, "dark", "ja"); err != nil {
			t.Fatalf("同じ値でのUpsertがエラーになりました: %v", err)
		}
	})
}
