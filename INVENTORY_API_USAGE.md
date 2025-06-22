# 在庫管理API使用ガイド

## 概要

ec-extensionプロジェクトに在庫管理エンドポイントと権限判定機能を追加しました。

## 実装した機能

### 1. 権限判定システム
- **認証ミドルウェア**: Bearer tokenによる認証
- **認可ミドルウェア**: ロールベースのアクセス制御
- **ユーザーロール**:
  - `admin`: 全ての操作が可能
  - `employee`: 在庫管理操作が可能（バッチ更新以外）
  - `customer`: 読み取り専用操作のみ

### 2. 在庫管理エンドポイント

| エンドポイント | メソッド | 必要な権限 | 説明 |
|---|---|---|---|
| `/api/v1/inventory` | GET | 認証済み全ユーザー | 在庫一覧取得 |
| `/api/v1/inventory/products/:product_id` | GET | 認証済み全ユーザー | 特定商品の在庫情報取得 |
| `/api/v1/inventory/products/:product_id` | PUT | Admin/Employee | 在庫数更新 |
| `/api/v1/inventory/transactions` | POST | Admin/Employee | 在庫トランザクション実行 |
| `/api/v1/inventory/batch-update` | POST | Admin | バッチ更新 |
| `/api/v1/inventory/alerts` | GET | Admin/Employee | 在庫アラート取得 |
| `/api/v1/inventory/stats` | GET | Admin/Employee | 在庫統計取得 |
| `/api/v1/inventory/products/:product_id/history` | GET | 認証済み全ユーザー | 在庫履歴取得 |

## 認証トークン（テスト用）

```bash
# 管理者
Authorization: Bearer admin-token-123

# 従業員
Authorization: Bearer employee-token-456

# 顧客
Authorization: Bearer customer-token-789
```

## APIリクエスト例

### 1. 在庫一覧取得

```bash
curl -X GET "http://localhost:8080/api/v1/inventory" \
  -H "Authorization: Bearer admin-token-123" \
  -H "Content-Type: application/json"
```

### 2. 特定商品の在庫情報取得

```bash
curl -X GET "http://localhost:8080/api/v1/inventory/products/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer admin-token-123" \
  -H "Content-Type: application/json"
```

### 3. 在庫数更新（管理者・従業員のみ）

```bash
curl -X PUT "http://localhost:8080/api/v1/inventory/products/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer admin-token-123" \
  -H "Content-Type: application/json" \
  -d '{
    "stock_quantity": 50,
    "reason": "Manual adjustment"
  }'
```

### 4. 在庫トランザクション実行（管理者・従業員のみ）

```bash
# 在庫増加
curl -X POST "http://localhost:8080/api/v1/inventory/transactions" \
  -H "Authorization: Bearer admin-token-123" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "transaction_type": "increase",
    "quantity": 10,
    "reason": "New stock arrival"
  }'

# 在庫減少
curl -X POST "http://localhost:8080/api/v1/inventory/transactions" \
  -H "Authorization: Bearer admin-token-123" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "transaction_type": "decrease",
    "quantity": 5,
    "reason": "Product return"
  }'
```

### 5. バッチ更新（管理者のみ）

```bash
curl -X POST "http://localhost:8080/api/v1/inventory/batch-update" \
  -H "Authorization: Bearer admin-token-123" \
  -H "Content-Type: application/json" \
  -d '{
    "updates": [
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440000",
        "stock_quantity": 100,
        "reason": "Monthly inventory adjustment"
      },
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440001",
        "stock_quantity": 75,
        "reason": "Monthly inventory adjustment"
      }
    ]
  }'
```

### 6. 在庫アラート取得（管理者・従業員のみ）

```bash
curl -X GET "http://localhost:8080/api/v1/inventory/alerts?threshold=10" \
  -H "Authorization: Bearer admin-token-123" \
  -H "Content-Type: application/json"
```

### 7. 在庫統計取得（管理者・従業員のみ）

```bash
curl -X GET "http://localhost:8080/api/v1/inventory/stats" \
  -H "Authorization: Bearer admin-token-123" \
  -H "Content-Type: application/json"
```

## 既存エンドポイントの権限追加

既存の商品・顧客情報エンドポイントにも権限判定を追加しました：

### 推薦エンドポイント（認証が必要）

```bash
curl -X GET "http://localhost:8080/api/v1/recommendations?customer_id=550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer customer-token-789" \
  -H "Content-Type: application/json"
```

### 顧客プロファイル（認証が必要）

```bash
curl -X GET "http://localhost:8080/api/v1/customers/550e8400-e29b-41d4-a716-446655440000/profile" \
  -H "Authorization: Bearer customer-token-789" \
  -H "Content-Type: application/json"
```

### 商品エンドポイント（認証が必要）

```bash
curl -X GET "http://localhost:8080/api/v1/products/trending" \
  -H "Authorization: Bearer customer-token-789" \
  -H "Content-Type: application/json"
```

## エラーレスポンス例

### 認証エラー
```json
{
  "error": "Authorization header is required",
  "code": 401,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 権限不足エラー
```json
{
  "error": "Insufficient permissions",
  "code": 403,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 無効なリクエストエラー
```json
{
  "error": "Invalid product_id format",
  "code": 400,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 検索・フィルタリング機能

在庫一覧取得では以下のクエリパラメータが使用可能です：

```bash
curl -X GET "http://localhost:8080/api/v1/inventory?query=iPhone&category_id=1&stock_status=low_stock&min_stock=1&max_stock=10&page=1&page_size=20&sort_by=stock_quantity&sort_direction=asc" \
  -H "Authorization: Bearer admin-token-123" \
  -H "Content-Type: application/json"
```

## セキュリティ考慮事項

1. **本番環境での認証**: 現在は簡易的なトークン認証を使用していますが、本番環境ではJWTトークンやOAuth2.0などの適切な認証システムを実装してください。

2. **データアクセス制御**: 顧客は自分のデータのみアクセス可能にするよう、追加の制御ロジックを実装することを推奨します。

3. **ログ記録**: 在庫操作のログ記録機能は基本的な実装のみ含まれています。本格的な監査ログシステムの実装を検討してください。

## 開発・テスト

アプリケーションを起動すると、コンソールに利用可能なエンドポイントとテスト用トークンが表示されます：

```bash
cd cmd/server
go run main.go
```

これで在庫系エンドポイントに権限判定が追加され、既存のエンドポイントにも適切な認証・認可が実装されました。
