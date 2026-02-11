# AWS Checker

AWSリソースのセキュリティ設定を監査するCLIツール

## 概要

`awscheck`は、AWS環境のセキュリティ設定を自動的にチェックし、潜在的なリスクを発見するためのコマンドラインツールです。

## 機能

### S3バケット監査

S3バケットのPublicAccessBlock設定をチェックし、以下の3種類のステータスで判定します：

- ✅ **SAFE**: すべてのパブリックアクセスブロック設定が有効
- ⚠️  **WARNING: NO CONFIGURATION**: PublicAccessBlock設定が存在しない（潜在的なリスク）
- ⚠️  **WARNING: PARTIAL CONFIGURATION**: 一部の設定のみ有効（要確認）

## インストール

### 前提条件

- Go 1.21以上
- AWS認証情報の設定
- make（オプション、ビルドを簡単にするため）

### ビルド

```bash
# リポジトリをクローン
git clone <repository-url>
cd my-aws-checker

# ビルド（Makefile使用）
make build

# または、直接Goコマンドで
go build -o awscheck
```

## 使い方

### 1. AWS認証情報の設定

以下のいずれかの方法でAWS認証情報を設定してください：

#### 方法A: 環境変数

```bash
export AWS_ACCESS_KEY_ID=your-access-key
export AWS_SECRET_ACCESS_KEY=your-secret-key
export AWS_REGION=ap-northeast-1
```

#### 方法B: AWS CLIの設定ファイル

```bash
# ~/.aws/credentials
[default]
aws_access_key_id = your-access-key
aws_secret_access_key = your-secret-key

# ~/.aws/config
[default]
region = ap-northeast-1
```

### 2. S3バケットのチェック

```bash
./awscheck s3
```

## 出力例

```
S3 check start
BUCKET NAME                                REGION           STATUS
-----------                                ------           ------
cf-templates-511105rkg24l-ap-northeast-1   ap-northeast-1   SAFE
codepipeline-ap-northeast-1-495187879750   ap-northeast-1   WARNING: NO CONFIGURATION
discord-bot-14-bucket                      ap-northeast-1   SAFE
webhosting-cloudfront24-226                ap-northeast-1   WARNING: PARTIAL CONFIGURATION
```

### ステータスの意味

#### SAFE
すべてのパブリックアクセスブロック設定が有効化されています：
- `BlockPublicAcls`: 有効
- `BlockPublicPolicy`: 有効
- `IgnorePublicAcls`: 有効
- `RestrictPublicBuckets`: 有効

#### WARNING: NO CONFIGURATION
PublicAccessBlock設定が存在しません。バケットが意図せず公開される可能性があります。

#### WARNING: PARTIAL CONFIGURATION
一部の設定のみ有効です。設定を確認して、必要に応じて有効化してください。

## アーキテクチャ

```
my-aws-checker/
├── cmd/                 # Cobraコマンド定義
│   ├── root.go         # ルートコマンド
│   └── s3.go           # S3チェックコマンド
├── internal/
│   ├── aws/            # AWS設定の共通処理
│   │   └── config.go   # AWS設定読み込み
│   └── s3/             # S3監査ロジック
│       ├── models.go        # ドメインモデル
│       ├── interface.go     # インターフェース定義
│       ├── client.go        # AWS SDKラッパー
│       ├── client_test.go   # Clientのテスト
│       ├── checker.go       # ビジネスロジック
│       └── checker_test.go  # Checkerのテスト
├── main.go
├── Makefile
└── README.md
```

### 設計思想

- **レイヤー分離アーキテクチャ**: コマンド層、ビジネスロジック層、データアクセス層を分離
- **インターフェースベースの設計**: テスト可能な設計
- **依存性の注入**: モックを使った単体テスト

## 開発

### テストの実行

```bash
# すべてのテスト
make test

# または
go test ./...

# 詳細出力付き
go test -v ./internal/s3/
```

## ライセンス

MIT License
