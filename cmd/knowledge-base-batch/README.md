# Knowledge Base バッチ処理

Amazon Bedrock Knowledge Basesのために、DBからS3へ商品データをアップロードするバッチ処理です。

## 概要

このバッチ処理は以下を実行します：

1. データベースからアクティブな商品データを取得
2. 商品、カテゴリ、レビュー情報を統合
3. ナレッジベース用のテキストコンテンツを生成
4. 構造化されたJSONドキュメントとしてS3にアップロード

## 前提条件

### 必要な依存関係のインストール

```bash
# AWS SDK v2の必要なサービスを追加
go get github.com/aws/aws-sdk-go-v2/service/s3@latest
go get github.com/aws/aws-sdk-go-v2/service/bedrockagent@latest
go get github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime@latest

# 依存関係を整理
go mod tidy
```

### AWS設定

1. AWS CLIがインストールされ、適切に設定されていること
2. 以下のIAMアクセス許可が必要：
   - S3への読み書きアクセス
   - Bedrock Knowledge Basesへのアクセス

### 環境変数

以下の環境変数を設定してください：

```bash
# データベース接続
export DATABASE_URL="postgres://username:password@localhost/ec_recommend?sslmode=disable"

# S3設定
export S3_BUCKET_NAME="your-bedrock-kb-bucket"
export S3_REGION="us-east-1"
export S3_KEY_PREFIX="products/"

# バッチ処理設定
export BATCH_SIZE="100"
export ENABLE_DEBUG="true"
```

## 実行方法

### 1. バッチ処理の実行

```bash
cd cmd/knowledge-base-batch
go run main.go
```

### 2. ログの確認

バッチ処理は以下の情報をログ出力します：

- 処理対象商品数
- バッチごとの進捗状況
- 成功/失敗件数
- エラー詳細

## 出力されるドキュメント形式

S3にアップロードされるマークダウンドキュメントの構造：

```markdown
---
product_id: prod-123
product_name: 商品名
brand: ブランド名
category: カテゴリ名
category_id: 1
price: 1000.00
rating_average: 4.5
rating_count: 100
tags: [タグ1, タグ2]
is_active: true
generated_at: 2024-01-01T12:00:00Z
---

# 商品名

## 基本情報

**商品ID:** prod-123
**ブランド:** ブランド名
**カテゴリ:** カテゴリ名

## 価格情報

**現在価格:** 1000.00円

## 商品説明

商品の詳細説明がここに表示されます。

## 商品特徴

- 特徴1
- 特徴2
- 特徴3

## タグ

`タグ1` `タグ2`

## 評価情報

**平均評価:** 4.5/5 ⭐
**評価件数:** 100件

## 顧客レビュー

### レビュー 1

**評価:** 5/5 ⭐
**タイトル:** とても良い商品です
**内容:** 期待以上の品質で満足しています...

## 在庫情報

**在庫数:** 50個
```

## S3ファイル構造

```
your-bucket/
└── products/
    └── 2024/01/01/
        ├── product_prod-001.md
        ├── product_prod-002.md
        └── ...
```

## Amazon Bedrock Knowledge Basesとの連携

### 1. Knowledge Baseの作成

AWS Consoleで以下の手順でKnowledge Baseを作成：

1. Amazon Bedrock コンソールを開く
2. 「Knowledge bases」を選択
3. 「Create knowledge base」をクリック
4. データソースとしてS3バケットを指定
5. **ファイル形式**: マークダウン(.md)ファイルがサポートされています
6. エンベディングモデルを選択（推奨：Amazon Titan Text Embeddings v2）
7. ベクトルストアを設定（推奨：Amazon OpenSearch Serverless）

### 対応フォーマット

Amazon Bedrock Knowledge Basesは以下のフォーマットをサポートしています：

- プレーンテキスト (.txt)
- **マークダウン (.md)** ← このバッチ処理で使用
- HTML (.html)
- Microsoft Word (.doc/.docx)
- CSV (.csv)
- Microsoft Excel (.xls/.xlsx)
- PDF (.pdf)

**注意**: JSONファイル(.json)はサポートされていません。

### 2. データ同期

バッチ処理完了後：

1. Knowledge BaseのData Sourceで「Sync」を実行
2. 同期完了まで待機
3. テスト機能で正常に動作することを確認

## トラブルシューティング

### よくある問題

1. **データベース接続エラー**
   - DATABASE_URLの確認
   - データベースが稼働中か確認

2. **S3アップロードエラー**
   - AWSクレデンシャルの確認
   - S3バケットのアクセス許可確認
   - リージョン設定の確認

3. **メモリ不足**
   - BATCH_SIZEを小さくする
   - 大量データの場合は分割処理を検討

### ログ例

```
2024/01/01 12:00:00 Starting Knowledge Base batch process...
2024/01/01 12:00:01 Starting product data extraction...
2024/01/01 12:00:02 Found 1500 products to process
2024/01/01 12:00:05 Processed 100/1500 products (Success: 100, Error: 0)
2024/01/01 12:00:08 Processed 200/1500 products (Success: 200, Error: 0)
...
2024/01/01 12:05:30 Batch processing completed. Success: 1485, Error: 15
2024/01/01 12:05:30 Knowledge Base batch process completed successfully
```

## 次のステップ

1. このバッチを定期実行する場合は、cron jobやAWS Batchを使用
2. データ更新の差分検知機能の実装
3. Custom Connectorsを使用したリアルタイム同期の検討

## 関連ドキュメント

- [Amazon Bedrock Knowledge Bases User Guide](https://docs.aws.amazon.com/bedrock/latest/userguide/knowledge-base.html)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)
