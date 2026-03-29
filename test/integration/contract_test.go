package integration_test

import (
	"strings"
	"testing"

	"github.com/example/upsert-example/internal/repository"
)

func TestContractRepository_Upsert(t *testing.T) {
	db := setupDB(t)
	repo := repository.NewContractRepository(db)

	t.Run("新規契約が作成される", func(t *testing.T) {
		cleanTable(t, db, "contracts")

		// 新規のcompany_idでUpsertを実行する
		rows, err := repo.Upsert(100, "standard", "2024-01-01")
		if err != nil {
			t.Fatalf("Upsertに失敗しました: %v", err)
		}

		// INSERT時はRowsAffectedが1になることを確認する
		if rows != 1 {
			t.Errorf("RowsAffected: got %d, want 1", rows)
		}

		// レコードが作成されていることを確認する
		got, err := repo.FindByCompanyID(100)
		if err != nil {
			t.Fatalf("FindByCompanyIDに失敗しました: %v", err)
		}
		if got == nil {
			t.Fatal("レコードが存在しません")
		}
		if got.Plan != "standard" {
			t.Errorf("Plan: got %q, want %q", got.Plan, "standard")
		}
	})

	t.Run("既存会社IDでUPSERTすると既存の契約内容が上書きされる", func(t *testing.T) {
		cleanTable(t, db, "contracts")

		// 先に"standard"プランで契約を登録する
		if _, err := repo.Upsert(200, "standard", "2024-01-01"); err != nil {
			t.Fatalf("1回目のUpsertに失敗しました: %v", err)
		}

		// 同じcompany_idで"enterprise"プランに変更するUpsertを実行する
		// これが「意図しない上書き」の検証パターン：
		// signed_atも含めて既存レコードの内容がすべて新しい値で上書きされる
		rows, err := repo.Upsert(200, "enterprise", "2025-06-01")
		if err != nil {
			t.Fatalf("2回目のUpsertに失敗しました: %v", err)
		}

		// UPDATE時はRowsAffectedが2になることを確認する
		if rows != 2 {
			t.Errorf("RowsAffected: got %d, want 2", rows)
		}

		// planとsigned_atが上書きされていることを確認する
		got, err := repo.FindByCompanyID(200)
		if err != nil {
			t.Fatalf("FindByCompanyIDに失敗しました: %v", err)
		}
		if got == nil {
			t.Fatal("レコードが存在しません")
		}
		if got.Plan != "enterprise" {
			t.Errorf("Plan: got %q, want %q", got.Plan, "enterprise")
		}
	})

	t.Run("通常INSERTは重複時にエラーになる", func(t *testing.T) {
		cleanTable(t, db, "contracts")

		// 1回目のInsertで初期データを作成する
		if err := repo.Insert(300, "standard", "2024-01-01"); err != nil {
			t.Fatalf("1回目のInsertに失敗しました: %v", err)
		}

		// 同じcompany_idで2回目のInsertを実行するとエラーになることを確認する
		err := repo.Insert(300, "enterprise", "2025-01-01")
		if err == nil {
			t.Fatal("重複INSERT時にエラーが発生しませんでした")
		}

		// Duplicate keyエラーであることを確認する
		if !strings.Contains(err.Error(), "Duplicate entry") {
			t.Errorf("期待するDuplicate keyエラーではありません: %v", err)
		}
	})
}
