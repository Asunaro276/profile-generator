# Randomuser Go実装

このプロジェクトは[Randomuser.me-Node](https://github.com/RandomAPI/Randomuser.me-Node)をGo言語に移植したものです。ランダムなユーザーデータを生成するAPIを提供します。

## 要件
- Go 1.24以上

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

4. ブラウザで `http://localhost:8080/api` にアクセス

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
GET /api/?seed=12345
```

### 性別の指定
```
GET /api/?gender=male
```

### ページ番号の指定
```
GET /api/?page=2
```

## ディレクトリ構造

```
randomuser-go/
├── cmd/
│   └── server/
│       └── main.go                 # アプリケーションのエントリーポイント
├── internal/
│   ├── config/                     # 設定管理
│   ├── data/                       # ユーザー情報
│   ├── generator/                  # ユーザー生成機能
│   ├── infrastructure/controller/  # ユーザー生成APIのコントローラー
│   └── model/                      # ユーザー情報のモデル
└── go.mod

## ライセンス

MIT

## 謝辞

このプロジェクトは[Randomuser.me-Node](https://github.com/RandomAPI/Randomuser.me-Node)をベースにしています。オリジナルの作者に感謝します。 
