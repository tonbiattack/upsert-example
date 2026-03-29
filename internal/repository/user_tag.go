package repository

import "database/sql"

// UserTagRepository はuser_tagsテーブルへのアクセスを担う。
type UserTagRepository struct {
	db *sql.DB
}

// NewUserTagRepository はUserTagRepositoryを生成して返す。
func NewUserTagRepository(db *sql.DB) *UserTagRepository {
	return &UserTagRepository{db: db}
}

// InsertIgnore はINSERT IGNOREで実行する。
// 主キーが重複している場合はエラーを返さずスキップされ、RowsAffectedが0になる。
// 新規挿入時はRowsAffectedが1になる。
func (r *UserTagRepository) InsertIgnore(userID, tagID int) (int64, error) {
	result, err := r.db.Exec(`
		INSERT IGNORE INTO user_tags (user_id, tag_id)
		VALUES (?, ?)
	`, userID, tagID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Exists は指定ユーザーID・タグIDの組み合わせが存在するかを返す。
func (r *UserTagRepository) Exists(userID, tagID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM user_tags WHERE user_id = ? AND tag_id = ?
	`, userID, tagID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
