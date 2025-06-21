package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ec-recommend/internal/dto"
)

// PromptGenerator generates context-aware prompts for AI models
type PromptGenerator struct {
	templates       map[string]*PromptTemplate
	outputFormatter *OutputFormatter
	config          *PromptConfig
}

// PromptConfig contains configuration for prompt generation
type PromptConfig struct {
	MaxTokens            int     `json:"max_tokens"`
	Temperature          float64 `json:"temperature"`
	EnableFewShot        bool    `json:"enable_few_shot"`
	MaxExamples          int     `json:"max_examples"`
	PersonalizationLevel string  `json:"personalization_level"` // "basic", "medium", "advanced"
	IncludeReasoningChain bool   `json:"include_reasoning_chain"`
}

// NewPromptGenerator creates a new prompt generator instance
func NewPromptGenerator(config *PromptConfig) *PromptGenerator {
	if config == nil {
		config = &PromptConfig{
			MaxTokens:            2000,
			Temperature:          0.7,
			EnableFewShot:        true,
			MaxExamples:          3,
			PersonalizationLevel: "medium",
			IncludeReasoningChain: false,
		}
	}

	pg := &PromptGenerator{
		templates:       make(map[string]*PromptTemplate),
		outputFormatter: NewOutputFormatter(),
		config:          config,
	}

	pg.initializeTemplates()
	return pg
}

// GenerateRecommendationPrompt generates a recommendation prompt based on context
func (pg *PromptGenerator) GenerateRecommendationPrompt(
	ctx context.Context,
	products []dto.ProductRecommendationV2,
	profile *dto.CustomerProfile,
	contextType string,
) (string, error) {

	// 1. Select appropriate template
	template, err := pg.selectTemplate(contextType, profile)
	if err != nil {
		return "", fmt.Errorf("failed to select template: %w", err)
	}

	// 2. Build context information
	contextInfo := pg.buildContextInfo(profile, contextType)

	// 3. Select relevant examples if few-shot is enabled
	var examples string
	if pg.config.EnableFewShot && len(template.Examples) > 0 {
		examples = pg.selectAndFormatExamples(template, profile, products)
	}

	// 4. Format output schema
	schemaName := pg.getSchemaNameForContext(contextType)
	outputSchema := pg.outputFormatter.FormatSchema(schemaName)

	// 5. Build variables
	variables := &PromptVariables{
		CustomerProfile: pg.formatCustomerProfile(profile),
		Products:        pg.formatProducts(products),
		OutputSchema:    outputSchema,
		Examples:        examples,
		ContextInfo:     contextInfo,
	}

	// 6. Assemble final prompt
	prompt := pg.assemblePrompt(template, variables)

	return prompt, nil
}

// selectTemplate selects the most appropriate template for the given context
func (pg *PromptGenerator) selectTemplate(contextType string, profile *dto.CustomerProfile) (*PromptTemplate, error) {
	templateID := pg.getTemplateIDForContext(contextType)

	template, exists := pg.templates[templateID]
	if !exists {
		// Fallback to default template
		template, exists = pg.templates["default_recommendation"]
		if !exists {
			return nil, fmt.Errorf("no suitable template found for context: %s", contextType)
		}
	}

	return template, nil
}

// getTemplateIDForContext maps context types to template IDs
func (pg *PromptGenerator) getTemplateIDForContext(contextType string) string {
	switch contextType {
	case "homepage":
		return "homepage_recommendations"
	case "product_page", "product_detail":
		return "product_detail_recommendations"
	case "cart":
		return "cart_recommendations"
	case "checkout":
		return "checkout_recommendations"
	case "search_results":
		return "search_recommendations"
	default:
		return "default_recommendation"
	}
}

// getSchemaNameForContext maps context types to output schema names
func (pg *PromptGenerator) getSchemaNameForContext(contextType string) string {
	switch contextType {
	case "homepage":
		return "homepage_recommendations"
	case "product_page", "product_detail":
		return "product_detail_recommendations"
	default:
		return "product_recommendation"
	}
}

// buildContextInfo builds context information string
func (pg *PromptGenerator) buildContextInfo(profile *dto.CustomerProfile, contextType string) string {
	var builder strings.Builder

	builder.WriteString("## コンテキスト情報\n\n")
	builder.WriteString(fmt.Sprintf("**閲覧状況**: %s\n", pg.getContextDescription(contextType)))
	builder.WriteString(fmt.Sprintf("**現在時刻**: %s\n", time.Now().Format("2006年01月02日 15:04")))

	// Add seasonal context
	season := pg.getCurrentSeason()
	builder.WriteString(fmt.Sprintf("**季節**: %s\n", season))

	// Add personalization level specific information
	if pg.config.PersonalizationLevel != "basic" {
		if profile != nil {
			if len(profile.PreferredCategories) > 0 {
				builder.WriteString(fmt.Sprintf("**好みのカテゴリ**: %v\n", profile.PreferredCategories))
			}

			if pg.config.PersonalizationLevel == "advanced" {
				// Add more detailed context for advanced personalization
				if len(profile.PurchaseHistory) > 0 {
					recentPurchases := pg.getRecentPurchaseCategories(profile.PurchaseHistory, 5)
					builder.WriteString(fmt.Sprintf("**最近の購入傾向**: %v\n", recentPurchases))
				}
			}
		}
	}

	builder.WriteString("\n")
	return builder.String()
}

// formatCustomerProfile formats customer profile for prompt
func (pg *PromptGenerator) formatCustomerProfile(profile *dto.CustomerProfile) string {
	if profile == nil {
		return "新規顧客（購買履歴なし）"
	}

	var builder strings.Builder
	builder.WriteString("### 顧客プロファイル\n\n")

	// Basic information
	builder.WriteString(fmt.Sprintf("**顧客タイプ**: %s\n", func() string {
		if profile.IsPremium {
			return "プレミアム顧客"
		}
		return "一般顧客"
	}()))

	if profile.TotalSpent > 0 {
		builder.WriteString(fmt.Sprintf("**累計購入金額**: %.0f円\n", profile.TotalSpent))
	}

	if profile.OrderCount > 0 {
		builder.WriteString(fmt.Sprintf("**注文回数**: %d回\n", profile.OrderCount))
	}

	// Purchase history summary
	if len(profile.PurchaseHistory) > 0 {
		builder.WriteString(fmt.Sprintf("**購買履歴**: %d件の購入歴\n", len(profile.PurchaseHistory)))

		// Recent purchases (last 3)
		recentCount := 3
		if len(profile.PurchaseHistory) < recentCount {
			recentCount = len(profile.PurchaseHistory)
		}

		builder.WriteString("**最近の購入商品**:\n")
		for i := 0; i < recentCount; i++ {
			purchase := profile.PurchaseHistory[i]
			builder.WriteString(fmt.Sprintf("- 商品ID: %s (カテゴリ: %d, 価格: %.0f円)\n", purchase.ProductID, purchase.CategoryID, purchase.Price))
		}
	} else {
		builder.WriteString("**購買履歴**: なし（新規顧客）\n")
	}

	// Preferred categories
	if len(profile.PreferredCategories) > 0 {
		builder.WriteString(fmt.Sprintf("**好みのカテゴリ**: %v\n", profile.PreferredCategories))
	}

	return builder.String()
}

// formatProducts formats products for prompt
func (pg *PromptGenerator) formatProducts(products []dto.ProductRecommendationV2) string {
	if len(products) == 0 {
		return "推薦対象商品なし"
	}

	var builder strings.Builder
	builder.WriteString("### 推薦対象商品\n\n")

	for i, product := range products {
		builder.WriteString(fmt.Sprintf("**商品 %d**\n", i+1))
		builder.WriteString(fmt.Sprintf("- ID: %s\n", product.ProductID))
		builder.WriteString(fmt.Sprintf("- 名前: %s\n", product.Name))
		builder.WriteString(fmt.Sprintf("- 価格: %.0f円\n", product.Price))
		builder.WriteString(fmt.Sprintf("- カテゴリ: %s\n", product.CategoryName))

		if product.Description != "" {
			// Limit description length
			description := product.Description
			if len(description) > 100 {
				description = description[:100] + "..."
			}
			builder.WriteString(fmt.Sprintf("- 説明: %s\n", description))
		}

		if product.RatingAverage > 0 {
			builder.WriteString(fmt.Sprintf("- 評価: %.1f/5 (%d件)\n", product.RatingAverage, product.RatingCount))
		}

		if len(product.Tags) > 0 {
			builder.WriteString(fmt.Sprintf("- タグ: %v\n", product.Tags))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// selectAndFormatExamples selects and formats relevant examples
func (pg *PromptGenerator) selectAndFormatExamples(
	template *PromptTemplate,
	profile *dto.CustomerProfile,
	products []dto.ProductRecommendationV2,
) string {
	if len(template.Examples) == 0 {
		return ""
	}

	// For Phase 1, use simple example selection
	// In Phase 2, this will be enhanced with similarity-based selection
	maxExamples := pg.config.MaxExamples
	if len(template.Examples) < maxExamples {
		maxExamples = len(template.Examples)
	}

	var builder strings.Builder
	builder.WriteString("## 参考例\n\n")
	builder.WriteString("以下の例を参考に、同様の品質と形式で推薦理由を生成してください：\n\n")

	for i := 0; i < maxExamples; i++ {
		example := template.Examples[i]
		builder.WriteString(fmt.Sprintf("### 例 %d\n", i+1))
		builder.WriteString(fmt.Sprintf("**入力**: %s\n\n", example.Input))
		builder.WriteString(fmt.Sprintf("**出力**: \n```json\n%s\n```\n\n", example.ExpectedOutput))

		if example.Explanation != "" {
			builder.WriteString(fmt.Sprintf("**解説**: %s\n\n", example.Explanation))
		}
	}

	return builder.String()
}

// assemblePrompt assembles the final prompt from template and variables
func (pg *PromptGenerator) assemblePrompt(template *PromptTemplate, variables *PromptVariables) string {
	prompt := template.BasePrompt

	// Replace template variables
	prompt = strings.ReplaceAll(prompt, "{{.CustomerProfile}}", variables.CustomerProfile)
	prompt = strings.ReplaceAll(prompt, "{{.Products}}", variables.Products)
	prompt = strings.ReplaceAll(prompt, "{{.OutputSchema}}", variables.OutputSchema)
	prompt = strings.ReplaceAll(prompt, "{{.Examples}}", variables.Examples)
	prompt = strings.ReplaceAll(prompt, "{{.ContextInfo}}", variables.ContextInfo)

	return prompt
}

// Helper methods

func (pg *PromptGenerator) getContextDescription(contextType string) string {
	switch contextType {
	case "homepage":
		return "ホームページ閲覧中"
	case "product_page", "product_detail":
		return "商品詳細ページ閲覧中"
	case "cart":
		return "カート画面"
	case "checkout":
		return "チェックアウト画面"
	case "search_results":
		return "検索結果ページ"
	default:
		return "サイト閲覧中"
	}
}

func (pg *PromptGenerator) getCurrentSeason() string {
	month := time.Now().Month()
	switch {
	case month >= 3 && month <= 5:
		return "春"
	case month >= 6 && month <= 8:
		return "夏"
	case month >= 9 && month <= 11:
		return "秋"
	default:
		return "冬"
	}
}

func (pg *PromptGenerator) getRecentPurchaseCategories(purchases []dto.PurchaseItem, limit int) []string {
	categories := make(map[string]bool)
	count := 0

	for _, purchase := range purchases {
		if count >= limit {
			break
		}

		categoryKey := fmt.Sprintf("カテゴリ%d", purchase.CategoryID)
		if !categories[categoryKey] {
			categories[categoryKey] = true
			count++
		}
	}

	result := make([]string, 0, len(categories))
	for category := range categories {
		result = append(result, category)
	}

	return result
}
