-- Categories test data for EC recommendation system

-- Base categories (avoid duplicates)
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
('automotive', '自動車・バイク')
ON CONFLICT (name) DO NOTHING;

-- Extended sub-categories (avoid duplicates)
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
('maintenance', 10, 'メンテナンス用品')
ON CONFLICT (name) DO NOTHING;
