# EC推薦システム改善分析レポート

## プロジェクト概要

`ec-recommend`プロジェクトは、AWS Bedrock と RAG (Retrieval-Augmented Generation) を活用した EC サイト向けの商品推薦システムです。現在の実装では以下の機能を提供しています：

- **複数の推薦手法**: ハイブリッド、セマンティック、協調フィルタリング、ベクトル検索、知識ベース
- **RAG 機能**: AWS Bedrock Knowledge Base を使用した商品検索と推薦理由生成
- **バッチ処理**: 商品データの Knowledge Base 用データ生成とS3アップロード
- **AI 強化**: Bedrock による推薦理由の自動生成

## 現在の実装の強み

### 1. 包括的なアーキテクチャ
- **レイヤー分離**: Handler → Service → Repository の明確な分離
- **インターフェース設計**: 依存関係逆転の原則に従った設計
- **v1/v2 併存**: 段階的な機能改善を可能にする設計

### 2. 高度なRAG機能
- **メタデータフィルタリング**: 詳細な商品属性による絞り込み
- **ハイブリッド検索**: セマンティック検索とキーワード検索の組み合わせ
- **コンテキスト考慮**: 顧客プロファイルとコンテキストに基づく個別化

### 3. 充実したデータ構造
- **詳細なDTO**: 推薦結果、メタデータ、パフォーマンス指標を含む包括的な応答
- **AI インサイト**: 商品の特徴、用途、ターゲット顧客などの AI 生成情報

## 改善点の特定

### 🔴 **重要度：高**

#### 1. エンベディングモデルの選択と最適化

**現在の問題**:
```go
// 固定でTitan Embeddingsを使用
EmbeddingModelID: getEnvWithDefault("EMBEDDING_MODEL_ID", "amazon.titan-embed-text-v1")
```

**改善案**:
- **多言語対応**: `amazon.titan-embed-text-v2`への移行
- **ドメイン特化**: EC商品に特化したファインチューニング済みモデルの検討
- **動的モデル選択**: コンテンツタイプに応じたモデル切り替え

```go
type EmbeddingConfig struct {
    ProductModel  string // 商品説明用
    ReviewModel   string // レビュー分析用
    QueryModel    string // 検索クエリ用
}
```

#### 2. Knowledge Base データ品質の向上

**現在の問題**:
```go
// 単純なマークダウン形式での保存
markdown.WriteString(fmt.Sprintf("# %s\n\n", document.ProductName))
markdown.WriteString(document.Content)
```

**改善案**:
- **構造化データ**: JSONドキュメントとマークダウンのハイブリッド構造
- **セマンティック情報の強化**: 商品の用途、シーン、ベネフィットを明示
- **関連性向上**: カテゴリ間の関係性やクロスセル情報の追加

```go
type EnhancedProductDocument struct {
    ProductInfo     ProductMetadata      `json:"product_info"`
    SemanticTags    []SemanticTag       `json:"semantic_tags"`
    UseCases        []UseCase           `json:"use_cases"`
    Relationships   ProductRelationship  `json:"relationships"`
    SearchKeywords  []string            `json:"search_keywords"`
}
```

#### 3. RAG プロンプトエンジニアリングの強化

**現在の問題**:
```go
// 単純なプロンプト構築
prompt := rs.createExplanationPrompt(recommendationID, profile)
```

**改善案**:
- **Few-shot プロンプト**: 良い推薦例を含めた学習
- **Chain-of-Thought**: 段階的な推論過程の明示
- **Dynamic prompting**: コンテキストに応じたプロンプト最適化

```go
type PromptTemplate struct {
    BasePrompt    string
    Examples      []ExampleCase
    ContextRules  []ContextRule
    OutputFormat  OutputSchema
}
```

### 🟡 **重要度：中**

#### 4. キャッシュ戦略の実装

**現在の問題**:
- エンベディング生成の重複処理
- 類似クエリの再計算
- Knowledge Base 検索結果のキャッシュなし

**改善案**:
```go
type RAGCacheManager interface {
    GetEmbedding(ctx context.Context, text string) ([]float64, bool)
    SetEmbedding(ctx context.Context, text string, embedding []float64, ttl time.Duration)
    GetSearchResults(ctx context.Context, query SearchQuery) ([]RAGSearchResult, bool)
    SetSearchResults(ctx context.Context, query SearchQuery, results []RAGSearchResult, ttl time.Duration)
}
```

#### 5. 評価・監視機能の強化

**現在の問題**:
- 推薦品質の定量的評価がない
- A/B テスト機能が限定的
- パフォーマンス監視が基本的

**改善案**:
```go
type RecommendationEvaluator struct {
    Metrics     []EvaluationMetric
    ABTester    ABTestManager
    Monitor     PerformanceMonitor
}

type EvaluationMetric interface {
    Calculate(recommendations []Recommendation, userActions []UserAction) float64
    Name() string
}
```

#### 6. エラーハンドリングとフォールバック機能

**現在の問題**:
```go
// Bedrock API エラー時のフォールバック戦略が不完全
if err != nil {
    log.Printf("Warning: failed to enhance recommendations with AI explanations: %v", err)
}
```

**改善案**:
```go
type FallbackStrategy interface {
    Execute(ctx context.Context, originalReq Request) (Response, error)
    ShouldFallback(err error) bool
}

type RecommendationService struct {
    Primary     RecommendationEngine
    Fallback    RecommendationEngine
    Strategy    FallbackStrategy
    CircuitBreaker CircuitBreaker
}
```

### 🟢 **重要度：低**

#### 7. コスト最適化

**改善案**:
- **トークン使用量監視**: Bedrock API 呼び出しの最適化
- **バッチ処理効率化**: 複数商品の同時処理
- **モデル選択最適化**: タスクに応じた適切なモデル選択

#### 8. 多言語対応の準備

**改善案**:
- **国際化対応**: 推薦理由の多言語生成
- **地域特化**: 地域別の商品推薦ロジック
- **文化的適応**: 地域の購買習慣を考慮した推薦

## 実装優先度

### Phase 1: 基盤強化 (1-2ヶ月)
1. エンベディングモデルの最適化
2. Knowledge Base データ品質向上
3. キャッシュ戦略の実装

### Phase 2: 機能拡張 (2-3ヶ月)
1. RAG プロンプトエンジニアリング強化
2. 評価・監視機能の実装
3. エラーハンドリング強化

### Phase 3: 最適化 (1-2ヶ月)
1. コスト最適化
2. パフォーマンスチューニング
3. 多言語対応準備

## 具体的な実装例

### 1. 強化されたRAGプロンプト

```go
func (rs *RecommendationServiceV2) createEnhancedRecommendationPrompt(
    products []ProductRecommendationV2,
    profile *dto.CustomerProfile,
    context string,
) string {
    return fmt.Sprintf(`
あなたは優秀なECサイトの販売アドバイザーです。以下の情報を基に、顧客に最適化された商品推薦理由を生成してください。

## 顧客プロファイル
- 購買履歴: %s
- 好みのカテゴリ: %s
- 価格帯: %s
- 閲覧コンテキスト: %s

## 推薦商品
%s

## 出力形式
各商品について、以下の形式でJSON応答を生成してください：
{
  "product_id": "商品ID",
  "recommendation_reason": "具体的で説得力のある推薦理由",
  "confidence_score": 0.0-1.0の信頼度,
  "key_benefits": ["利点1", "利点2", "利点3"],
  "usage_scenarios": ["利用シーン1", "利用シーン2"]
}

## 注意事項
- 顧客の過去の購入パターンを考慮
- 具体的なベネフィットを強調
- 感情的な訴求も含める
- 150文字以内で簡潔に
`,
        rs.formatPurchaseHistory(profile.PurchaseHistory),
        rs.formatPreferredCategories(profile.PreferredCategories),
        rs.formatPriceRange(profile),
        context,
        rs.formatProductsForPrompt(products),
    )
}
```

### 2. 包括的な評価フレームワーク

```go
type RecommendationMetrics struct {
    ClickThroughRate   float64 `json:"click_through_rate"`
    ConversionRate     float64 `json:"conversion_rate"`
    DiversityScore     float64 `json:"diversity_score"`
    NoveltyScore       float64 `json:"novelty_score"`
    SerendipityScore   float64 `json:"serendipity_score"`
    RelevanceScore     float64 `json:"relevance_score"`
    PersonalizationScore float64 `json:"personalization_score"`
}

func (e *RecommendationEvaluator) EvaluateRecommendations(
    ctx context.Context,
    recommendations []ProductRecommendationV2,
    userFeedback []UserFeedback,
) (*RecommendationMetrics, error) {
    metrics := &RecommendationMetrics{}

    // CTR calculation
    metrics.ClickThroughRate = e.calculateCTR(recommendations, userFeedback)

    // Diversity calculation
    metrics.DiversityScore = e.calculateDiversity(recommendations)

    // Novelty calculation
    metrics.NoveltyScore = e.calculateNovelty(recommendations, userFeedback)

    return metrics, nil
}
```

## まとめ

現在の `ec-recommend` プロジェクトは、RAG と Bedrock を活用した堅実な基盤を持っていますが、以下の分野で大幅な改善が可能です：

1. **AI/ML の最適化**: より適切なモデル選択とプロンプト設計
2. **データ品質**: Knowledge Base の構造化と意味情報の強化
3. **運用効率**: キャッシュ、監視、フォールバック機能の実装
4. **ビジネス価値**: より正確で説得力のある推薦システム

これらの改善により、顧客満足度の向上、コンバージョン率の改善、運用コストの削減が期待できます。
