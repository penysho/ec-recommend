package service

import (
	"fmt"
	"strings"
)

// OutputFormatter handles output format standardization
type OutputFormatter struct {
	schemas map[string]OutputSchema
}

// NewOutputFormatter creates a new output formatter instance
func NewOutputFormatter() *OutputFormatter {
	formatter := &OutputFormatter{
		schemas: make(map[string]OutputSchema),
	}
	formatter.initializeSchemas()
	return formatter
}

// initializeSchemas initializes predefined output schemas
func (of *OutputFormatter) initializeSchemas() {
	// Product recommendation schema
	of.schemas["product_recommendation"] = OutputSchema{
		Format: "json",
		Fields: []SchemaField{
			{
				Name:        "product_id",
				Type:        "string",
				Required:    true,
				Description: "推薦商品のUUID",
			},
			{
				Name:        "recommendation_reason",
				Type:        "string",
				Required:    true,
				Description: "推薦理由（150文字以内、具体的で説得力のある内容）",
			},
			{
				Name:        "confidence_score",
				Type:        "float",
				Required:    true,
				Description: "信頼度スコア（0.0-1.0の範囲）",
			},
			{
				Name:        "key_benefits",
				Type:        "array[string]",
				Required:    true,
				Description: "主要なベネフィット（最大3つ、具体的な利点）",
			},
			{
				Name:        "usage_scenarios",
				Type:        "array[string]",
				Required:    true,
				Description: "使用シーン（最大3つ、実際の利用場面）",
			},
			{
				Name:        "emotional_appeal",
				Type:        "string",
				Required:    false,
				Description: "感情的な訴求ポイント",
			},
		},
		Constraints: []string{
			"recommendation_reasonは150文字以内",
			"confidence_scoreは0.0-1.0の範囲",
			"key_benefitsは最大3要素、各要素は50文字以内",
			"usage_scenariosは最大3要素、各要素は50文字以内",
		},
		Example: `{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "recommendation_reason": "過去のスキンケア商品の購入履歴から美容への関心が高く、この高級美容液は年齢に応じたアンチエイジング効果が期待できます。",
  "confidence_score": 0.85,
  "key_benefits": ["アンチエイジング効果", "肌質改善", "高級感あるケア"],
  "usage_scenarios": ["朝のスキンケア", "夜のスペシャルケア", "特別な日の準備"],
  "emotional_appeal": "自分への投資として、毎日のケアをワンランク上げる特別感"
}`,
	}

	// Homepage recommendations schema
	of.schemas["homepage_recommendations"] = OutputSchema{
		Format: "json_array",
		Fields: []SchemaField{
			{
				Name:        "product_id",
				Type:        "string",
				Required:    true,
				Description: "推薦商品のUUID",
			},
			{
				Name:        "recommendation_reason",
				Type:        "string",
				Required:    true,
				Description: "ホームページ向けの推薦理由（100文字以内、魅力的で簡潔）",
			},
			{
				Name:        "confidence_score",
				Type:        "float",
				Required:    true,
				Description: "信頼度スコア（0.0-1.0の範囲）",
			},
			{
				Name:        "appeal_type",
				Type:        "string",
				Required:    true,
				Description: "訴求タイプ（'trending', 'personalized', 'seasonal', 'value'のいずれか）",
			},
			{
				Name:        "priority_score",
				Type:        "float",
				Required:    true,
				Description: "表示優先度（0.0-1.0の範囲）",
			},
		},
		Constraints: []string{
			"recommendation_reasonは100文字以内",
			"appeal_typeは指定された値のみ",
			"priority_scoreは0.0-1.0の範囲",
		},
		Example: `[
  {
    "product_id": "123e4567-e89b-12d3-a456-426614174000",
    "recommendation_reason": "今話題のスキンケアアイテム。あなたの肌質にぴったりの美容液です。",
    "confidence_score": 0.88,
    "appeal_type": "personalized",
    "priority_score": 0.9
  }
]`,
	}

	// Product detail recommendations schema
	of.schemas["product_detail_recommendations"] = OutputSchema{
		Format: "json_array",
		Fields: []SchemaField{
			{
				Name:        "product_id",
				Type:        "string",
				Required:    true,
				Description: "関連商品のUUID",
			},
			{
				Name:        "recommendation_reason",
				Type:        "string",
				Required:    true,
				Description: "関連商品としての推薦理由（120文字以内）",
			},
			{
				Name:        "relationship_type",
				Type:        "string",
				Required:    true,
				Description: "関係性タイプ（'complement', 'alternative', 'upgrade', 'accessory'のいずれか）",
			},
			{
				Name:        "confidence_score",
				Type:        "float",
				Required:    true,
				Description: "信頼度スコア（0.0-1.0の範囲）",
			},
			{
				Name:        "cross_sell_potential",
				Type:        "float",
				Required:    true,
				Description: "クロスセル可能性（0.0-1.0の範囲）",
			},
		},
		Constraints: []string{
			"recommendation_reasonは120文字以内",
			"relationship_typeは指定された値のみ",
			"confidence_scoreは0.0-1.0の範囲",
			"cross_sell_potentialは0.0-1.0の範囲",
		},
		Example: `[
  {
    "product_id": "456e7890-e89b-12d3-a456-426614174001",
    "recommendation_reason": "この美容液と一緒に使うことで、より効果的なスキンケアが可能になります。",
    "relationship_type": "complement",
    "confidence_score": 0.82,
    "cross_sell_potential": 0.75
  }
]`,
	}
}

// GetSchema retrieves a schema by name
func (of *OutputFormatter) GetSchema(schemaName string) (OutputSchema, bool) {
	schema, exists := of.schemas[schemaName]
	return schema, exists
}

// FormatSchema formats a schema as a prompt instruction
func (of *OutputFormatter) FormatSchema(schemaName string) string {
	schema, exists := of.schemas[schemaName]
	if !exists {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("## 出力形式\n\n")
	builder.WriteString(fmt.Sprintf("以下の%s形式で正確に出力してください：\n\n", schema.Format))

	// JSON structure
	if schema.Format == "json" {
		builder.WriteString("```json\n{\n")
	} else if schema.Format == "json_array" {
		builder.WriteString("```json\n[\n  {\n")
	}

	for i, field := range schema.Fields {
		indent := "  "
		if schema.Format == "json_array" {
			indent = "    "
		}

		builder.WriteString(fmt.Sprintf("%s\"%s\": ", indent, field.Name))

		switch field.Type {
		case "string":
			builder.WriteString(fmt.Sprintf("\"%s\"", field.Description))
		case "float":
			builder.WriteString(fmt.Sprintf("0.0 // %s", field.Description))
		case "array[string]":
			builder.WriteString(fmt.Sprintf("[\"要素1\", \"要素2\"] // %s", field.Description))
		default:
			builder.WriteString(fmt.Sprintf("\"%s\"", field.Description))
		}

		if field.Required {
			builder.WriteString(" (必須)")
		}

		if i < len(schema.Fields)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}

	if schema.Format == "json" {
		builder.WriteString("}\n```\n\n")
	} else if schema.Format == "json_array" {
		builder.WriteString("  }\n]\n```\n\n")
	}

	// Constraints
	if len(schema.Constraints) > 0 {
		builder.WriteString("### 制約事項\n")
		for _, constraint := range schema.Constraints {
			builder.WriteString(fmt.Sprintf("- %s\n", constraint))
		}
		builder.WriteString("\n")
	}

	// Example
	if schema.Example != "" {
		builder.WriteString("### 出力例\n")
		builder.WriteString("```json\n")
		builder.WriteString(schema.Example)
		builder.WriteString("\n```\n\n")
	}

	return builder.String()
}

// RegisterSchema registers a new output schema
func (of *OutputFormatter) RegisterSchema(schemaName string, schema OutputSchema) {
	of.schemas[schemaName] = schema
}

// ValidateSchema validates if a schema is properly formatted
func (of *OutputFormatter) ValidateSchema(schema OutputSchema) error {
	if schema.Format == "" {
		return fmt.Errorf("schema format cannot be empty")
	}

	if len(schema.Fields) == 0 {
		return fmt.Errorf("schema must have at least one field")
	}

	for _, field := range schema.Fields {
		if field.Name == "" {
			return fmt.Errorf("field name cannot be empty")
		}
		if field.Type == "" {
			return fmt.Errorf("field type cannot be empty for field %s", field.Name)
		}
	}

	return nil
}

// GetAvailableSchemas returns all available schema names
func (of *OutputFormatter) GetAvailableSchemas() []string {
	schemas := make([]string, 0, len(of.schemas))
	for name := range of.schemas {
		schemas = append(schemas, name)
	}
	return schemas
}
