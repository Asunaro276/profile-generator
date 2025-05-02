.PHONY: build run test clean docker-build docker-run docker-clean

# Go実行ファイル
BINARY_NAME=randomuser-server

# ビルド
build:
	@echo "Building..."
	go build -o $(BINARY_NAME) ./cmd/server

# 実行
run:
	@echo "Running..."
	go run ./cmd/server/main.go

# テスト
test:
	@echo "Testing..."
	go test -v ./...
