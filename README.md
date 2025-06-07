# EC-Recommend: AWS Bedrock AI API Server

GolangとGinフレームワークを使用して構築されたAWS Bedrock生成AIアプリケーションです。

## 機能

- **質問回答API**: 単発の質問に対してAIが回答を生成
- **チャットAPI**: 会話履歴を含むチャット形式での対話
- **ヘルスチェック**: アプリケーションの稼働状況確認
- **グレースフルシャットダウン**: 安全なサーバー停止
- **CORS対応**: クロスオリジンリクエストサポート
- **構造化ログ**: リクエスト/レスポンスのログ出力

## アーキテクチャ

```
cmd/server/           # アプリケーションエントリーポイント
internal/
  ├── config/         # 設定管理
  ├── dto/           # データ転送オブジェクト
  ├── handler/       # HTTPハンドラー
  ├── middleware/    # Ginミドルウェア
  ├── router/        # ルーティング設定
  └── service/       # ビジネスロジック（AWS Bedrock連携）
examples/             # APIテスト用サンプル
```

## 前提条件

- Go 1.21以上
- AWS アカウントとBedrock利用権限
- AWS認証情報の設定

## セットアップ

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd ec-recommend
```

### 2. 依存関係のインストール

```bash
make deps
```

### 3. 環境変数の設定

```bash
# 必須
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your-access-key-id
export AWS_SECRET_ACCESS_KEY=your-secret-access-key

# オプション
export PORT=8080
export BEDROCK_MODEL_ID=anthropic.claude-3-haiku-20240307-v1:0
export LOG_LEVEL=info
```

### 4. ビルドと実行

```bash
# ビルド
make build

# 実行
make run

# または直接実行
./ec-recommend
```

## API エンドポイント

### ヘルスチェック

```http
GET /health
```

### 質問回答API

```http
POST /api/v1/ai/ask
Content-Type: application/json

{
  "question": "GoでWebアプリケーションを作る際のベストプラクティスを教えてください。"
}
```

**レスポンス例:**

```json
{
  "answer": "Goでのウェブアプリケーション開発のベストプラクティスをご紹介します...",
  "usage": {
    "input_tokens": 142,
    "output_tokens": 368
  },
  "timestamp": "2024-01-20T10:30:00Z"
}
```

### チャットAPI

```http
POST /api/v1/ai/chat
Content-Type: application/json

{
  "messages": [
    {
      "role": "user",
      "content": "AWS Bedrockとは何ですか？"
    }
  ]
}
```

**レスポンス例:**

```json
{
  "message": "AWS Bedrockは、Amazon Web Servicesが提供する完全マネージド型の生成AI基盤サービスです...",
  "usage": {
    "input_tokens": 89,
    "output_tokens": 156
  },
  "timestamp": "2024-01-20T10:30:00Z"
}
```

## 開発

### テストの実行

```bash
make test
```

### コードフォーマット

```bash
make fmt
```

### リンター実行

```bash
make lint
```

### 開発モード（ホットリロード）

```bash
# air をインストール
go install github.com/cosmtrek/air@latest

# 開発モード実行
make dev
```

## Docker

### イメージのビルド

```bash
make docker-build
```

### コンテナ実行

```bash
docker run -p 8080:8080 \
  -e AWS_REGION=us-east-1 \
  -e AWS_ACCESS_KEY_ID=your-access-key \
  -e AWS_SECRET_ACCESS_KEY=your-secret-key \
  ec-recommend
```

## 対応Bedrockモデル

現在、以下のモデルに対応しています：

### Claude モデル

- `anthropic.claude-3-haiku-20240307-v1:0`
- `anthropic.claude-3-sonnet-20240229-v1:0`
- `anthropic.claude-3-opus-20240229-v1:0`

### Amazon Nova モデル

- `amazon.nova-micro-v1:0`
- `amazon.nova-lite-v1:0`
- `amazon.nova-pro-v1:0`

### Amazon Titan モデル

- `amazon.titan-text-express-v1`
- `amazon.titan-text-lite-v1`

## トラブルシューティング

### 1. AWS認証エラー

```
Failed to load AWS configuration
```

- AWS認証情報が正しく設定されているか確認
- IAMユーザーにBedrock利用権限があるか確認

### 2. Bedrock利用権限エラー

```
AccessDeniedException
```

- AWSコンソールでBedrock利用申請が完了しているか確認
- 指定リージョンでBedrockが利用可能か確認

### 3. モデルアクセスエラー

```
ValidationException: The provided model identifier is invalid
```

- 指定したモデルIDが正しいか確認
- Bedrockコンソールでモデルアクセスが有効になっているか確認

## コントリビューション

1. フォークしてください
2. フィーチャーブランチを作成してください (`git checkout -b feature/amazing-feature`)
3. 変更をコミットしてください (`git commit -m 'Add amazing feature'`)
4. ブランチにプッシュしてください (`git push origin feature/amazing-feature`)
5. プルリクエストを開いてください

## ライセンス

MIT License - 詳細は[LICENSE](LICENSE)ファイルを参照してください。
