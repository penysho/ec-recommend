package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"ec-recommend/internal/repository/db/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Config struct {
	DatabaseURL  string
	S3BucketName string
	S3Region     string
	S3KeyPrefix  string
	BatchSize    int
	EnableDebug  bool
}

// CSV用の商品データ構造体
type ProductCSVRecord struct {
	ProductID       string `csv:"product_id"`
	ProductName     string `csv:"product_name"`
	Description     string `csv:"description"` // コンテンツフィールド
	Brand           string `csv:"brand"`
	Category        string `csv:"category"`
	CategoryID      string `csv:"category_id"`
	Price           string `csv:"price"`
	OriginalPrice   string `csv:"original_price"`
	RatingAverage   string `csv:"rating_average"`
	RatingCount     string `csv:"rating_count"`
	PopularityScore string `csv:"popularity_score"`
	StockQuantity   string `csv:"stock_quantity"`
	StockStatus     string `csv:"stock_status"`
	Tags            string `csv:"tags"`
	Features        string `csv:"features"`
	ReviewSummary   string `csv:"review_summary"`
	IsActive        string `csv:"is_active"`
	CreatedAt       string `csv:"created_at"`
	UpdatedAt       string `csv:"updated_at"`
}

// メタデータファイル用の構造体
type CSVMetadata struct {
	MetadataAttributes             map[string]interface{}         `json:"metadataAttributes"`
	DocumentStructureConfiguration DocumentStructureConfiguration `json:"documentStructureConfiguration"`
}

type DocumentStructureConfiguration struct {
	Type                         string                       `json:"type"`
	RecordBasedStructureMetadata RecordBasedStructureMetadata `json:"recordBasedStructureMetadata"`
}

type RecordBasedStructureMetadata struct {
	ContentFields               []FieldName                 `json:"contentFields"`
	MetadataFieldsSpecification MetadataFieldsSpecification `json:"metadataFieldsSpecification"`
}

type FieldName struct {
	FieldName string `json:"fieldName"`
}

type MetadataFieldsSpecification struct {
	FieldsToInclude []FieldName `json:"fieldsToInclude"`
}

type S3Uploader struct {
	client *s3.Client
	bucket string
	prefix string
}

func main() {
	log.Println("Starting Knowledge Base CSV batch process...")

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
	processor := NewCSVKnowledgeBaseBatchProcessor(db, uploader, config)

	if err := processor.ProcessAllProducts(ctx); err != nil {
		log.Fatalf("CSV batch process failed: %v", err)
	}

	log.Println("Knowledge Base CSV batch process completed successfully")
}

func loadConfig() *Config {
	return &Config{
		DatabaseURL:  getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5436/postgres?sslmode=disable"),
		S3BucketName: getEnvOrDefault("S3_BUCKET_NAME", "ec-recommend-knowledge-base-local"),
		S3Region:     getEnvOrDefault("S3_REGION", "ap-northeast-1"),
		S3KeyPrefix:  getEnvOrDefault("S3_KEY_PREFIX", "products-csv/"),
		BatchSize:    getIntEnvOrDefault("BATCH_SIZE", 1000), // CSVの場合は大きなバッチサイズが効率的
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

type CSVKnowledgeBaseBatchProcessor struct {
	db       *sql.DB
	uploader *S3Uploader
	config   *Config
}

func NewCSVKnowledgeBaseBatchProcessor(db *sql.DB, uploader *S3Uploader, config *Config) *CSVKnowledgeBaseBatchProcessor {
	return &CSVKnowledgeBaseBatchProcessor{
		db:       db,
		uploader: uploader,
		config:   config,
	}
}

func (p *CSVKnowledgeBaseBatchProcessor) ProcessAllProducts(ctx context.Context) error {
	log.Println("Starting CSV product data extraction...")

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

	// CSV形式でバッチ処理
	if err := p.generateAndUploadCSV(ctx, products); err != nil {
		return fmt.Errorf("failed to generate and upload CSV: %w", err)
	}

	log.Println("CSV batch processing completed successfully")
	return nil
}

func (p *CSVKnowledgeBaseBatchProcessor) generateAndUploadCSV(ctx context.Context, products []*models.Product) error {
	// CSVデータの生成
	csvRecords := make([]ProductCSVRecord, 0, len(products))

	for _, product := range products {
		record, err := p.productToCSVRecord(product)
		if err != nil {
			log.Printf("Failed to convert product %s to CSV record: %v", product.ID, err)
			continue
		}
		csvRecords = append(csvRecords, *record)
	}

	// CSVファイル名を生成（タイムスタンプ付き）
	timestamp := time.Now().Format("20060102_150405")
	csvFileName := fmt.Sprintf("products_%s.csv", timestamp)
	metadataFileName := fmt.Sprintf("%s.metadata.json", csvFileName)

	// CSVコンテンツを生成
	csvContent, err := p.generateCSVContent(csvRecords)
	if err != nil {
		return fmt.Errorf("failed to generate CSV content: %w", err)
	}

	// メタデータコンテンツを生成
	metadataContent, err := p.generateMetadataContent(len(csvRecords))
	if err != nil {
		return fmt.Errorf("failed to generate metadata content: %w", err)
	}

	// S3にアップロード
	if err := p.uploader.UploadCSVFiles(ctx, csvFileName, csvContent, metadataFileName, metadataContent); err != nil {
		return fmt.Errorf("failed to upload CSV files: %w", err)
	}

	log.Printf("Successfully uploaded CSV with %d records", len(csvRecords))
	return nil
}

func (p *CSVKnowledgeBaseBatchProcessor) productToCSVRecord(product *models.Product) (*ProductCSVRecord, error) {
	// 商品説明を生成（コンテンツフィールド）
	description := p.generateProductDescription(product)

	// カテゴリ情報の取得
	var categoryName string
	if product.R != nil && product.R.Category != nil {
		categoryName = product.R.Category.Name
	}

	// 在庫状況の判定
	stockStatus := "在庫切れ"
	if product.StockQuantity.Valid && product.StockQuantity.Int > 0 {
		stockStatus = "在庫あり"
	}

	// タグをカンマ区切り文字列に変換
	tagsStr := ""
	if len(product.Tags) > 0 {
		tagsStr = strings.Join(product.Tags, ", ")
	}

	// 特徴を文字列に変換
	featuresStr := ""
	if product.Features.Valid {
		var features []string
		if err := json.Unmarshal(product.Features.JSON, &features); err == nil {
			featuresStr = strings.Join(features, ", ")
		}
	}

	// レビューサマリーを生成
	reviewSummary := p.generateReviewSummary(product)

	// 日時の安全な変換
	createdAt := ""
	if product.CreatedAt.Valid {
		createdAt = product.CreatedAt.Time.Format(time.RFC3339)
	}

	updatedAt := ""
	if product.UpdatedAt.Valid {
		updatedAt = product.UpdatedAt.Time.Format(time.RFC3339)
	}

	record := &ProductCSVRecord{
		ProductID:       product.ID,
		ProductName:     product.Name,
		Description:     description, // これがコンテンツフィールド
		Brand:           product.Brand.String,
		Category:        categoryName,
		CategoryID:      strconv.Itoa(product.CategoryID),
		Price:           product.Price.String(),
		OriginalPrice:   product.OriginalPrice.String(),
		RatingAverage:   product.RatingAverage.String(),
		RatingCount:     strconv.Itoa(product.RatingCount.Int),
		PopularityScore: strconv.Itoa(product.PopularityScore.Int),
		StockQuantity:   strconv.Itoa(product.StockQuantity.Int),
		StockStatus:     stockStatus,
		Tags:            tagsStr,
		Features:        featuresStr,
		ReviewSummary:   reviewSummary,
		IsActive:        strconv.FormatBool(product.IsActive.Bool),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}

	return record, nil
}

func (p *CSVKnowledgeBaseBatchProcessor) generateProductDescription(product *models.Product) string {
	var desc strings.Builder

	// 基本説明
	if product.Description.Valid {
		desc.WriteString(product.Description.String)
		desc.WriteString(" ")
	}

	// 特徴の追加
	if product.Features.Valid {
		var features []string
		if err := json.Unmarshal(product.Features.JSON, &features); err == nil && len(features) > 0 {
			desc.WriteString("主な特徴: ")
			desc.WriteString(strings.Join(features, ", "))
			desc.WriteString(" ")
		}
	}

	// レビューからの情報追加
	if product.R != nil && len(product.R.ProductReviews) > 0 {
		positiveReviews := []string{}
		for _, review := range product.R.ProductReviews {
			if review.Rating >= 4 && review.Content.Valid {
				// 簡潔なレビュー内容を追加
				content := review.Content.String
				if len(content) > 100 {
					content = content[:100] + "..."
				}
				positiveReviews = append(positiveReviews, content)
				if len(positiveReviews) >= 3 { // 最大3件
					break
				}
			}
		}
		if len(positiveReviews) > 0 {
			desc.WriteString("お客様の声: ")
			desc.WriteString(strings.Join(positiveReviews, " | "))
		}
	}

	return desc.String()
}

func (p *CSVKnowledgeBaseBatchProcessor) generateReviewSummary(product *models.Product) string {
	if product.R == nil || len(product.R.ProductReviews) == 0 {
		return ""
	}

	var summary strings.Builder
	reviewCount := len(product.R.ProductReviews)

	// 評価分布の計算
	ratingCounts := make(map[int]int)
	for _, review := range product.R.ProductReviews {
		ratingCounts[review.Rating]++
	}

	summary.WriteString(fmt.Sprintf("総レビュー数: %d件", reviewCount))
	if len(ratingCounts) > 0 {
		summary.WriteString(" | 評価分布: ")
		for rating := 5; rating >= 1; rating-- {
			if count, exists := ratingCounts[rating]; exists {
				summary.WriteString(fmt.Sprintf("%d星:%d件 ", rating, count))
			}
		}
	}

	return strings.TrimSpace(summary.String())
}

func (p *CSVKnowledgeBaseBatchProcessor) generateCSVContent(records []ProductCSVRecord) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// ヘッダーを書き込み
	headers := []string{
		"product_id", "product_name", "description", "brand", "category", "category_id",
		"price", "original_price", "rating_average", "rating_count", "popularity_score",
		"stock_quantity", "stock_status", "tags", "features", "review_summary",
		"is_active", "created_at", "updated_at",
	}

	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// データ行を書き込み
	for _, record := range records {
		row := []string{
			record.ProductID, record.ProductName, record.Description, record.Brand,
			record.Category, record.CategoryID, record.Price, record.OriginalPrice,
			record.RatingAverage, record.RatingCount, record.PopularityScore,
			record.StockQuantity, record.StockStatus, record.Tags, record.Features,
			record.ReviewSummary, record.IsActive, record.CreatedAt, record.UpdatedAt,
		}

		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.String(), nil
}

func (p *CSVKnowledgeBaseBatchProcessor) generateMetadataContent(recordCount int) (string, error) {
	metadata := CSVMetadata{
		MetadataAttributes: map[string]interface{}{
			"data_source":    "ec_product_catalog",
			"last_updated":   time.Now().Format(time.RFC3339),
			"version":        "v1.0",
			"total_products": recordCount,
			"batch_id":       fmt.Sprintf("batch_%s", time.Now().Format("20060102_150405")),
			"content_type":   "ecommerce_products",
		},
		DocumentStructureConfiguration: DocumentStructureConfiguration{
			Type: "RECORD_BASED_STRUCTURE_METADATA",
			RecordBasedStructureMetadata: RecordBasedStructureMetadata{
				ContentFields: []FieldName{
					{FieldName: "description"}, // descriptionをコンテンツフィールドに指定
				},
				MetadataFieldsSpecification: MetadataFieldsSpecification{
					FieldsToInclude: []FieldName{
						{FieldName: "product_id"},
						{FieldName: "product_name"},
						{FieldName: "brand"},
						{FieldName: "category"},
						{FieldName: "category_id"},
						{FieldName: "price"},
						{FieldName: "original_price"},
						{FieldName: "rating_average"},
						{FieldName: "rating_count"},
						{FieldName: "popularity_score"},
						{FieldName: "stock_quantity"},
						{FieldName: "stock_status"},
						{FieldName: "tags"},
						{FieldName: "features"},
						{FieldName: "review_summary"},
						{FieldName: "is_active"},
						{FieldName: "created_at"},
						{FieldName: "updated_at"},
					},
				},
			},
		},
	}

	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return string(metadataJSON), nil
}

func (u *S3Uploader) UploadCSVFiles(ctx context.Context, csvFileName, csvContent, metadataFileName, metadataContent string) error {
	// 日付ベースのプレフィックスを生成
	datePrefix := time.Now().Format("2006/01/02")

	// CSVファイルのアップロード
	csvKey := fmt.Sprintf("%s%s/%s", u.prefix, datePrefix, csvFileName)
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(csvKey),
		Body:        strings.NewReader(csvContent),
		ContentType: aws.String("text/csv"),
		Metadata: map[string]string{
			"content-type": "knowledge-base-csv",
			"generated-at": time.Now().Format(time.RFC3339),
			"file-type":    "product-data",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to upload CSV file: %w", err)
	}

	// メタデータファイルのアップロード
	metadataKey := fmt.Sprintf("%s%s/%s", u.prefix, datePrefix, metadataFileName)
	_, err = u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(metadataKey),
		Body:        strings.NewReader(metadataContent),
		ContentType: aws.String("application/json"),
		Metadata: map[string]string{
			"content-type": "knowledge-base-metadata",
			"generated-at": time.Now().Format(time.RFC3339),
			"file-type":    "metadata-config",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to upload metadata file: %w", err)
	}

	log.Printf("Successfully uploaded CSV file: %s", csvKey)
	log.Printf("Successfully uploaded metadata file: %s", metadataKey)

	return nil
}
