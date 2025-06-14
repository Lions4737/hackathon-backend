# 📊 DB Schema: Twitter風アプリ用

## ✅ users テーブル

| カラム名    | データ型      | 制約                          | 説明                          |
|:----------- |:-------------|:------------------------------|:------------------------------|
| id          | CHAR(26)     | PRIMARY KEY                   | ユーザーの一意なID（ULID）    |
| name        | VARCHAR(50)  | NOT NULL                      | 表示名                        |
| email       | VARCHAR(255) | UNIQUE, NOT NULL              | ログイン用メールアドレス      |
| created_at  | DATETIME     | DEFAULT CURRENT_TIMESTAMP     | 作成日時                      |

---

## ✅ tweets テーブル

| カラム名    | データ型      | 制約                                         | 説明                                      |
|:----------- |:-------------|:----------------------------------------------|:-------------------------------------------|
| id          | CHAR(26)     | PRIMARY KEY                                   | ツイートの一意なID（ULID）                |
| user_id     | CHAR(26)     | NOT NULL, FOREIGN KEY → users(id)             | 投稿したユーザーのID                      |
| content     | TEXT         | NOT NULL                                      | ツイート本文                              |
| reply_to    | CHAR(26)     | NULLABLE, FOREIGN KEY → tweets(id)            | リプライ先ツイートID（NULLなら通常投稿）  |
| created_at  | DATETIME     | DEFAULT CURRENT_TIMESTAMP                     | 作成日時                                  |

---

## ✅ likes テーブル

| カラム名    | データ型      | 制約                                         | 説明                         |
|:----------- |:-------------|:----------------------------------------------|:-----------------------------|
| id          | CHAR(26)     | PRIMARY KEY                                   | いいねの一意なID（ULID）     |
| tweet_id    | CHAR(26)     | NOT NULL, FOREIGN KEY → tweets(id)            | いいね対象のツイートID       |
| user_id     | CHAR(26)     | NOT NULL, FOREIGN KEY → users(id)             | いいねしたユーザーID         |
| created_at  | DATETIME     | DEFAULT CURRENT_TIMESTAMP                     | いいねした日時               |
