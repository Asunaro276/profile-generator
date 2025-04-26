FROM golang:1.20-alpine AS builder

WORKDIR /build

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o randomuser-server ./cmd/server

# 実行用の軽量イメージ
FROM alpine:latest

WORKDIR /app

# 必要なファイルのコピー
COPY --from=builder /build/randomuser-server .
COPY --from=builder /build/config.json .
COPY --from=builder /build/web ./web
COPY --from=builder /build/api ./api

# タイムゾーン設定
RUN apk --no-cache add tzdata && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
  echo "Asia/Tokyo" > /etc/timezone && \
  apk del tzdata

# アプリケーションの実行
EXPOSE 8080
CMD ["./randomuser-server"] 
