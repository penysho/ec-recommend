# EC商品レコメンド機能

BedrockとDBを使った商品レコメンド機能の実装です。

## 機能概要

### 実装されたレコメンド手法

1. **コンテンツベースフィルタリング**
   - 商品のタグ、カテゴリ、価格帯に基づく推薦
   - 顧客の過去の購入履歴から嗜好を分析

2. **協調フィルタリング**
   - 類似した購入パターンを持つ顧客を特定
   - 類似顧客が購入した商品を推薦

3. **ハイブリッド手法**
   - 複数の手法を組み合わせて精度向上
   - 重み付け: 協調フィルタリング40% + コンテンツベース40% + トレンド20%

4. **AI強化レコメンド**
   - Amazon Bedrockを使用して推薦理由を生成
   - 信頼度スコアの算出
   - パーソナライズされた説明文の生成

## API エンドポイント

### 1. 商品レコメンド取得

```bash
# GETリクエスト
GET /api/v1/recommendations?customer_id={uuid}&recommendation_type=hybrid&limit=10

# POSTリクエスト
POST /api/v1/recommendations
Content-Type: application/json

{
  "customer_id": "123e4567-e89b-12d3-a456-426614174000",
  "recommendation_type": "hybrid",
  "context_type": "homepage",
  "limit": 10,
  "exclude_owned": true
}
```

**レスポンス例:**

```json
{
  "customer_id": "123e4567-e89b-12d3-a456-426614174000",
  "recommendations": [
    {
      "product_id": "456e7890-e89b-12d3-a456-426614174001",
      "name": "ワイヤレスイヤホン",
      "price": 15000,
      "category_id": 1,
      "category_name": "電子機器",
      "rating_average": 4.5,
      "confidence_score": 0.85,
      "reason": "あなたの過去の電子機器購入履歴と高い評価商品への嗜好から推薦しました"
    }
  ],
  "recommendation_type": "hybrid",
  "context_type": "homepage",
  "generated_at": "2024-01-15T10:30:00Z",
  "metadata": {
    "algorithm_version": "hybrid_v1.0",
    "processing_time_ms": 250,
    "ai_model_used": "amazon.nova-lite-v1:0"
  }
}
```

### 2. 類似商品取得

```bash
GET /api/v1/products/{product_id}/similar?limit=5
```

### 3. トレンド商品取得

```bash
GET /api/v1/products/trending?category_id=1&limit=10
```

### 4. 顧客プロフィール取得

```bash
GET /api/v1/customers/{customer_id}/profile
```

### 5. レコメンド結果のログ記録

```bash
POST /api/v1/recommendations/interactions
Content-Type: application/json

{
  "recommendation_id": "789e1234-e89b-12d3-a456-426614174002",
  "customer_id": "123e4567-e89b-12d3-a456-426614174000",
  "recommended_products": ["456e7890-e89b-12d3-a456-426614174001"],
  "clicked_products": ["456e7890-e89b-12d3-a456-426614174001"],
  "purchased_products": []
}
```

## データベース設計

### 主要テーブル

- **customers**: 顧客情報と嗜好データ
- **products**: 商品情報（タグ、カテゴリ、評価等）
- **orders/order_items**: 注文履歴
- **customer_activities**: 顧客行動ログ
- **recommendation_logs**: レコメンド結果とパフォーマンス追跡

### 分析用ビュー

- **customer_purchase_summary**: 顧客の購入サマリー
- **product_popularity**: 商品の人気度指標

## 設定

### 環境変数

```bash
# データベース設定
DB_HOST=localhost
DB_PORT=5436
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres

# AWS設定
AWS_REGION=ap-northeast-1
BEDROCK_MODEL_ID=amazon.nova-lite-v1:0

# サーバー設定
PORT=8080
LOG_LEVEL=info
```

## 使用方法

### 1. データベースセットアップ

```bash
# データベース起動
make db-up

# スキーマ作成
make db-setup

# サンプルデータ投入
make db-seed
```

### 2. アプリケーション起動

```bash
# 開発環境起動
make dev

# または直接実行
go run cmd/server/main.go
```

### 3. レコメンド機能テスト

```bash
# レコメンド機能のテスト
make recommend-test

# 特定顧客へのレコメンド取得
curl "http://localhost:8080/api/v1/recommendations?customer_id=123e4567-e89b-12d3-a456-426614174000&recommendation_type=hybrid&limit=5"
```

## アーキテクチャ

### レイヤー構成

```
cmd/server/          # アプリケーションエントリーポイント
├── main.go

internal/
├── interfaces/      # インターフェース定義
├── handler/         # HTTPハンドラー
├── service/         # ビジネスロジック
├── repository/      # データアクセス層
├── dto/            # データ転送オブジェクト
├── types/          # 共通型定義
├── config/         # 設定管理
└── middleware/     # HTTPミドルウェア
```

### 依存関係

- **Gin**: HTTPルーター
- **PostgreSQL**: データベース
- **AWS Bedrock**: AI推論
- **lib/pq**: PostgreSQLドライバー
- **google/uuid**: UUID生成

## 拡張性

### 新しいレコメンド手法の追加

1. `RecommendationService`に新しいメソッドを追加
2. `recommendation_type`パラメータに新しい値を追加
3. 必要に応じてデータベースクエリを追加

### パフォーマンス最適化

- Redis等のキャッシュ層の追加
- レコメンド結果の事前計算
- 機械学習モデルの導入

### A/Bテスト対応

- `recommendation_logs`テーブルを活用
- アルゴリズムバージョンの管理
- 成果指標の追跡

## 注意事項

- AWS認証情報が適切に設定されている必要があります
- Bedrockの利用には適切なIAM権限が必要です
- 本番環境では適切なレート制限とエラーハンドリングを実装してください
