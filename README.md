# HomeNOC WebSystem Backend

![Go-Develop](https://github.com/homenoc/dsbd-backend/workflows/Go-Develop/badge.svg)
## 特徴
* ユーザ認証・登録システム
* グループ認証
* ネットワーク管理システム
* メール送信
* Slack通知
* 各システムとの連携

## MySQL
```
"Error 1406: Data too long for column 'question' at row 1"
```
桁数があふれる可能性があるので、以下のように設定
```
[mysqld]
sql_mode=''
```