# HomeNOC WebSystem Backend

![Go-Develop](https://github.com/homenoc/dsbd-backend/workflows/Go-Develop/badge.svg)
## 特徴
* ユーザ認証・登録システム
* グループ認証
* ネットワーク管理システム
* メール送信
* Slack通知
* 各システムとの連携

## 初回実行時
```
cd cmd/backend
go run . init database --config config.json
go run . init database 
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