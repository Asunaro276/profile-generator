# Randomuser Go実装

このプロジェクトは[Randomuser.me-Node](https://github.com/RandomAPI/Randomuser.me-Node)をGo言語に移植したものです。ランダムなユーザーデータを生成するAPIを提供します。

## 機能

- 複数フォーマットのサポート (JSON, XML, CSV)
- 複数の国籍のデータに対応
- レート制限機能
- 統計情報の記録

## 要件

- Go 1.24以上
- Docker & Docker Compose (推奨)

## 使い方
1. リポジトリをクローン
```bash
git clone https://github.com/yourname/randomuser-go.git
cd randomuser-go
```

2. 依存関係のインストール
```bash
go mod download
```

3. サーバーの起動
```bash
make run
```

4. ブラウザで `http://localhost:8080` にアクセス

## API使用例

### 基本的な使用法
```
GET /api
```

### 結果数の指定
```
GET /api?results=5
```

### シード値の指定（同じ結果を再現）
```
GET /api/?seed=abc
```

### 性別の指定
```
GET /api/?gender=male
```

## ディレクトリ構造

```
randomuser-go/
├── cmd/
│   └── server/
│       └── main.go                 # アプリケーションのエントリーポイント
├── internal/
│   ├── data/                       # ユーザー情報
│   ├── config/
│   │   └── config.go               # 設定管理
│   ├── generator/
│   │   └── generator.go            # ユーザー生成機能
│   └── infrastructure/controller/
│       └── generateuser.go         # ユーザー生成APIのコントローラー
└── go.mod

## ライセンス

MIT

## 謝辞

このプロジェクトは[Randomuser.me-Node](https://github.com/RandomAPI/Randomuser.me-Node)をベースにしています。オリジナルの作者に感謝します。 
