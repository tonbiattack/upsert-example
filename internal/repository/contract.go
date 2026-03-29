package repository

import "database/sql"

// Contract はcontractsテーブルの1レコードを表す。
type Contract struct {
	ID        int
	CompanyID int
	Plan      string
	SignedAt  string
}

// ContractRepository はcontractsテーブルへのアクセスを担う。
type ContractRepository struct {
	db *sql.DB
}

// NewContractRepository はContractRepositoryを生成して返す。
func NewContractRepository(db *sql.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

// Upsert は記事のサンプルSQLをそのまま実行する。
// company_idが重複した場合、ON DUPLICATE KEY UPDATEによってplanとsigned_atが上書きされる。
// これは「意図しない上書き」の検証パターンとして使用する。
// 戻り値はRowsAffectedで、INSERT時は1、UPDATE時は2となる。
func (r *ContractRepository) Upsert(companyID int, plan, signedAt string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO contracts (company_id, plan, signed_at)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
		  plan      = VALUES(plan),
		  signed_at = VALUES(signed_at)
	`, companyID, plan, signedAt)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// FindByCompanyID は指定会社IDの契約を返す。レコードが存在しない場合はnilを返す。
func (r *ContractRepository) FindByCompanyID(companyID int) (*Contract, error) {
	row := r.db.QueryRow(`SELECT id, company_id, plan, signed_at FROM contracts WHERE company_id = ?`, companyID)
	c := &Contract{}
	err := row.Scan(&c.ID, &c.CompanyID, &c.Plan, &c.SignedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

// Insert は通常のINSERTを実行する。company_idが重複している場合はDuplicate keyエラーを返す。
func (r *ContractRepository) Insert(companyID int, plan, signedAt string) error {
	_, err := r.db.Exec(`
		INSERT INTO contracts (company_id, plan, signed_at)
		VALUES (?, ?, ?)
	`, companyID, plan, signedAt)
	return err
}
