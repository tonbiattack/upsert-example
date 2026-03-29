# upsert-example

MySQL 8.0 における UPSERT（`INSERT ... ON DUPLICATE KEY UPDATE`）と `INSERT IGNORE` の動作を、Go の統合テストで検証するサンプルプロジェクトです。

Qiita 記事「UPSERTはいつ使ってよいか、乱用が設計を壊すパターン」のサンプルコードに対応しています。

## 検証内容

| テストファイル | 検証する操作 |
|---|---|
| `user_settings_test.go` | 設定テーブルへの UPSERT（新規挿入・値の更新・同値再実行） |
| `monthly_order_summary_test.go` | 集計テーブルへのバッチ UPSERT（新規集計・再集計・期間外除外） |
| `contract_test.go` | UPSERT による意図しない上書きパターン・通常 INSERT の重複エラー |
| `user_tag_test.go` | `INSERT IGNORE` のスキップ挙動（RowsAffected=0） |
| `row_count_test.go` | `ROW_COUNT()` の戻り値（INSERT=1・UPDATE=2・変更なし=0） |

## 構成

```
upsert-example/
├── docker-compose.yml
├── go.mod
├── sql/
│   └── schema.sql                              # テーブル定義（Docker 起動時に自動適用）
├── internal/
│   └── repository/
│       ├── user_settings.go
│       ├── monthly_order_summary.go
│       ├── contract.go
│       └── user_tag.go
└── test/
    └── integration/
        ├── testhelper_test.go                  # DB 接続・テーブルクリーンアップ
        ├── user_settings_test.go
        ├── monthly_order_summary_test.go
        ├── contract_test.go
        ├── user_tag_test.go
        └── row_count_test.go
```

## 前提

- Docker がインストールされていること
- Go 1.22 以上がインストールされていること

## 使い方

### 1. MySQL を起動する

```bash
docker compose up -d
```

`sql/schema.sql` が自動的に適用され、テーブルが作成されます。

### 2. テストを実行する

```bash
go test ./test/integration/... -v
```

デフォルトの接続先は `localhost:3308`（docker-compose.yml のポート設定）です。
別の接続先を使う場合は `TEST_DSN` 環境変数で指定できます。

```bash
TEST_DSN="root:password@tcp(localhost:3308)/upsert_test?parseTime=true" \
  go test ./test/integration/... -v
```

### 3. 後片付け

```bash
docker compose down
```

## 動作確認済み環境

- MySQL 8.0
- Go 1.22

## 補足

### `VALUES(col)` 非推奨構文について

MySQL 8.0.20 以降、`ON DUPLICATE KEY UPDATE` での `VALUES(col)` による参照は非推奨になりました。
このプロジェクトのサンプルは、挿入行に別名を付けて参照する推奨の書き方に統一しています。

```sql
INSERT INTO user_settings (user_id, theme, language, updated_at)
VALUES (?, ?, ?, NOW()) AS new
ON DUPLICATE KEY UPDATE
  theme      = new.theme,
  language   = new.language,
  updated_at = NOW();
```

### `ROW_COUNT()` と `updated_at = NOW()` の関係

`ON DUPLICATE KEY UPDATE` での `ROW_COUNT()` は、更新対象のカラムがすべて同じ値だった場合に 0 を返します。
`updated_at = NOW()` のように実行のたびに値が変わるカラムを含めている場合、業務的な値が変わっていなくても UPDATE は常に発動し、ROW_COUNT は 2 になります。
`row_count_test.go` の「変更なし」ケースでは、この挙動を正確に検証するため `updated_at` に固定値を使う専用 SQL で確認しています。
