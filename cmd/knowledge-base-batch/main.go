package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"ec-recommend/internal/repository/models"
)

type Config struct {
	DatabaseURL  string
	S3BucketName string
	S3Region     string
	S3KeyPrefix  string
	BatchSize    int
	EnableDebug  bool
}

type ProductKnowledgeDocument struct {
	ProductID   string          `json:"product_id"`
	ProductName string          `json:"product_name"`
	Content     string          `json:"content"`
	Metadata    ProductMetadata `json:"metadata"`
	GeneratedAt time.Time       `json:"generated_at"`
}

type ProductMetadata struct {
	Brand           string   `json:"brand,omitempty"`
	Category        string   `json:"category"`
	CategoryID      int      `json:"category_id"`
	Price           string   `json:"price"`
	OriginalPrice   string   `json:"original_price,omitempty"`
	RatingAverage   string   `json:"rating_average,omitempty"`
	RatingCount     int      `json:"rating_count,omitempty"`
	PopularityScore int      `json:"popularity_score,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	IsActive        bool     `json:"is_active"`
}

type ReviewSummary struct {
	Rating  int    `json:"rating"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type S3Uploader struct {
	client *s3.Client
	bucket string
	prefix string
}

func main() {
	log.Println("Starting Knowledge Base batch process...")

	config := loadConfig()

	// データベース接続
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// AWS設定
	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(config.S3Region))
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// S3クライアント初期化
	uploader := &S3Uploader{
		client: s3.NewFromConfig(awsConfig),
		bucket: config.S3BucketName,
		prefix: config.S3KeyPrefix,
	}

	// バッチ処理実行
	ctx := context.Background()
	processor := NewKnowledgeBaseBatchProcessor(db, uploader, config)

	if err := processor.ProcessAllProducts(ctx); err != nil {
		log.Fatalf("Batch process failed: %v", err)
	}

	log.Println("Knowledge Base batch process completed successfully")
}

func loadConfig() *Config {
	return &Config{
		DatabaseURL:  getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5436/postgres?sslmode=disable"),
		S3BucketName: getEnvOrDefault("S3_BUCKET_NAME", "ec-recommend-knowledge-base-local"),
		S3Region:     getEnvOrDefault("S3_REGION", "ap-northeast-1"),
		S3KeyPrefix:  getEnvOrDefault("S3_KEY_PREFIX", "products/"),
		BatchSize:    getIntEnvOrDefault("BATCH_SIZE", 100),
		EnableDebug:  getBoolEnvOrDefault("ENABLE_DEBUG", true),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

type KnowledgeBaseBatchProcessor struct {
	db       *sql.DB
	uploader *S3Uploader
	config   *Config
}

func NewKnowledgeBaseBatchProcessor(db *sql.DB, uploader *S3Uploader, config *Config) *KnowledgeBaseBatchProcessor {
	return &KnowledgeBaseBatchProcessor{
		db:       db,
		uploader: uploader,
		config:   config,
	}
}

func (p *KnowledgeBaseBatchProcessor) ProcessAllProducts(ctx context.Context) error {
	log.Println("Starting product data extraction...")

	// アクティブな商品のみを取得
	products, err := models.Products(
		qm.Where("is_active = ?", true),
		qm.Load("Category"),
		qm.Load("ProductReviews", qm.OrderBy("created_at DESC"), qm.Limit(10)), // 最新10件のレビュー
	).All(ctx, p.db)

	if err != nil {
		return fmt.Errorf("failed to fetch products: %w", err)
	}

	log.Printf("Found %d products to process", len(products))

	// バッチ処理でアップロード
	successCount := 0
	errorCount := 0

	for i, product := range products {
		if p.config.EnableDebug {
			log.Printf("Processing product %d/%d: %s", i+1, len(products), product.Name)
		}

		// ナレッジベース用ドキュメントを生成
		document, err := p.generateProductDocument(product)
		if err != nil {
			log.Printf("Failed to generate document for product %s: %v", product.ID, err)
			errorCount++
			continue
		}

		// S3にアップロード
		if err := p.uploader.UploadDocument(ctx, document); err != nil {
			log.Printf("Failed to upload document for product %s: %v", product.ID, err)
			errorCount++
			continue
		}

		successCount++

		// バッチサイズごとに進捗を出力
		if (i+1)%p.config.BatchSize == 0 {
			log.Printf("Processed %d/%d products (Success: %d, Error: %d)",
				i+1, len(products), successCount, errorCount)
		}
	}

	log.Printf("Batch processing completed. Success: %d, Error: %d", successCount, errorCount)
	return nil
}

func (p *KnowledgeBaseBatchProcessor) generateProductDocument(product *models.Product) (*ProductKnowledgeDocument, error) {
	var content strings.Builder

	// 基本商品情報（マークダウン形式）
	content.WriteString("## 基本情報\n\n")
	content.WriteString(fmt.Sprintf("**商品ID:** %s  \n", product.ID))

	// ブランド情報
	if product.Brand.Valid {
		content.WriteString(fmt.Sprintf("**ブランド:** %s  \n", product.Brand.String))
	}

	// カテゴリ情報
	var categoryName string
	if product.R != nil && product.R.Category != nil {
		categoryName = product.R.Category.Name
		content.WriteString(fmt.Sprintf("**カテゴリ:** %s  \n", categoryName))
		if product.R.Category.Description.Valid {
			content.WriteString(fmt.Sprintf("**カテゴリ説明:** %s  \n", product.R.Category.Description.String))
		}
	}

	// 価格情報
	content.WriteString("\n## 価格情報\n\n")
	content.WriteString(fmt.Sprintf("**現在価格:** %s円  \n", product.Price.String()))
	if !product.OriginalPrice.IsZero() {
		content.WriteString(fmt.Sprintf("**元価格:** %s円  \n", product.OriginalPrice.String()))
	}

	// 商品説明
	if product.Description.Valid {
		content.WriteString("\n## 商品説明\n\n")
		content.WriteString(fmt.Sprintf("%s\n\n", product.Description.String))
	}

	// 特徴をマークダウンリスト化
	if product.Features.Valid {
		var features []string
		if err := json.Unmarshal(product.Features.JSON, &features); err == nil && len(features) > 0 {
			content.WriteString("## 商品特徴\n\n")
			for _, feature := range features {
				content.WriteString(fmt.Sprintf("- %s\n", feature))
			}
			content.WriteString("\n")
		}
	}

	// タグ情報
	if len(product.Tags) > 0 {
		content.WriteString("## タグ\n\n")
		for _, tag := range product.Tags {
			content.WriteString(fmt.Sprintf("`%s` ", tag))
		}
		content.WriteString("\n\n")
	}

	// 評価情報
	if !product.RatingAverage.IsZero() || product.RatingCount.Valid {
		content.WriteString("## 評価情報\n\n")
		if !product.RatingAverage.IsZero() {
			content.WriteString(fmt.Sprintf("**平均評価:** %s/5 ⭐  \n", product.RatingAverage.String()))
		}
		if product.RatingCount.Valid {
			content.WriteString(fmt.Sprintf("**評価件数:** %d件  \n\n", product.RatingCount.Int))
		}
	}

	// レビュー情報
	if product.R != nil && len(product.R.ProductReviews) > 0 {
		content.WriteString("## 顧客レビュー\n\n")
		for i, review := range product.R.ProductReviews {
			if i >= 5 { // 最大5件のレビューを含める
				break
			}
			content.WriteString(fmt.Sprintf("### レビュー %d\n\n", i+1))
			content.WriteString(fmt.Sprintf("**評価:** %d/5 ⭐  \n", review.Rating))
			if review.Title.Valid {
				content.WriteString(fmt.Sprintf("**タイトル:** %s  \n", review.Title.String))
			}
			if review.Content.Valid {
				// レビュー内容が長い場合は制限
				reviewContent := review.Content.String
				if len(reviewContent) > 200 {
					reviewContent = reviewContent[:200] + "..."
				}
				content.WriteString(fmt.Sprintf("**内容:** %s\n\n", reviewContent))
			}
		}
	}

	// 在庫情報
	if product.StockQuantity.Valid {
		content.WriteString("## 在庫情報\n\n")
		if product.StockQuantity.Int > 0 {
			content.WriteString(fmt.Sprintf("**在庫数:** %d個  \n", product.StockQuantity.Int))
		} else {
			content.WriteString("**在庫状況:** 品切れ ❌  \n")
		}
		content.WriteString("\n")
	}

	// メタデータ構築
	metadata := ProductMetadata{
		CategoryID: product.CategoryID,
		Category:   categoryName,
		Price:      product.Price.String(),
		IsActive:   product.IsActive.Valid && product.IsActive.Bool,
		Tags:       product.Tags,
	}

	if product.Brand.Valid {
		metadata.Brand = product.Brand.String
	}

	if !product.OriginalPrice.IsZero() {
		metadata.OriginalPrice = product.OriginalPrice.String()
	}

	if !product.RatingAverage.IsZero() {
		metadata.RatingAverage = product.RatingAverage.String()
	}
	if product.RatingCount.Valid {
		metadata.RatingCount = product.RatingCount.Int
	}
	if product.PopularityScore.Valid {
		metadata.PopularityScore = product.PopularityScore.Int
	}

	return &ProductKnowledgeDocument{
		ProductID:   product.ID,
		ProductName: product.Name,
		Content:     content.String(),
		Metadata:    metadata,
		GeneratedAt: time.Now(),
	}, nil
}

func (u *S3Uploader) UploadDocument(ctx context.Context, document *ProductKnowledgeDocument) error {
	// マークダウン形式でドキュメントを生成
	var markdown strings.Builder

	// メタデータをフロントマターとして追加
	markdown.WriteString("---\n")
	markdown.WriteString(fmt.Sprintf("product_id: %s\n", document.ProductID))
	markdown.WriteString(fmt.Sprintf("product_name: %s\n", document.ProductName))
	markdown.WriteString(fmt.Sprintf("brand: %s\n", document.Metadata.Brand))
	markdown.WriteString(fmt.Sprintf("category: %s\n", document.Metadata.Category))
	markdown.WriteString(fmt.Sprintf("category_id: %d\n", document.Metadata.CategoryID))
	markdown.WriteString(fmt.Sprintf("price: %s\n", document.Metadata.Price))
	if document.Metadata.OriginalPrice != "" {
		markdown.WriteString(fmt.Sprintf("original_price: %s\n", document.Metadata.OriginalPrice))
	}
	if document.Metadata.RatingAverage != "" {
		markdown.WriteString(fmt.Sprintf("rating_average: %s\n", document.Metadata.RatingAverage))
	}
	if document.Metadata.RatingCount > 0 {
		markdown.WriteString(fmt.Sprintf("rating_count: %d\n", document.Metadata.RatingCount))
	}
	if document.Metadata.PopularityScore > 0 {
		markdown.WriteString(fmt.Sprintf("popularity_score: %d\n", document.Metadata.PopularityScore))
	}
	if len(document.Metadata.Tags) > 0 {
		markdown.WriteString(fmt.Sprintf("tags: [%s]\n", strings.Join(document.Metadata.Tags, ", ")))
	}
	markdown.WriteString(fmt.Sprintf("is_active: %t\n", document.Metadata.IsActive))
	markdown.WriteString(fmt.Sprintf("generated_at: %s\n", document.GeneratedAt.Format(time.RFC3339)))
	markdown.WriteString("---\n\n")

	// メインコンテンツをマークダウンとして追加
	markdown.WriteString(fmt.Sprintf("# %s\n\n", document.ProductName))
	markdown.WriteString(document.Content)

	// S3キーを生成 (日付ごとにフォルダ分け、拡張子を.mdに変更)
	datePrefix := time.Now().Format("2006/01/02")
	key := fmt.Sprintf("%s%s/product_%s.md", u.prefix, datePrefix, document.ProductID)

	// S3にアップロード
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		Body:        strings.NewReader(markdown.String()),
		ContentType: aws.String("text/markdown"),
		Metadata: map[string]string{
			"product-id":   document.ProductID,
			"product-name": document.ProductName,
			"generated-at": document.GeneratedAt.Format(time.RFC3339),
			"content-type": "knowledge-base-document",
		},
	})

	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}
