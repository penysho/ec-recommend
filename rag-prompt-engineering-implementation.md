# RAGプロンプトエンジニアリング強化実装レポート

## 概要

現在のec-recommendプロジェクトのRAGプロンプトを強化し、より精度の高い商品推薦と説得力のある理由生成を実現するための具体的な実装方法を示します。

## 現在の課題分析

### 既存実装の問題点

```go
// 現在の単純なプロンプト実装
func (rs *RecommendationServiceV2) createExplanationPrompt(recommendationID uuid.UUID, profile *dto.CustomerProfile) string {
    // 基本的なプロンプト構築のみ
    prompt := fmt.Sprintf("商品%sを%sに推薦する理由を説明してください", recommendationID, profile.CustomerID)
    return prompt
}
```

**問題点:**
- コンテキストの考慮が不十分
- Few-shot学習の活用なし
- 出力形式の一貫性不足
- 推論過程の可視化なし

## 実装方針

### 1. プロンプトテンプレートシステムの構築

```go
// internal/service/prompt_template.go
package service

import (
    "context"
    "fmt"
    "strings"
    "ec-recommend/internal/dto"
)

type PromptTemplate struct {
    ID            string
    Name          string
    BasePrompt    string
    Examples      []ExampleCase
    ContextRules  []ContextRule
    OutputSchema  OutputSchema
    Version       string
}

type ExampleCase struct {
    Input           string
    ExpectedOutput  string
    Explanation     string
}

type ContextRule struct {
    Condition   string
    Modification string
    Priority    int
}

type OutputSchema struct {
    Format      string
    Fields      []SchemaField
    Constraints []string
}

type SchemaField struct {
    Name        string
    Type        string
    Required    bool
    Description string
}
```

### 2. コンテキスト対応プロンプトジェネレーター

```go
// internal/service/prompt_generator.go
package service

type PromptGenerator struct {
    templates map[string]*PromptTemplate
    config    *PromptConfig
}

type PromptConfig struct {
    MaxTokens         int
    Temperature       float64
    EnableChainOfThought bool
    EnableFewShot     bool
    PersonalizationLevel string
}

func NewPromptGenerator(config *PromptConfig) *PromptGenerator {
    return &PromptGenerator{
        templates: make(map[string]*PromptTemplate),
        config:    config,
    }
}

func (pg *PromptGenerator) GenerateRecommendationPrompt(
    ctx context.Context,
    products []dto.ProductRecommendationV2,
    profile *dto.CustomerProfile,
    contextType string,
) (string, error) {

    // 1. テンプレート選択
    template := pg.selectTemplate(contextType, profile)

    // 2. コンテキスト情報の構築
    contextInfo := pg.buildContextInfo(profile, contextType)

    // 3. Few-shot例の選択
    examples := pg.selectRelevantExamples(template, profile, products)

    // 4. Chain-of-Thought構造の構築
    reasoning := pg.buildReasoningStructure(contextType)

    // 5. 最終プロンプトの組み立て
    prompt := pg.assemblePrompt(template, contextInfo, examples, reasoning, products)

    return prompt, nil
}
```

### 3. 強化されたプロンプトテンプレートの定義

```go
// internal/service/prompt_templates.go
package service

func (pg *PromptGenerator) initializeTemplates() {
    // ホームページ向けプロンプト
    pg.templates["homepage_recommendations"] = &PromptTemplate{
        ID:   "homepage_recommendations",
        Name: "ホームページ商品推薦",
        BasePrompt: `
あなたは経験豊富なECサイトの販売アドバイザーです。
顧客の購買履歴と嗜好を分析し、最適な商品を推薦してください。

## 顧客プロファイル
{{.CustomerProfile}}

## 推薦対象商品
{{.Products}}

## 推薦理由生成の指針
1. 顧客の過去の購入パターンを分析
2. 類似する他の顧客の行動パターンを考慮
3. 季節性やトレンドを反映
4. 具体的なベネフィットを明示
5. 感情的な訴求を含める

## 出力形式
{{.OutputSchema}}
`,
        Examples: []ExampleCase{
            {
                Input: `顧客：30代女性、過去にスキンケア商品を購入
商品：高級美容液`,
                ExpectedOutput: `{
    "recommendation_reason": "過去のスキンケア商品の購入履歴から、美容への関心が高いことが分かります。この高級美容液は、年齢に応じたアンチエイジング効果が期待でき、朝晩のスキンケアルーチンをワンランク上げることができます。",
    "confidence_score": 0.85,
    "key_benefits": ["アンチエイジング効果", "肌質改善", "高級感"],
    "usage_scenarios": ["朝のスキンケア", "夜のスペシャルケア"]
}`,
                Explanation: "購買履歴と年齢を考慮した具体的な提案",
            },
        },
    }

    // 商品詳細ページ向けプロンプト
    pg.templates["product_detail_recommendations"] = &PromptTemplate{
        ID:   "product_detail_recommendations",
        Name: "商品詳細ページ関連商品推薦",
        BasePrompt: `
あなたは商品知識豊富な販売スタッフです。
顧客が閲覧中の商品と関連性の高い商品を推薦してください。

## 現在閲覧中の商品
{{.CurrentProduct}}

## 顧客情報
{{.CustomerProfile}}

## 関連商品推薦
{{.Products}}

## 推薦戦略
1. 補完関係（一緒に使うと効果的）
2. 代替品（予算に応じた選択肢）
3. アップグレード（より高品質な選択肢）
4. クロスセル（関連カテゴリの商品）

## 推論過程を示してください
Step 1: 現在の商品の特徴分析
Step 2: 顧客ニーズの推定
Step 3: 最適な推薦理由の構築

## 出力形式
{{.OutputSchema}}
`,
    }
}
```

### 4. Chain-of-Thought推論の実装

```go
// internal/service/reasoning_chain.go
package service

type ReasoningChain struct {
    Steps []ReasoningStep
}

type ReasoningStep struct {
    StepNumber  int
    Description string
    Analysis    string
    Conclusion  string
    Confidence  float64
}

func (pg *PromptGenerator) buildReasoningStructure(contextType string) *ReasoningChain {
    switch contextType {
    case "homepage":
        return &ReasoningChain{
            Steps: []ReasoningStep{
                {
                    StepNumber:  1,
                    Description: "顧客の購買履歴分析",
                    Analysis:    "過去の購入商品から嗜好パターンを特定",
                    Conclusion:  "顧客の主要な関心領域を把握",
                    Confidence:  0.8,
                },
                {
                    StepNumber:  2,
                    Description: "商品マッチング評価",
                    Analysis:    "顧客嗜好と商品特徴の適合度を評価",
                    Conclusion:  "最適な商品候補を選定",
                    Confidence:  0.75,
                },
                {
                    StepNumber:  3,
                    Description: "推薦理由の構築",
                    Analysis:    "具体的なベネフィットと使用シーンを提示",
                    Conclusion:  "説得力のある推薦メッセージを生成",
                    Confidence:  0.85,
                },
            },
        }
    case "product_detail":
        return &ReasoningChain{
            Steps: []ReasoningStep{
                {
                    StepNumber:  1,
                    Description: "現在商品の特徴分析",
                    Analysis:    "閲覧中商品の機能、価格帯、用途を分析",
                    Conclusion:  "商品の核となる価値を特定",
                    Confidence:  0.9,
                },
                {
                    StepNumber:  2,
                    Description: "関連性評価",
                    Analysis:    "補完、代替、アップグレード関係を評価",
                    Conclusion:  "最適な関連商品タイプを決定",
                    Confidence:  0.8,
                },
                {
                    StepNumber:  3,
                    Description: "購買意欲向上戦略",
                    Analysis:    "クロスセル・アップセル機会を分析",
                    Conclusion:  "購買行動を促進する提案を構築",
                    Confidence:  0.75,
                },
            },
        }
    default:
        return &ReasoningChain{Steps: []ReasoningStep{}}
    }
}
```

### 5. 動的Few-Shot例選択システム

```go
// internal/service/example_selector.go
package service

type ExampleSelector struct {
    exampleDatabase map[string][]ExampleCase
    similarityCalculator SimilarityCalculator
}

type SimilarityCalculator interface {
    CalculateSimilarity(profile1, profile2 *dto.CustomerProfile) float64
    CalculateProductSimilarity(product1, product2 *dto.ProductRecommendationV2) float64
}

func (es *ExampleSelector) SelectRelevantExamples(
    template *PromptTemplate,
    profile *dto.CustomerProfile,
    products []dto.ProductRecommendationV2,
    maxExamples int,
) []ExampleCase {

    relevantExamples := make([]ExampleCase, 0)

    // 1. 顧客プロファイル類似度による例選択
    profileExamples := es.selectByCustomerProfile(profile, maxExamples/2)
    relevantExamples = append(relevantExamples, profileExamples...)

    // 2. 商品類似度による例選択
    productExamples := es.selectByProductSimilarity(products, maxExamples/2)
    relevantExamples = append(relevantExamples, productExamples...)

    // 3. 多様性を考慮した最終選択
    finalExamples := es.diversifyExamples(relevantExamples, maxExamples)

    return finalExamples
}

func (es *ExampleSelector) selectByCustomerProfile(
    profile *dto.CustomerProfile,
    maxCount int,
) []ExampleCase {
    examples := make([]ExampleCase, 0)

    // 年齢層での絞り込み
    if profile.Age >= 20 && profile.Age < 30 {
        examples = append(examples, es.exampleDatabase["age_20s"]...)
    } else if profile.Age >= 30 && profile.Age < 40 {
        examples = append(examples, es.exampleDatabase["age_30s"]...)
    }

    // 購買履歴での絞り込み
    for _, purchase := range profile.PurchaseHistory {
        if categoryExamples, exists := es.exampleDatabase[fmt.Sprintf("category_%d", purchase.CategoryID)]; exists {
            examples = append(examples, categoryExamples...)
        }
    }

    // 類似度でソートして上位を選択
    return es.selectTopExamples(examples, maxCount)
}
```

### 6. 出力形式標準化システム

```go
// internal/service/output_formatter.go
package service

type OutputFormatter struct {
    schemas map[string]OutputSchema
}

func NewOutputFormatter() *OutputFormatter {
    formatter := &OutputFormatter{
        schemas: make(map[string]OutputSchema),
    }
    formatter.initializeSchemas()
    return formatter
}

func (of *OutputFormatter) initializeSchemas() {
    // 商品推薦用スキーマ
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
                Description: "推薦理由（150文字以内）",
            },
            {
                Name:        "confidence_score",
                Type:        "float",
                Required:    true,
                Description: "信頼度スコア（0.0-1.0）",
            },
            {
                Name:        "key_benefits",
                Type:        "array[string]",
                Required:    true,
                Description: "主要なベネフィット（最大3つ）",
            },
            {
                Name:        "usage_scenarios",
                Type:        "array[string]",
                Required:    true,
                Description: "使用シーン（最大3つ）",
            },
            {
                Name:        "reasoning_chain",
                Type:        "array[object]",
                Required:    false,
                Description: "推論過程（Chain-of-Thought）",
            },
        },
        Constraints: []string{
            "recommendation_reasonは150文字以内",
            "confidence_scoreは0.0-1.0の範囲",
            "key_benefitsは最大3要素",
            "usage_scenariosは最大3要素",
        },
    }
}

func (of *OutputFormatter) FormatSchema(schemaName string) string {
    schema, exists := of.schemas[schemaName]
    if !exists {
        return ""
    }

    var builder strings.Builder
    builder.WriteString("以下のJSON形式で出力してください：\n\n")
    builder.WriteString("```json\n{\n")

    for i, field := range schema.Fields {
        builder.WriteString(fmt.Sprintf("  \"%s\": \"%s\"", field.Name, field.Description))
        if field.Required {
            builder.WriteString(" (必須)")
        }
        if i < len(schema.Fields)-1 {
            builder.WriteString(",")
        }
        builder.WriteString("\n")
    }

    builder.WriteString("}\n```\n\n")
    builder.WriteString("制約事項：\n")
    for _, constraint := range schema.Constraints {
        builder.WriteString(fmt.Sprintf("- %s\n", constraint))
    }

    return builder.String()
}
```

### 7. メインサービスへの統合

```go
// internal/service/recommendation_service_v2.go への追加実装

func (rs *RecommendationServiceV2) enhanceWithAdvancedAIExplanations(
    ctx context.Context,
    recommendations []dto.ProductRecommendationV2,
    profile *dto.CustomerProfile,
    contextType string,
    metrics *dto.PerformanceMetrics,
) ([]dto.ProductRecommendationV2, error) {

    startTime := time.Now()

    // 1. 強化されたプロンプト生成
    prompt, err := rs.promptGenerator.GenerateRecommendationPrompt(
        ctx, recommendations, profile, contextType,
    )
    if err != nil {
        return recommendations, fmt.Errorf("failed to generate enhanced prompt: %w", err)
    }

    // 2. AI応答の生成
    response, err := rs.chatService.GenerateResponse(ctx, prompt)
    if err != nil {
        return recommendations, fmt.Errorf("failed to generate AI response: %w", err)
    }

    // 3. 構造化された応答の解析
    enhancedRecommendations, err := rs.parseEnhancedAIResponse(
        response.Content, recommendations,
    )
    if err != nil {
        log.Printf("Warning: failed to parse enhanced AI response: %v", err)
        return recommendations, nil // フォールバック
    }

    // 4. 品質検証
    validatedRecommendations := rs.validateRecommendations(enhancedRecommendations)

    metrics.AIProcessingTimeMs = time.Since(startTime).Milliseconds()

    return validatedRecommendations, nil
}

func (rs *RecommendationServiceV2) parseEnhancedAIResponse(
    content string,
    originalRecommendations []dto.ProductRecommendationV2,
) ([]dto.ProductRecommendationV2, error) {

    // JSON形式の応答を解析
    var aiResponses []struct {
        ProductID          string            `json:"product_id"`
        RecommendationReason string         `json:"recommendation_reason"`
        ConfidenceScore    float64          `json:"confidence_score"`
        KeyBenefits        []string         `json:"key_benefits"`
        UsageScenarios     []string         `json:"usage_scenarios"`
        ReasoningChain     []ReasoningStep  `json:"reasoning_chain"`
    }

    if err := json.Unmarshal([]byte(content), &aiResponses); err != nil {
        return nil, fmt.Errorf("failed to unmarshal AI response: %w", err)
    }

    // 元の推薦リストと統合
    enhancedRecommendations := make([]dto.ProductRecommendationV2, len(originalRecommendations))
    copy(enhancedRecommendations, originalRecommendations)

    for i, rec := range enhancedRecommendations {
        for _, aiResp := range aiResponses {
            if aiResp.ProductID == rec.ProductID.String() {
                rec.Reason = aiResp.RecommendationReason
                rec.ConfidenceScore = aiResp.ConfidenceScore

                // AI インサイトの強化
                if rec.AIInsights == nil {
                    rec.AIInsights = &dto.ProductAIInsights{}
                }
                rec.AIInsights.KeyFeatures = aiResp.KeyBenefits
                rec.AIInsights.UseCases = aiResp.UsageScenarios

                enhancedRecommendations[i] = rec
                break
            }
        }
    }

    return enhancedRecommendations, nil
}
```

## 実装ステップ

### Phase 1: 基盤構築（1-2週間）
1. **プロンプトテンプレートシステム**の実装
2. **出力形式標準化**の実装
3. **基本的なコンテキスト対応**の実装

### Phase 2: 高度化（2-3週間）
1. **Few-Shot例選択システム**の実装
2. **Chain-of-Thought推論**の実装
3. **動的プロンプト生成**の実装

### Phase 3: 最適化（1-2週間）
1. **品質検証システム**の実装
2. **パフォーマンス最適化**
3. **A/Bテスト対応**

## 期待される効果

### 推薦精度の向上
- **信頼度スコア**: 15-20%向上
- **クリック率**: 10-15%向上
- **コンバージョン率**: 8-12%向上

### ユーザー体験の改善
- **説明の説得力**: より具体的で理解しやすい理由
- **個別化レベル**: 顧客固有のニーズへの対応
- **一貫性**: 統一された出力形式

### 運用効率の向上
- **品質の安定性**: 構造化された出力
- **保守性**: テンプレートベースの管理
- **拡張性**: 新しいコンテキストへの対応

## 次のステップ

1. **プロトタイプ開発**: 小規模な実装で効果を検証
2. **A/Bテスト**: 従来方式との比較評価
3. **段階的展開**: 成功したパターンの全面適用
4. **継続的改善**: ユーザーフィードバックによる最適化

この実装により、ec-recommendシステムの推薦品質と説得力を大幅に向上させることが可能になります。
