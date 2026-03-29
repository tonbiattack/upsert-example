package repository

import "database/sql"

// UserSettings はuser_settingsテーブルの1レコードを表す。
type UserSettings struct {
	UserID    int
	Theme     string
	Language  string
	UpdatedAt string
}

// UserSettingsRepository はuser_settingsテーブルへのアクセスを担う。
type UserSettingsRepository struct {
	db *sql.DB
}

// NewUserSettingsRepository はUserSettingsRepositoryを生成して返す。
func NewUserSettingsRepository(db *sql.DB) *UserSettingsRepository {
	return &UserSettingsRepository{db: db}
}

// Upsert は記事のサンプルSQLをそのまま実行する。
// 対象ユーザーが存在しなければINSERT、存在すればtheme・language・updated_atをUPDATEする。
func (r *UserSettingsRepository) Upsert(userID int, theme, language string) error {
	_, err := r.db.Exec(`
		INSERT INTO user_settings (user_id, theme, language, updated_at)
		VALUES (?, ?, ?, NOW()) AS new
		ON DUPLICATE KEY UPDATE
		  theme      = new.theme,
		  language   = new.language,
		  updated_at = NOW()
	`, userID, theme, language)
	return err
}

// UpsertWithRowsAffected はUpsertと同じSQLを実行し、RowsAffectedも返す。
// ROW_COUNT()の戻り値検証用。
func (r *UserSettingsRepository) UpsertWithRowsAffected(userID int, theme, language string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO user_settings (user_id, theme, language, updated_at)
		VALUES (?, ?, ?, NOW()) AS new
		ON DUPLICATE KEY UPDATE
		  theme      = new.theme,
		  language   = new.language,
		  updated_at = NOW()
	`, userID, theme, language)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// FindByUserID は指定ユーザーIDの設定を返す。レコードが存在しない場合はnilを返す。
func (r *UserSettingsRepository) FindByUserID(userID int) (*UserSettings, error) {
	row := r.db.QueryRow(`SELECT user_id, theme, language, updated_at FROM user_settings WHERE user_id = ?`, userID)
	s := &UserSettings{}
	err := row.Scan(&s.UserID, &s.Theme, &s.Language, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}
