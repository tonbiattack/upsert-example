package repository

import "database/sql"

// MonthlyOrderSummary はmonthly_order_summariesテーブルの1レコードを表す。
type MonthlyOrderSummary struct {
	UserID       int
	TargetMonth  string
	OrderCount   int
	TotalAmount  float64
	AggregatedAt string
}

// MonthlyOrderSummaryRepository はmonthly_order_summariesテーブルへのアクセスを担う。
type MonthlyOrderSummaryRepository struct {
	db *sql.DB
}

// NewMonthlyOrderSummaryRepository はMonthlyOrderSummaryRepositoryを生成して返す。
func NewMonthlyOrderSummaryRepository(db *sql.DB) *MonthlyOrderSummaryRepository {
	return &MonthlyOrderSummaryRepository{db: db}
}

// UpsertFromOrders はordersテーブルから集計し、monthly_order_summariesにUPSERTする。
// rangeStartとrangeEndは "2024-01-01 00:00:00" 形式で指定する。
// 対象期間はrangeStart以上rangeEnd未満とする。
func (r *MonthlyOrderSummaryRepository) UpsertFromOrders(rangeStart, rangeEnd string) error {
	_, err := r.db.Exec(`
		INSERT INTO monthly_order_summaries
		  (user_id, target_month, order_count, total_amount, aggregated_at)
		SELECT
		  src.user_id,
		  src.target_month,
		  src.order_count,
		  src.total_amount,
		  src.aggregated_at
		FROM (
		  SELECT
		    user_id,
		    DATE(ordered_at - INTERVAL (DAY(ordered_at) - 1) DAY) AS target_month,
		    COUNT(*) AS order_count,
		    SUM(amount) AS total_amount,
		    NOW() AS aggregated_at
		  FROM orders
		  WHERE ordered_at >= ? AND ordered_at < ?
		  GROUP BY
		    user_id,
		    DATE(ordered_at - INTERVAL (DAY(ordered_at) - 1) DAY)
		) AS src
		ON DUPLICATE KEY UPDATE
		  order_count   = src.order_count,
		  total_amount  = src.total_amount,
		  aggregated_at = src.aggregated_at
	`, rangeStart, rangeEnd)
	return err
}

// FindByUserAndMonth は指定ユーザーID・対象月の集計レコードを返す。
// targetMonthは "2024-01-01" 形式で指定する。レコードが存在しない場合はnilを返す。
func (r *MonthlyOrderSummaryRepository) FindByUserAndMonth(userID int, targetMonth string) (*MonthlyOrderSummary, error) {
	row := r.db.QueryRow(`
		SELECT user_id, target_month, order_count, total_amount, aggregated_at
		FROM monthly_order_summaries
		WHERE user_id = ? AND target_month = ?
	`, userID, targetMonth)
	s := &MonthlyOrderSummary{}
	err := row.Scan(&s.UserID, &s.TargetMonth, &s.OrderCount, &s.TotalAmount, &s.AggregatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}
