.PHONY: tidy build run-user run-admin migrate seed test test-integration lint db-up db-down db-clean db-logs docker-up docker-down docker-clean docker-migrate docker-seed docker-logs setup setup-docker clean

# 依存関係の整理
tidy:
	go mod tidy

# ビルド
build:
	go build -o tmp/main cmd/backend/main.go

# ローカル実行
run-user: build
	./tmp/main start user --config tmp/config.json

run-admin: build
	./tmp/main start admin --config tmp/config.json

# マイグレーション（ローカル）
migrate:
	go run cmd/backend/main.go init database --config tmp/config.json

# シードデータ作成（ローカル）
seed:
	go run cmd/backend/main.go init seed --config tmp/config.json

# テスト
test:
	go test ./...

# 統合テスト（要: make db-up）。ゴールデン更新は UPDATE_GOLDEN=1 make test-integration
test-integration:
	DSBD_TEST_DB=1 go test -count=1 -v ./pkg/api/...

# Lint
lint:
	golangci-lint run

# DB only (ローカル開発用)
db-up:
	docker compose -f compose.dev.yaml up -d

db-down:
	docker compose -f compose.dev.yaml down

db-clean:
	docker compose -f compose.dev.yaml down -v

db-logs:
	docker compose -f compose.dev.yaml logs -f

# Docker Compose (全サービス)
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-clean:
	docker compose down -v

docker-logs:
	docker compose logs -f

# Docker内でマイグレーション
docker-migrate:
	docker compose exec user-api go run cmd/backend/main.go init database --config tmp/config.json

# Docker内でシードデータ作成
docker-seed:
	docker compose exec user-api go run cmd/backend/main.go init seed --config tmp/config.json

# 初期セットアップ（ローカル用）
setup:
	mkdir -p tmp
	cp configs/config.json tmp/config.json

# 初期セットアップ（Docker用）
setup-docker:
	mkdir -p tmp
	cp configs/config.docker.json tmp/config.json

# クリーン
clean:
	rm -rf tmp/main
