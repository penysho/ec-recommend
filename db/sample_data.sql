-- Extended sample data for EC recommendation system testing

-- Categories (keep existing base categories)
INSERT INTO categories (name, description) VALUES
('electronics', '電子機器・デバイス'),
('fashion', 'ファッション・アパレル'),
('home', '家具・インテリア'),
('books', '書籍・雑誌'),
('sports', 'スポーツ・アウトドア'),
('beauty', '美容・コスメ'),
('food', '食品・飲料'),
('toys', 'おもちゃ・ゲーム'),
('health', 'ヘルス・医療'),
('automotive', '自動車・バイク');

-- Extended sub-categories
INSERT INTO categories (name, parent_id, description) VALUES
-- Electronics (id=1)
('smartphones', 1, 'スマートフォン'),
('laptops', 1, 'ノートパソコン'),
('headphones', 1, 'ヘッドフォン・イヤホン'),
('tablets', 1, 'タブレット'),
('smartwatch', 1, 'スマートウォッチ'),
('cameras', 1, 'カメラ・撮影機器'),
('gaming_console', 1, 'ゲーム機'),

-- Fashion (id=2)
('mens_fashion', 2, 'メンズファッション'),
('womens_fashion', 2, 'レディースファッション'),
('shoes', 2, '靴・シューズ'),
('accessories', 2, 'アクセサリー'),
('bags', 2, 'バッグ・かばん'),

-- Home (id=3)
('furniture', 3, '家具'),
('kitchenware', 3, 'キッチン用品'),
('bedding', 3, '寝具'),
('lighting', 3, '照明器具'),
('storage', 3, '収納用品'),

-- Books (id=4)
('fiction', 4, '小説・文学'),
('business', 4, 'ビジネス書'),
('technical', 4, '技術書'),
('manga', 4, '漫画'),
('magazine', 4, '雑誌'),

-- Sports (id=5)
('running', 5, 'ランニング'),
('fitness', 5, 'フィットネス'),
('outdoor', 5, 'アウトドア'),
('team_sports', 5, 'チームスポーツ'),
('water_sports', 5, 'ウォータースポーツ'),

-- Beauty (id=6)
('skincare', 6, 'スキンケア'),
('makeup', 6, 'メイクアップ'),
('hair_care', 6, 'ヘアケア'),
('fragrance', 6, '香水・フレグランス'),

-- Food (id=7)
('organic_food', 7, 'オーガニック食品'),
('beverages', 7, '飲料'),
('snacks', 7, 'スナック・お菓子'),
('supplements', 7, 'サプリメント'),

-- Toys (id=8)
('board_games', 8, 'ボードゲーム'),
('video_games', 8, 'ビデオゲーム'),
('educational_toys', 8, '知育玩具'),
('action_figures', 8, 'フィギュア'),

-- Health (id=9)
('medical_devices', 9, '医療機器'),
('wellness', 9, 'ウェルネス'),
('vitamins', 9, 'ビタミン・栄養'),

-- Automotive (id=10)
('car_accessories', 10, 'カー用品'),
('maintenance', 10, 'メンテナンス用品');

-- Extended products (50+ products)
INSERT INTO products (name, description, category_id, price, original_price, brand, sku, stock_quantity, features, tags, rating_average, rating_count, popularity_score) VALUES
-- Electronics
('iPhone 15 Pro', '最新のiPhone Pro シリーズ', 11, 159800, 179800, 'Apple', 'IPH15P-128', 50, '{"storage": "128GB", "color": "Natural Titanium", "camera": "48MP"}', '{"smartphone", "ios", "premium"}', 4.5, 120, 95),
('iPhone 15', 'スタンダードモデル', 11, 124800, 139800, 'Apple', 'IPH15-128', 80, '{"storage": "128GB", "color": "Blue", "camera": "48MP"}', '{"smartphone", "ios", "standard"}', 4.3, 95, 88),
('Galaxy S24', 'Samsung最新フラグシップ', 11, 149800, 169800, 'Samsung', 'GXS24-256', 60, '{"storage": "256GB", "color": "Titanium Gray", "camera": "50MP"}', '{"smartphone", "android", "premium"}', 4.4, 87, 82),
('Pixel 8', 'Google純正Android', 11, 112800, 129800, 'Google', 'PXL8-128', 70, '{"storage": "128GB", "color": "Hazel", "camera": "50MP"}', '{"smartphone", "android", "google"}', 4.2, 76, 79),

('MacBook Air M3', '軽量高性能ノートパソコン', 12, 164800, 184800, 'Apple', 'MBA-M3-13', 30, '{"cpu": "M3", "ram": "8GB", "storage": "256GB"}', '{"laptop", "mac", "portable"}', 4.7, 85, 88),
('MacBook Pro M3', 'プロ仕様ノートパソコン', 12, 248800, 268800, 'Apple', 'MBP-M3-14', 25, '{"cpu": "M3 Pro", "ram": "18GB", "storage": "512GB"}', '{"laptop", "mac", "professional"}', 4.8, 92, 91),
('ThinkPad X1', 'ビジネス向け高性能ノート', 12, 198800, 218800, 'Lenovo', 'TPX1-I7-16', 40, '{"cpu": "Intel i7", "ram": "16GB", "storage": "512GB"}', '{"laptop", "windows", "business"}', 4.5, 78, 85),
('Surface Laptop 5', 'Microsoft純正ノート', 12, 158800, 178800, 'Microsoft', 'SFL5-I5-8', 35, '{"cpu": "Intel i5", "ram": "8GB", "storage": "256GB"}', '{"laptop", "windows", "microsoft"}', 4.3, 69, 77),

('Sony WH-1000XM5', 'ノイズキャンセリングヘッドフォン', 13, 49500, 54500, 'Sony', 'WH1000XM5-B', 75, '{"noise_cancelling": true, "battery_life": "30h", "wireless": true}', '{"headphones", "wireless", "noise-cancelling"}', 4.6, 203, 92),
('AirPods Pro 2', 'Apple純正ワイヤレスイヤホン', 13, 39800, 44800, 'Apple', 'APP2-USB-C', 90, '{"noise_cancelling": true, "battery_life": "6h", "wireless": true}', '{"earphones", "wireless", "apple"}', 4.5, 156, 89),
('Bose QC45', 'プレミアムノイズキャンセリング', 13, 42800, 47800, 'Bose', 'QC45-BLACK', 65, '{"noise_cancelling": true, "battery_life": "24h", "wireless": true}', '{"headphones", "wireless", "premium"}', 4.4, 134, 86),

('iPad Air', 'バランス型タブレット', 14, 92800, 102800, 'Apple', 'IPA-M2-64', 45, '{"cpu": "M2", "storage": "64GB", "screen": "10.9inch"}', '{"tablet", "apple", "portable"}', 4.4, 98, 83),
('Surface Pro 9', 'ノートPC代替タブレット', 14, 148800, 168800, 'Microsoft', 'SPR9-I5-8', 30, '{"cpu": "Intel i5", "ram": "8GB", "storage": "256GB"}', '{"tablet", "windows", "2in1"}', 4.3, 74, 79),

('Apple Watch Series 9', '最新スマートウォッチ', 15, 59800, 69800, 'Apple', 'AWS9-GPS-41', 55, '{"gps": true, "health_monitoring": true, "battery_life": "18h"}', '{"smartwatch", "health", "apple"}', 4.5, 167, 87),
('Galaxy Watch 6', 'Samsung最新ウォッチ', 15, 47800, 52800, 'Samsung', 'GW6-BT-40', 60, '{"bluetooth": true, "health_monitoring": true, "battery_life": "40h"}', '{"smartwatch", "health", "samsung"}', 4.3, 89, 81),

-- Fashion
('プレミアムダウンジャケット', '高品質ダウンジャケット', 18, 29800, 39800, 'UNIQLO', 'PDJ-M-BK-L', 100, '{"material": "90% down", "water_resistant": true, "size": "L"}', '{"outerwear", "winter", "mens"}', 4.3, 156, 78),
('カシミヤコート', 'エレガントなカシミヤコート', 18, 89800, 109800, 'ZARA', 'CSH-CT-BG-M', 25, '{"material": "100% cashmere", "color": "beige", "size": "M"}', '{"outerwear", "luxury", "mens"}', 4.6, 67, 85),
('デニムジャケット', 'カジュアルデニムジャケット', 18, 12800, 15800, 'LEVIS', 'DNM-JK-BL-L', 80, '{"material": "cotton", "color": "blue", "size": "L"}', '{"outerwear", "casual", "mens"}', 4.2, 134, 72),

('シルクブラウス', 'エレガントなシルクブラウス', 19, 15800, 19800, 'ZARA', 'SLK-BLS-WH-M', 45, '{"material": "100% silk", "color": "white", "size": "M"}', '{"blouse", "silk", "elegant", "womens"}', 4.4, 89, 71),
('ニットワンピース', '暖かいニットワンピース', 19, 24800, 29800, 'H&M', 'KNT-OP-GR-S', 60, '{"material": "wool blend", "color": "gray", "size": "S"}', '{"dress", "winter", "womens"}', 4.3, 76, 69),
('リネンシャツ', '涼しいリネンシャツ', 19, 8800, 11800, 'MUJI', 'LNN-SH-WH-M', 90, '{"material": "100% linen", "color": "white", "size": "M"}', '{"shirt", "summer", "womens"}', 4.1, 112, 65),

('ランニングシューズ', '軽量高性能ランニングシューズ', 20, 12800, 16800, 'Nike', 'NKE-RUN-001-27', 60, '{"size": "27.0cm", "weight": "250g", "cushioning": "Air Zoom"}', '{"running", "shoes", "lightweight"}', 4.4, 187, 89),
('ビジネスシューズ', '本革ビジネスシューズ', 20, 28800, 35800, 'Cole Haan', 'CH-BIZ-BK-27', 40, '{"material": "genuine leather", "color": "black", "size": "27.0cm"}', '{"business", "leather", "formal"}', 4.5, 98, 82),
('スニーカー', 'カジュアルスニーカー', 20, 9800, 12800, 'Adidas', 'ADI-SNK-WH-26', 85, '{"color": "white", "size": "26.5cm", "style": "casual"}', '{"casual", "sneakers", "daily"}', 4.2, 156, 75),

-- Home & Kitchen
('北欧スタイルソファ', 'モダンな北欧デザインソファ', 23, 89800, 109800, 'IKEA', 'NRD-SFA-GY-3P', 12, '{"seats": 3, "material": "fabric", "color": "gray"}', '{"sofa", "nordic", "furniture"}', 4.2, 67, 65),
('エルゴノミックチェア', '人間工学に基づいたオフィスチェア', 23, 68800, 78800, 'Herman Miller', 'HM-ERG-BK', 20, '{"adjustable": true, "lumbar_support": true, "color": "black"}', '{"chair", "ergonomic", "office"}', 4.7, 89, 88),
('ダイニングテーブル', '無垢材ダイニングテーブル', 23, 128800, 148800, 'Karimoku', 'KRM-DT-OAK-4', 8, '{"material": "oak", "seats": 4, "finish": "natural"}', '{"table", "dining", "wood"}', 4.6, 45, 79),

('ステンレス製鍋セット', '高品質調理器具セット', 24, 25800, 32800, 'T-fal', 'SS-POT-SET-5', 80, '{"pieces": 5, "material": "stainless_steel", "induction_compatible": true}', '{"cookware", "kitchen", "stainless"}', 4.5, 142, 83),
('ブレンダー', '高性能ミキサーブレンダー', 24, 18800, 23800, 'Vitamix', 'VTX-BLD-001', 35, '{"power": "1200W", "capacity": "2L", "speed_settings": 10}', '{"blender", "kitchen", "healthy"}', 4.4, 76, 78),
('コーヒーメーカー', '全自動コーヒーメーカー', 24, 45800, 55800, 'Nespresso', 'NSP-CM-BK', 25, '{"type": "automatic", "milk_frother": true, "color": "black"}', '{"coffee", "automatic", "premium"}', 4.3, 134, 81),

-- Books
('プログラミング入門', 'Go言語でWebアプリケーション開発', 27, 3200, null, '技術評論社', 'BK-GO-WEB-001', 200, '{"pages": 400, "level": "beginner", "language": "japanese"}', '{"programming", "go", "web", "beginner"}', 4.1, 78, 72),
('AI機械学習実践', 'PythonではじめるAI開発', 27, 4200, null, 'オライリー', 'BK-AI-PY-001', 150, '{"pages": 520, "level": "intermediate", "language": "japanese"}', '{"ai", "python", "machine_learning"}', 4.3, 89, 85),
('クラウド設計パターン', 'AWS実践アーキテクチャ', 27, 3800, null, '翔泳社', 'BK-AWS-ARC-001', 120, '{"pages": 360, "level": "advanced", "language": "japanese"}', '{"aws", "cloud", "architecture"}', 4.4, 67, 81),

('投資の心理学', '行動経済学から学ぶ投資術', 26, 2800, null, '日経BP', 'BK-INV-PSY-001', 150, '{"pages": 320, "genre": "finance", "language": "japanese"}', '{"investment", "psychology", "finance"}', 4.3, 94, 68),
('マーケティング戦略', 'デジタル時代のブランド構築', 26, 3200, null, 'ダイヤモンド社', 'BK-MKT-DIG-001', 180, '{"pages": 280, "genre": "marketing", "language": "japanese"}', '{"marketing", "digital", "strategy"}', 4.2, 76, 74),

('鬼滅の刃 全巻セット', '人気漫画コンプリートセット', 29, 12800, 15800, '集英社', 'MG-KMT-SET-23', 50, '{"volumes": 23, "genre": "action", "complete_set": true}', '{"manga", "action", "complete"}', 4.8, 456, 95),
('ワンピース 最新巻', '冒険漫画最新巻', 29, 550, null, '集英社', 'MG-OP-108', 200, '{"volume": 108, "genre": "adventure", "series": "ongoing"}', '{"manga", "adventure", "popular"}', 4.7, 234, 89),

-- Sports & Fitness
('ヨガマット', 'プレミアムヨガマット', 32, 6800, 8800, 'Manduka', 'YGA-MAT-PRP-6MM', 40, '{"thickness": "6mm", "material": "PVC", "color": "purple"}', '{"yoga", "fitness", "mat"}', 4.6, 156, 76),
('ダンベルセット', '可変重量ダンベル', 32, 28800, 35800, 'Bowflex', 'BWF-DB-24KG', 25, '{"weight_range": "2-24kg", "adjustable": true, "space_saving": true}', '{"weights", "strength", "adjustable"}', 4.5, 98, 84),
('エクササイズバイク', '室内用フィットネスバイク', 32, 89800, 109800, 'Peloton', 'PLT-BK-001', 15, '{"resistance_levels": 32, "bluetooth": true, "display": "22inch"}', '{"cardio", "indoor", "connected"}', 4.4, 67, 79),

('テント', '4人用キャンプテント', 33, 38800, 45800, 'Coleman', 'CLM-TNT-4P', 30, '{"capacity": 4, "waterproof": true, "easy_setup": true}', '{"camping", "outdoor", "family"}', 4.3, 89, 77),
('登山リュック', '大容量ハイキングバックパック', 33, 24800, 29800, 'Osprey', 'OSP-BP-50L', 40, '{"capacity": "50L", "waterproof": true, "ergonomic": true}', '{"hiking", "backpack", "outdoor"}', 4.6, 76, 82),

-- Beauty
('保湿クリーム', 'オーガニック保湿クリーム', 35, 4800, 5800, 'SK-II', 'SKII-MOIST-50ML', 90, '{"volume": "50ml", "organic": true, "skin_type": "all"}', '{"skincare", "organic", "moisturizer"}', 4.7, 234, 91),
('美容液', 'ヒアルロン酸配合美容液', 35, 8800, 11800, 'Shiseido', 'SHI-SER-30ML', 60, '{"volume": "30ml", "ingredient": "hyaluronic_acid", "anti_aging": true}', '{"skincare", "serum", "anti_aging"}', 4.6, 167, 87),
('日焼け止め', 'SPF50+ PA++++', 35, 2800, null, 'Anessa', 'ANS-SS-60ML', 120, '{"volume": "60ml", "spf": 50, "waterproof": true}', '{"skincare", "sunscreen", "protection"}', 4.4, 198, 83),

('リップスティック', 'マットフィニッシュリップ', 36, 3200, null, 'MAC', 'MAC-LIP-RED-01', 120, '{"color": "classic_red", "finish": "matte", "long_lasting": true}', '{"makeup", "lipstick", "matte"}', 4.5, 167, 82),
('アイシャドウパレット', '12色アイシャドウ', 36, 6800, 8800, 'Urban Decay', 'UD-ES-12COL', 80, '{"colors": 12, "finish": "mixed", "pigmented": true}', '{"makeup", "eyeshadow", "palette"}', 4.4, 145, 79),
('ファンデーション', 'リキッドファンデーション', 36, 5800, null, 'NARS', 'NRS-FD-30ML', 100, '{"volume": "30ml", "coverage": "medium", "finish": "natural"}', '{"makeup", "foundation", "natural"}', 4.3, 123, 76);

-- Extended customers (20 customers)
INSERT INTO customers (email, first_name, last_name, phone, date_of_birth, gender, preferred_categories, price_range_min, price_range_max, preferred_brands, location, lifestyle_tags) VALUES
('tanaka@example.com', '太郎', '田中', '090-1234-5678', '1990-05-15', 'male', '{1,5,8}', 1000, 50000, '{"Apple", "Nike", "Sony"}', '{"prefecture": "Tokyo", "city": "Shibuya"}', '{"tech_enthusiast", "fitness", "urban"}'),
('sato@example.com', '花子', '佐藤', '090-2345-6789', '1985-08-22', 'female', '{2,6,7}', 2000, 30000, '{"ZARA", "SK-II", "UNIQLO"}', '{"prefecture": "Osaka", "city": "Namba"}', '{"fashion_lover", "beauty", "organic"}'),
('yamada@example.com', '次郎', '山田', '090-3456-7890', '1993-12-10', 'male', '{4,5}', 500, 20000, '{"Nike", "Manduka"}', '{"prefecture": "Kanagawa", "city": "Yokohama"}', '{"bookworm", "fitness", "minimalist"}'),
('watanabe@example.com', '美咲', '渡辺', '090-4567-8901', '1988-03-28', 'female', '{3,6}', 3000, 100000, '{"IKEA", "MAC", "T-fal"}', '{"prefecture": "Tokyo", "city": "Setagaya"}', '{"home_designer", "cooking", "premium"}'),
('suzuki@example.com', '健太', '鈴木', '090-5678-9012', '1995-07-03', 'male', '{1,4,8}', 1500, 200000, '{"Apple", "Sony", "Nintendo"}', '{"prefecture": "Aichi", "city": "Nagoya"}', '{"gamer", "tech", "student"}'),

-- Additional customers
('takahashi@example.com', '愛子', '高橋', '090-6789-0123', '1992-11-18', 'female', '{2,6,9}', 2500, 80000, '{"Shiseido", "MAC", "ZARA"}', '{"prefecture": "Tokyo", "city": "Ginza"}', '{"beauty_enthusiast", "fashion", "professional"}'),
('ito@example.com', '博', '伊藤', '090-7890-1234', '1987-04-25', 'male', '{1,10,5}', 5000, 150000, '{"Apple", "BMW", "Nike"}', '{"prefecture": "Kanagawa", "city": "Kawasaki"}', '{"tech_professional", "car_enthusiast", "active"}'),
('kobayashi@example.com', '麻衣', '小林', '090-8901-2345', '1991-09-12', 'female', '{3,7,4}', 1000, 40000, '{"MUJI", "IKEA", "オーガニック"}', '{"prefecture": "Kyoto", "city": "Central"}', '{"minimalist", "health_conscious", "reader"}'),
('nakamura@example.com', '大輔', '中村', '090-9012-3456', '1989-12-03', 'male', '{5,4,8}', 2000, 60000, '{"Coleman", "Osprey", "Nike"}', '{"prefecture": "Nagano", "city": "Matsumoto"}', '{"outdoorsman", "adventure", "nature_lover"}'),
('kato@example.com', '由美', '加藤', '090-0123-4567', '1994-07-30', 'female', '{6,2,7}', 3000, 70000, '{"SK-II", "H&M", "オーガニック"}', '{"prefecture": "Fukuoka", "city": "Tenjin"}', '{"beauty_lover", "trendy", "health_conscious"}'),

('hayashi@example.com', '竜也', '林', '090-1234-5679', '1986-02-14', 'male', '{1,4,5}', 8000, 200000, '{"Apple", "Sony", "Nike"}', '{"prefecture": "Tokyo", "city": "Roppongi"}', '{"tech_executive", "fitness", "luxury"}'),
('mori@example.com', '千尋', '森', '090-2345-6780', '1993-06-21', 'female', '{3,6,2}', 1500, 50000, '{"IKEA", "MUJI", "UNIQLO"}', '{"prefecture": "Hiroshima", "city": "Central"}', '{"homemaker", "minimalist", "practical"}'),
('yoshida@example.com', '和也', '吉田', '090-3456-7891', '1990-10-08', 'male', '{8,4,1}', 1000, 30000, '{"Nintendo", "Sony", "集英社"}', '{"prefecture": "Osaka", "city": "Den Den Town"}', '{"gamer", "anime_fan", "tech"}'),
('matsumoto@example.com', '理恵', '松本', '090-4567-8902', '1985-01-16', 'female', '{9,6,7}', 4000, 120000, '{"Shiseido", "オーガニック", "サプリメント"}', '{"prefecture": "Tokyo", "city": "Daikanyama"}', '{"health_expert", "premium", "wellness"}'),
('inoue@example.com', '信一', '井上', '090-5678-9013', '1988-08-27', 'male', '{10,1,5}', 10000, 300000, '{"BMW", "Apple", "Nike"}', '{"prefecture": "Kanagawa", "city": "Yokohama"}', '{"luxury_lover", "car_enthusiast", "premium"}'),

('kimura@example.com', '涼子', '木村', '090-6789-0124', '1992-03-19', 'female', '{2,6,3}', 2000, 60000, '{"ZARA", "MAC", "H&M"}', '{"prefecture": "Tokyo", "city": "Harajuku"}', '{"fashionista", "trendy", "social"}'),
('yamazaki@example.com', '健一', '山崎', '090-7890-1235', '1991-05-07', 'male', '{4,5,1}', 1500, 45000, '{"技術評論社", "Nike", "Apple"}', '{"prefecture": "Kyoto", "city": "University Area"}', '{"engineer", "fitness", "intellectual"}'),
('sasaki@example.com', '美穂', '佐々木', '090-8901-2346', '1989-11-23', 'female', '{7,9,6}', 3500, 90000, '{"オーガニック", "SK-II", "サプリメント"}', '{"prefecture": "Kanagawa", "city": "Kamakura"}', '{"health_guru", "organic", "premium"}'),
('ishii@example.com', '翔太', '石井', '090-9012-3457', '1995-09-15', 'male', '{5,8,1}', 500, 25000, '{"Nike", "Nintendo", "Sony"}', '{"prefecture": "Chiba", "city": "Funabashi"}', '{"student", "sports", "gamer"}'),
('fujita@example.com', '恵美', '藤田', '090-0123-4568', '1987-12-11', 'female', '{3,4,6}', 2500, 80000, '{"IKEA", "集英社", "Shiseido"}', '{"prefecture": "Osaka", "city": "Umeda"}', '{"working_mom", "practical", "beauty"}');

-- Generate more realistic order data with multiple orders per customer
WITH customer_orders AS (
    SELECT
        c.id as customer_id,
        c.email,
        generate_series(1, CASE
            WHEN c.email LIKE '%tanaka%' THEN 8
            WHEN c.email LIKE '%sato%' THEN 5
            WHEN c.email LIKE '%yamada%' THEN 3
            WHEN c.email LIKE '%watanabe%' THEN 6
            WHEN c.email LIKE '%suzuki%' THEN 7
            WHEN c.email LIKE '%takahashi%' THEN 4
            WHEN c.email LIKE '%ito%' THEN 6
            WHEN c.email LIKE '%kobayashi%' THEN 3
            WHEN c.email LIKE '%nakamura%' THEN 5
            WHEN c.email LIKE '%kato%' THEN 4
            WHEN c.email LIKE '%hayashi%' THEN 9
            WHEN c.email LIKE '%mori%' THEN 3
            WHEN c.email LIKE '%yoshida%' THEN 6
            WHEN c.email LIKE '%matsumoto%' THEN 7
            WHEN c.email LIKE '%inoue%' THEN 8
            WHEN c.email LIKE '%kimura%' THEN 5
            WHEN c.email LIKE '%yamazaki%' THEN 4
            WHEN c.email LIKE '%sasaki%' THEN 5
            WHEN c.email LIKE '%ishii%' THEN 2
            ELSE 3
        END) as order_number
    FROM customers c
),
order_data AS (
    INSERT INTO orders (customer_id, order_number, status, subtotal, tax_amount, shipping_fee, total_amount, payment_method, ordered_at, delivered_at)
    SELECT
        co.customer_id,
        'ORD-' || LPAD((ROW_NUMBER() OVER())::text, 5, '0'),
        CASE WHEN RANDOM() < 0.9 THEN 'delivered' ELSE 'shipped' END,
        (RANDOM() * 80000 + 5000)::numeric(10,2),
        ((RANDOM() * 80000 + 5000) * 0.1)::numeric(10,2),
        CASE WHEN RANDOM() < 0.3 THEN 500 ELSE 0 END,
        ((RANDOM() * 80000 + 5000) * 1.1 + CASE WHEN RANDOM() < 0.3 THEN 500 ELSE 0 END)::numeric(10,2),
        CASE
            WHEN RANDOM() < 0.7 THEN 'credit_card'
            WHEN RANDOM() < 0.9 THEN 'bank_transfer'
            ELSE 'mobile_payment'
        END,
        CURRENT_TIMESTAMP - INTERVAL '365 days' + (RANDOM() * INTERVAL '365 days'),
        CURRENT_TIMESTAMP - INTERVAL '360 days' + (RANDOM() * INTERVAL '360 days')
    FROM customer_orders co
    RETURNING id, customer_id, order_number
)
-- Generate order items for each order
INSERT INTO order_items (order_id, product_id, quantity, unit_price, total_price)
SELECT
    od.id,
    p.id,
    CASE WHEN RANDOM() < 0.8 THEN 1 ELSE 2 END,
    p.price,
    p.price * CASE WHEN RANDOM() < 0.8 THEN 1 ELSE 2 END
FROM order_data od
CROSS JOIN LATERAL (
    SELECT id, price
    FROM products
    ORDER BY RANDOM()
    LIMIT CASE WHEN RANDOM() < 0.6 THEN 1 WHEN RANDOM() < 0.9 THEN 2 ELSE 3 END
) p;

-- Update customer totals based on actual orders
UPDATE customers SET
    total_spent = COALESCE(order_totals.total, 0),
    order_count = COALESCE(order_totals.count, 0),
    is_premium = COALESCE(order_totals.total, 0) > 100000
FROM (
    SELECT
        o.customer_id,
        SUM(o.total_amount) as total,
        COUNT(*) as count
    FROM orders o
    GROUP BY o.customer_id
) order_totals
WHERE customers.id = order_totals.customer_id;

-- Generate extensive product reviews
INSERT INTO product_reviews (product_id, customer_id, order_id, rating, title, content, is_verified_purchase)
SELECT
    oi.product_id,
    o.customer_id,
    o.id,
    CASE
        WHEN RANDOM() < 0.4 THEN 5
        WHEN RANDOM() < 0.7 THEN 4
        WHEN RANDOM() < 0.9 THEN 3
        WHEN RANDOM() < 0.97 THEN 2
        ELSE 1
    END,
    CASE
        WHEN RANDOM() < 0.3 THEN '満足しています'
        WHEN RANDOM() < 0.5 THEN '良い商品です'
        WHEN RANDOM() < 0.7 THEN '期待通りでした'
        WHEN RANDOM() < 0.85 THEN '使いやすいです'
        ELSE 'おすすめです'
    END,
    CASE
        WHEN RANDOM() < 0.3 THEN '商品の品質が高く、期待を上回る性能でした。'
        WHEN RANDOM() < 0.5 THEN 'デザインも機能も満足しています。'
        WHEN RANDOM() < 0.7 THEN '価格に見合った良い商品だと思います。'
        WHEN RANDOM() < 0.85 THEN '使い心地が良く、日常的に活用しています。'
        ELSE '友人にもおすすめしたい商品です。'
    END,
    true
FROM orders o
JOIN order_items oi ON o.id = oi.order_id
WHERE RANDOM() < 0.6 AND o.status = 'delivered';

-- Generate extensive customer activities
INSERT INTO customer_activities (customer_id, activity_type, product_id, search_query, session_id, created_at)
SELECT
    c.id,
    CASE
        WHEN RANDOM() < 0.5 THEN 'view'
        WHEN RANDOM() < 0.7 THEN 'search'
        WHEN RANDOM() < 0.85 THEN 'add_to_cart'
        WHEN RANDOM() < 0.95 THEN 'wishlist_add'
        ELSE 'remove_from_cart'
    END,
    CASE WHEN RANDOM() < 0.8 THEN p.id ELSE NULL END,
    CASE
        WHEN RANDOM() < 0.2 THEN NULL
        WHEN RANDOM() < 0.3 THEN 'iPhone'
        WHEN RANDOM() < 0.4 THEN 'ランニングシューズ'
        WHEN RANDOM() < 0.5 THEN 'ソファ'
        WHEN RANDOM() < 0.6 THEN 'スキンケア'
        WHEN RANDOM() < 0.7 THEN 'MacBook'
        WHEN RANDOM() < 0.8 THEN 'ダウンジャケット'
        WHEN RANDOM() < 0.9 THEN 'ヘッドフォン'
        ELSE 'プログラミング本'
    END,
    uuid_generate_v4(),
    CURRENT_TIMESTAMP - INTERVAL '30 days' + (RANDOM() * INTERVAL '30 days')
 FROM customers c
 CROSS JOIN LATERAL (
     SELECT id FROM products ORDER BY RANDOM() LIMIT 1
 ) p
 CROSS JOIN generate_series(1, 50) gs;

-- Generate wishlist items
INSERT INTO wishlist_items (customer_id, product_id)
SELECT DISTINCT
    c.id,
    p.id
FROM customers c
CROSS JOIN products p
WHERE RANDOM() < 0.15
AND NOT EXISTS (
    SELECT 1 FROM order_items oi
    JOIN orders o ON oi.order_id = o.id
    WHERE o.customer_id = c.id AND oi.product_id = p.id
);

-- Generate cart items
INSERT INTO cart_items (customer_id, product_id, quantity)
SELECT DISTINCT
    c.id,
    p.id,
    CASE WHEN RANDOM() < 0.8 THEN 1 ELSE 2 END
FROM customers c
CROSS JOIN products p
WHERE RANDOM() < 0.08
AND NOT EXISTS (
    SELECT 1 FROM order_items oi
    JOIN orders o ON oi.order_id = o.id
    WHERE o.customer_id = c.id AND oi.product_id = p.id
)
AND NOT EXISTS (
    SELECT 1 FROM wishlist_items wi
    WHERE wi.customer_id = c.id AND wi.product_id = p.id
);

-- Generate recommendation logs for testing
INSERT INTO recommendation_logs (customer_id, session_id, recommendation_type, context_type, recommended_products, clicked_products, purchased_products, algorithm_version, confidence_scores, created_at)
SELECT
    c.id,
    uuid_generate_v4(),
    CASE
        WHEN RANDOM() < 0.3 THEN 'collaborative'
        WHEN RANDOM() < 0.6 THEN 'content_based'
        WHEN RANDOM() < 0.8 THEN 'hybrid'
        ELSE 'similar'
    END,
    CASE
        WHEN RANDOM() < 0.4 THEN 'homepage'
        WHEN RANDOM() < 0.7 THEN 'product_page'
        WHEN RANDOM() < 0.9 THEN 'cart'
        ELSE 'checkout'
    END,
    ARRAY(SELECT id FROM products ORDER BY RANDOM() LIMIT 5),
    ARRAY(SELECT id FROM products ORDER BY RANDOM() LIMIT CASE WHEN RANDOM() < 0.7 THEN 1 ELSE 2 END),
    CASE WHEN RANDOM() < 0.3 THEN ARRAY(SELECT id FROM products ORDER BY RANDOM() LIMIT 1) ELSE ARRAY[]::UUID[] END,
    'v1.0',
         ARRAY[0.85, 0.78, 0.72, 0.69, 0.65],
     CURRENT_TIMESTAMP - INTERVAL '7 days' + (RANDOM() * INTERVAL '7 days')
 FROM customers c
 CROSS JOIN generate_series(1, 10) gs;
