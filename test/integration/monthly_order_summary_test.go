package integration_test

import (
	"testing"

	"github.com/example/upsert-example/internal/repository"
)

func TestMonthlyOrderSummaryRepository_UpsertFromOrders(t *testing.T) {
	db := setupDB(t)
	repo := repository.NewMonthlyOrderSummaryRepository(db)

	t.Run("集計が新規作成される", func(t *testing.T) {
		cleanTable(t, db, "monthly_order_summaries", "orders")

		// ordersにテストデータを挿入する（2024年1月分）
		_, err := db.Exec(`
			INSERT INTO orders (user_id, amount, status, ordered_at) VALUES
			(1, 1000.00, 'completed', '2024-01-10 10:00:00'),
			(1, 2000.00, 'completed', '2024-01-20 15:00:00')
		`)
		if err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		// 2024年1月を対象に集計UPSERTを実行する
		err = repo.UpsertFromOrders("2024-01-01 00:00:00", "2024-02-01 00:00:00")
		if err != nil {
			t.Fatalf("UpsertFromOrdersに失敗しました: %v", err)
		}

		// 集計結果を確認する
		got, err := repo.FindByUserAndMonth(1, "2024-01-01")
		if err != nil {
			t.Fatalf("FindByUserAndMonthに失敗しました: %v", err)
		}
		if got == nil {
			t.Fatal("集計レコードが存在しません")
		}
		if got.OrderCount != 2 {
			t.Errorf("OrderCount: got %d, want %d", got.OrderCount, 2)
		}
		if got.TotalAmount != 3000.00 {
			t.Errorf("TotalAmount: got %f, want %f", got.TotalAmount, 3000.00)
		}
	})

	t.Run("再集計で既存集計が更新される", func(t *testing.T) {
		cleanTable(t, db, "monthly_order_summaries", "orders")

		// 1回目: 2件の注文を挿入して集計する
		_, err := db.Exec(`
			INSERT INTO orders (user_id, amount, status, ordered_at) VALUES
			(1, 1000.00, 'completed', '2024-01-10 10:00:00'),
			(1, 2000.00, 'completed', '2024-01-20 15:00:00')
		`)
		if err != nil {
			t.Fatalf("1回目のテストデータ挿入に失敗しました: %v", err)
		}

		err = repo.UpsertFromOrders("2024-01-01 00:00:00", "2024-02-01 00:00:00")
		if err != nil {
			t.Fatalf("1回目のUpsertFromOrdersに失敗しました: %v", err)
		}

		// 2回目: さらに注文を追加して再集計する
		_, err = db.Exec(`
			INSERT INTO orders (user_id, amount, status, ordered_at) VALUES
			(1, 500.00, 'completed', '2024-01-25 09:00:00')
		`)
		if err != nil {
			t.Fatalf("2回目のテストデータ挿入に失敗しました: %v", err)
		}

		err = repo.UpsertFromOrders("2024-01-01 00:00:00", "2024-02-01 00:00:00")
		if err != nil {
			t.Fatalf("2回目のUpsertFromOrdersに失敗しました: %v", err)
		}

		// order_countが3件に更新されていることを確認する
		got, err := repo.FindByUserAndMonth(1, "2024-01-01")
		if err != nil {
			t.Fatalf("FindByUserAndMonthに失敗しました: %v", err)
		}
		if got == nil {
			t.Fatal("集計レコードが存在しません")
		}
		if got.OrderCount != 3 {
			t.Errorf("OrderCount: got %d, want %d", got.OrderCount, 3)
		}
		if got.TotalAmount != 3500.00 {
			t.Errorf("TotalAmount: got %f, want %f", got.TotalAmount, 3500.00)
		}
	})

	t.Run("対象期間外の注文は集計されない", func(t *testing.T) {
		cleanTable(t, db, "monthly_order_summaries", "orders")

		// 2024年1月の注文と、対象外（2月）の注文を両方挿入する
		_, err := db.Exec(`
			INSERT INTO orders (user_id, amount, status, ordered_at) VALUES
			(1, 1000.00, 'completed', '2024-01-15 10:00:00'),
			(1, 9999.00, 'completed', '2024-02-01 10:00:00')
		`)
		if err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		// 2024年1月のみを対象として集計UPSERTを実行する
		err = repo.UpsertFromOrders("2024-01-01 00:00:00", "2024-02-01 00:00:00")
		if err != nil {
			t.Fatalf("UpsertFromOrdersに失敗しました: %v", err)
		}

		// 1月の集計に2月分が含まれていないことを確認する
		got, err := repo.FindByUserAndMonth(1, "2024-01-01")
		if err != nil {
			t.Fatalf("FindByUserAndMonthに失敗しました: %v", err)
		}
		if got == nil {
			t.Fatal("集計レコードが存在しません")
		}
		if got.OrderCount != 1 {
			t.Errorf("OrderCount: got %d, want %d (対象期間外の注文が混入している可能性があります)", got.OrderCount, 1)
		}
		if got.TotalAmount != 1000.00 {
			t.Errorf("TotalAmount: got %f, want %f", got.TotalAmount, 1000.00)
		}
	})
}
