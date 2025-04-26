.PHONY: build run test clean docker-build docker-run docker-clean

# Go実行ファイル
BINARY_NAME=randomuser-server

# Dockerイメージ名
IMAGE_NAME=randomuser-go
CONTAINER_NAME=randomuser-app

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

# クリーンアップ
clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)

# Dockerビルド
docker-build:
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME) .

# Dockerで実行
docker-run:
	@echo "Running with Docker Compose..."
	docker compose up -d

# Docker停止
docker-stop:
	@echo "Stopping Docker containers..."
	docker compose down

# Docker環境クリーンアップ（データボリュームを含む）
docker-clean:
	@echo "Cleaning Docker resources..."
	docker compose down -v

# 開発用に全て一度にセットアップして実行
dev: docker-build docker-run
	@echo "Development environment ready!"
	@echo "Access the API at: http://localhost:8080/api"
	@echo "Access the web interface at: http://localhost:8080"

# ヘルプ
help:
	@echo "Available commands:"
	@echo "  make build         - Build the Go binary"
	@echo "  make run           - Build and run locally"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run with Docker Compose"
	@echo "  make docker-stop   - Stop Docker containers"
	@echo "  make docker-clean  - Remove Docker resources including volumes"
	@echo "  make dev           - Setup and run development environment"
	@echo "  make help          - Show this help" 
