package service

// initializeTemplates initializes all prompt templates
func (pg *PromptGenerator) initializeTemplates() {
	// Default recommendation template
	pg.templates["default_recommendation"] = &PromptTemplate{
		ID:          "default_recommendation",
		Name:        "基本商品推薦",
		Category:    "recommendation",
		Version:     "1.0",
		Description: "汎用的な商品推薦プロンプト",
		BasePrompt: `あなたは経験豊富なECサイトの販売アドバイザーです。
顧客の情報と商品データを分析し、最適な商品推薦理由を生成してください。

{{.ContextInfo}}

{{.CustomerProfile}}

{{.Products}}

## 推薦理由生成の指針
1. 顧客の購買履歴と嗜好を考慮
2. 商品の特徴と顧客ニーズのマッチングを分析
3. 具体的なベネフィットを明示
4. 感情的な訴求も含める
5. 150文字以内で簡潔かつ説得力のある内容

{{.Examples}}

{{.OutputSchema}}

重要: 必ずJSON形式で出力し、すべての必須フィールドを含めてください。`,
		Examples: []ExampleCase{
			{
				Input: `顧客：30代女性、過去にスキンケア商品を購入
商品：高級美容液（15,000円）`,
				ExpectedOutput: `{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "recommendation_reason": "過去のスキンケア商品の購入履歴から美容への関心が高く、この高級美容液は年齢に応じたアンチエイジング効果が期待できます。",
  "confidence_score": 0.85,
  "key_benefits": ["アンチエイジング効果", "肌質改善", "高級感あるケア"],
  "usage_scenarios": ["朝のスキンケア", "夜のスペシャルケア", "特別な日の準備"],
  "emotional_appeal": "自分への投資として、毎日のケアをワンランク上げる特別感"
}`,
				Explanation: "顧客の購買履歴と年齢を考慮し、具体的なベネフィットと感情的な訴求を組み合わせた推薦",
				Context:     "default",
				Priority:    1,
			},
		},
	}

	// Homepage recommendations template
	pg.templates["homepage_recommendations"] = &PromptTemplate{
		ID:          "homepage_recommendations",
		Name:        "ホームページ商品推薦",
		Category:    "homepage",
		Version:     "1.0",
		Description: "ホームページ向けの魅力的で簡潔な商品推薦",
		BasePrompt: `あなたは顧客の興味を引くECサイトのマーケティング担当者です。
ホームページを訪問した顧客に対して、魅力的で行動を促す商品推薦を生成してください。

{{.ContextInfo}}

{{.CustomerProfile}}

{{.Products}}

## ホームページ推薦の特徴
1. **簡潔で魅力的**: 100文字以内で興味を引く
2. **行動促進**: クリックしたくなる表現
3. **トレンド反映**: 今話題の要素を含める
4. **個別化**: 顧客の特徴に合わせた訴求
5. **多様性**: 異なる訴求タイプで変化をつける

## 訴求タイプ
- **trending**: 今話題の商品
- **personalized**: 個人の嗜好に合わせた提案
- **seasonal**: 季節に合った商品
- **value**: お得感のある商品

{{.Examples}}

{{.OutputSchema}}

重要: 複数商品の配列形式で出力し、各商品に適切な訴求タイプを設定してください。`,
		Examples: []ExampleCase{
			{
				Input: `顧客：25歳女性、過去にファッション小物を購入
商品：トレンドバッグ（8,000円）`,
				ExpectedOutput: `[
  {
    "product_id": "123e4567-e89b-12d3-a456-426614174000",
    "recommendation_reason": "今シーズン大注目のデザイン！あなたのファッションセンスにぴったりです",
    "confidence_score": 0.88,
    "appeal_type": "trending",
    "priority_score": 0.9
  }
]`,
				Explanation: "若い女性の顧客に対してトレンド性と個人の嗜好を組み合わせた魅力的な表現",
				Context:     "homepage",
				Priority:    1,
			},
		},
	}

	// Product detail recommendations template
	pg.templates["product_detail_recommendations"] = &PromptTemplate{
		ID:          "product_detail_recommendations",
		Name:        "商品詳細ページ関連商品推薦",
		Category:    "product_detail",
		Version:     "1.0",
		Description: "商品詳細ページでの関連商品推薦に特化したプロンプト",
		BasePrompt: `あなたは商品知識豊富な販売スタッフです。
顧客が現在閲覧している商品と関連性の高い商品を推薦し、クロスセルとアップセルの機会を創出してください。

{{.ContextInfo}}

{{.CustomerProfile}}

{{.Products}}

## 関連商品推薦戦略
1. **補完関係（complement）**: 一緒に使うとより効果的
2. **代替品（alternative）**: 予算や好みに応じた選択肢
3. **アップグレード（upgrade）**: より高品質・高機能な選択肢
4. **アクセサリー（accessory）**: 関連アクセサリーや周辺商品

## 関連性の評価基準
- 機能的な関連性（使用場面、用途）
- デザイン的な関連性（スタイル、色合い）
- 価格帯の適切さ
- 顧客の購買履歴との整合性

{{.Examples}}

{{.OutputSchema}}

重要: 各商品の関係性タイプを明確に指定し、クロスセル可能性を数値で評価してください。`,
		Examples: []ExampleCase{
			{
				Input: `現在閲覧中：高級美容液
顧客：30代女性、スキンケアに関心
関連商品：保湿クリーム（補完商品）`,
				ExpectedOutput: `[
  {
    "product_id": "456e7890-e89b-12d3-a456-426614174001",
    "recommendation_reason": "この美容液と一緒に使うことで、より効果的なスキンケアが可能になります。",
    "relationship_type": "complement",
    "confidence_score": 0.82,
    "cross_sell_potential": 0.75
  }
]`,
				Explanation: "美容液との相乗効果を強調した補完商品としての推薦",
				Context:     "product_detail",
				Priority:    1,
			},
		},
	}

	// Cart recommendations template
	pg.templates["cart_recommendations"] = &PromptTemplate{
		ID:          "cart_recommendations",
		Name:        "カート商品推薦",
		Category:    "cart",
		Version:     "1.0",
		Description: "カート画面での最後のクロスセル機会を活用",
		BasePrompt: `あなたは購買意欲の高い顧客に対するクロスセルの専門家です。
カートに商品を入れた顧客に、購買プロセスを邪魔せずに価値のある追加商品を提案してください。

{{.ContextInfo}}

{{.CustomerProfile}}

## カート内商品
{{.CurrentProduct}}

{{.Products}}

## カート推薦の原則
1. **低摩擦**: 購買プロセスを邪魔しない
2. **即座の価値**: すぐに理解できるメリット
3. **適切な価格**: カート合計の20%以下推奨
4. **関連性**: カート内商品との明確な関連
5. **緊急性**: 今追加する理由を提示

{{.Examples}}

{{.OutputSchema}}`,
		Examples: []ExampleCase{
			{
				Input: `カート内：ワイヤレスイヤホン
追加提案：充電ケース`,
				ExpectedOutput: `[
  {
    "product_id": "789e1234-e89b-12d3-a456-426614174002",
    "recommendation_reason": "イヤホンと一緒にご購入で、外出先でも安心して充電できます。今なら送料も一回で済みます。",
    "relationship_type": "accessory",
    "confidence_score": 0.9,
    "cross_sell_potential": 0.85
  }
]`,
				Explanation: "実用性と経済的メリットを組み合わせた提案",
				Context:     "cart",
				Priority:    1,
			},
		},
	}

	// Search recommendations template
	pg.templates["search_recommendations"] = &PromptTemplate{
		ID:          "search_recommendations",
		Name:        "検索結果商品推薦",
		Category:    "search",
		Version:     "1.0",
		Description: "検索意図を考慮した商品推薦",
		BasePrompt: `あなたは顧客の検索意図を理解する検索エンジンの専門家です。
顧客の検索キーワードと意図を分析し、最も関連性の高い商品を推薦してください。

{{.ContextInfo}}

{{.CustomerProfile}}

{{.Products}}

## 検索推薦の重点
1. **検索意図の理解**: 何を探しているかを正確に把握
2. **関連性の高さ**: 検索キーワードとの適合度
3. **多様な選択肢**: 価格帯やブランドのバリエーション
4. **絞り込み支援**: 選択を助ける具体的な情報

{{.Examples}}

{{.OutputSchema}}`,
		Examples: []ExampleCase{
			{
				Input: `検索キーワード：「防水 スマートウォッチ」
商品：スポーツ用スマートウォッチ`,
				ExpectedOutput: `[
  {
    "product_id": "abc1234-e89b-12d3-a456-426614174003",
    "recommendation_reason": "防水機能付きでスポーツに最適。GPS搭載で運動記録も正確に測定できます。",
    "relationship_type": "complement",
    "confidence_score": 0.92,
    "cross_sell_potential": 0.7
  }
]`,
				Explanation: "検索キーワードの防水要求とスポーツ用途を的確に反映",
				Context:     "search",
				Priority:    1,
			},
		},
	}
}
