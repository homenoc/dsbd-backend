# HomeNOC WebSystem Backend

![Go-Develop](https://github.com/homenoc/dsbd-backend/workflows/Go-Develop/badge.svg)
## 特徴
* ユーザ認証・登録システム
* グループ認証
* ネットワーク管理システム
* メール送信
* Slack通知
* 各システムとの連携

## セットアップ

### 開発ツールのインストール (aqua)
[aqua](https://aquaproj.github.io/)を使用して開発ツール（Go, Air, golangci-lint）を管理しています。

```shell
# aquaのインストール（未インストールの場合）
brew install aquaproj/aqua/aqua

# PATHの設定（.zshrc または .bashrc に追加）
export PATH="$(aqua root-dir)/bin:$PATH"

# 開発ツールのインストール
aqua i
```

### ローカル開発（推奨）
DBのみDockerで起動し、Goアプリはローカルで実行します。

```shell
# 初期設定（ローカル用config）
make setup

# DB起動
make db-up

# マイグレーション（初回のみ）
make migrate

# シードデータ作成（テストユーザー等）
make seed

# User API起動 (port 8080)
make run-user

# Admin API起動 (port 8081)
make run-admin

# DB停止
make db-down
```

### Docker Compose（全サービス）
全てのサービスをDockerで起動します。

```shell
# 初期設定（Docker用config）
make setup-docker

# コンテナ起動
make docker-up

# マイグレーション（初回のみ）
make docker-migrate

# シードデータ作成（テストユーザー等）
make docker-seed

# ログ確認
make docker-logs

# コンテナ停止
make docker-down
```

- UserAPI  : http://localhost:8080
- AdminAPI : http://localhost:8081

### テストアカウント
`make seed` で作成されるテストアカウント（User API用）：

| 種別 | メールアドレス | パスワード | 権限 |
|------|---------------|-----------|------|
| マスター | master@example.com | password | Level 1: 申請・変更・閲覧 |
| 一般メンバー | member@example.com | password | Level 2: 一般メンバー |

※ 反社チェック未同意状態で作成されます（`PUT /api/v1/user/antisocial/agree` で同意）

Admin APIは `config.json` のBasic認証を使用（デフォルト: admin/admin）

### テスト
```shell
make test
```

## MySQL
```
"Error 1406: Data too long for column 'question' at row 1"
```
桁数があふれる可能性があるので、以下のように設定
```
[mysqld]
sql_mode=''
```

## Database
https://drawsql.app/y-net/diagrams/dsbd-backend/embed