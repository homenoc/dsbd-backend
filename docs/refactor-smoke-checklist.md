# リファクタ用 手動スモークチェックリスト

ゴールデンテスト（`make test-integration`）はHTTPレスポンスJSONを凍結するが、
**goroutine起動・WebSocket・外部連携（Slack/mail/Stripe）はHTTPテストに映らない**。
各フェーズのチェックポイントで以下を手動確認する。

## 前提

```bash
make db-up && make migrate && make seed
make run-user    # :8080
make run-admin   # :8081
```

シードアカウント: master@example.com / password（Level 1）, member@example.com / password（Level 2）
Admin: config.json の admin/admin。

## チェック項目

- [ ] **user ログイン / ログアウト**（:8080）: GET /login → tmp token → POST /login（HASH_PASS）→ access token 取得 → logout
- [ ] **admin ログイン**（:8081）: POST /login（USER/PASS ヘッダ）→ access token
- [ ] **group / info 表示**: GET /info でユーザ・グループ・サービスが返る
- [ ] **service 追加（user）** → admin 一覧（GET /service）に反映される
- [ ] **connection 追加（user）** → admin 一覧（GET /connection）に反映される
- [ ] **承認フロー（admin）**: service を pass=true に更新 → connection を open=true → ステータス遷移が正しい
  - P5以降: VLAN 自動採番が承認時に一度だけ行われる
- [ ] **サポートチケット WS チャット双方向**: user がチケット作成 → admin が返信 → 双方向で届く
  （`/ws/v1/support`。HandleMessages / HandleMessagesByAdmin の goroutine 起動に依存。**ゴールデンに映らない最重要項目**）
- [ ] **notice**: admin が notice 追加 → user ダッシュボードに表示
- [ ] **Slack 通知**: dev では notifier のログ行、本番相当では該当チャンネルへ投稿されること

## リリーストレイン時（wire 変更を伴うデプロイ）

seed 済みローカルで backend + dsbd-web + dsbd-web-admin を同時起動し、
上記を全て通してからデプロイする。
