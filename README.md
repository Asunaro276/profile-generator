# Randomuser Go実装

このプロジェクトは[Randomuser.me-Node](https://github.com/RandomAPI/Randomuser.me-Node)をGo言語に移植したものです。ランダムなユーザーデータを生成するAPIを提供します。

## 機能

- 複数フォーマットのサポート (JSON, XML, CSV)
- 複数の国籍のデータに対応
- レート制限機能
- 統計情報の記録

## 要件

- Go 1.20以上
- Docker & Docker Compose (推奨)

## 使い方

### Dockerでの実行 (推奨)

1. リポジトリをクローン
```bash
git clone https://github.com/yourname/randomuser-go.git
cd randomuser-go
```

2. Docker Composeで起動
```bash
docker compose up -d
```

3. ブラウザで `http://localhost:8080` にアクセス

### 手動での実行

1. リポジトリをクローン
```bash
git clone https://github.com/yourname/randomuser-go.git
cd randomuser-go
```

2. 依存関係のインストール
```bash
go mod download
```

3. MongoDBの起動
```bash
# MongoDBをローカルで実行する必要があります
```

4. サーバーの起動
```bash
go run cmd/server/main.go
```

5. ブラウザで `http://localhost:8080` にアクセス

## API使用例

### 基本的な使用法
```
GET /api/
```

### 結果数の指定
```
GET /api/?results=5
```

### 特定の国籍の指定
```
GET /api/?nat=us,gb
```

### フォーマット指定
```
GET /api/?format=xml
GET /api/?format=csv
GET /api/?format=json
```

### シード値の指定（同じ結果を再現）
```
GET /api/?seed=abc
```

## ディレクトリ構造

```
randomuser-go/
├── cmd/
│   └── server/
│       └── main.go       # アプリケーションのエントリーポイント
├── internal/
│   ├── api/
│   │   ├── handlers.go   # APIハンドラー
│   │   └── routes.go     # ルーティング設定
│   ├── config/
│   │   └── config.go     # 設定管理
│   ├── db/
│   │   └── mongodb.go    # MongoDBとの接続
│   ├── generator/
│   │   └── generator.go  # ユーザー生成機能
│   └── models/
│       └── request.go    # リクエストモデル
├── web/
│   ├── static/           # 静的ファイル
│   └── templates/        # HTMLテンプレート
├── api/                  # 国籍別データ
├── Dockerfile            # Dockerビルド設定
├── compose.yml           # Docker Compose設定
├── config.json           # アプリケーション設定
└── go.mod
```

## Docker構成

- `compose.yml` - Docker Compose設定ファイル
- `Dockerfile` - アプリケーションのDockerイメージビルド設定
- MongoDB用のボリュームマウント (データ永続化)

### 環境変数

Docker環境では以下の環境変数が利用可能です：

- `MONGODB_URI` - MongoDB接続URI
- `MONGODB_DATABASE` - 使用するデータベース名

## ライセンス

MIT

## 謝辞

このプロジェクトは[Randomuser.me-Node](https://github.com/RandomAPI/Randomuser.me-Node)をベースにしています。オリジナルの作者に感謝します。 
