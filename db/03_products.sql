-- Extended products (50+ products)
INSERT INTO products (name, description, category_id, price, original_price, brand, sku, stock_quantity, features, tags, rating_average, rating_count, popularity_score) VALUES
-- Electronics
('iPhone 15 Pro', '最新のiPhone Pro シリーズ', 11, 159800, 179800, 'Apple', 'IPH15P-128', 50, '{"storage": "128GB", "color": "Natural Titanium", "camera": "48MP"}', '{"smartphone", "ios", "premium"}', 4.5, 120, 95),
('iPhone 15', 'スタンダードモデル', 11, 124800, 139800, 'Apple', 'IPH15-128', 80, '{"storage": "128GB", "color": "Blue", "camera": "48MP"}', '{"smartphone", "ios", "standard"}', 4.3, 95, 88),
('Galaxy S24', 'Samsung最新フラグシップ', 11, 149800, 169800, 'Samsung', 'GXS24-256', 60, '{"storage": "256GB", "color": "Titanium Gray", "camera": "50MP"}', '{"smartphone", "android", "premium"}', 4.4, 87, 82),
('Pixel 8', 'Google純正Android', 11, 112800, 129800, 'Google', 'PXL8-128', 70, '{"storage": "128GB", "color": "Hazel", "camera": "50MP"}', '{"smartphone", "android", "google"}', 4.2, 76, 79),
('Xperia 1 V', 'ソニーフラグシップスマートフォン', 11, 194800, 214800, 'Sony', 'XPR1V-512', 35, '{"storage": "512GB", "color": "Platinum Silver", "camera": "48MP", "display": "4K"}', '{"smartphone", "android", "camera", "premium"}', 4.3, 68, 78),
('iPhone 14', '前世代iPhone', 11, 112800, 124800, 'Apple', 'IPH14-128', 95, '{"storage": "128GB", "color": "Midnight", "camera": "12MP"}', '{"smartphone", "ios", "standard"}', 4.2, 156, 85),
('AQUOS sense8', 'シャープミドルレンジ', 11, 59800, 69800, 'Sharp', 'AQS8-128', 120, '{"storage": "128GB", "color": "Light Copper", "battery": "5000mAh"}', '{"smartphone", "android", "budget"}', 4.1, 89, 72),
('Redmi Note 13', 'コストパフォーマンス重視', 11, 29800, 34800, 'Xiaomi', 'RMN13-128', 150, '{"storage": "128GB", "color": "Midnight Black", "camera": "108MP"}', '{"smartphone", "android", "budget"}', 4.0, 112, 69),

('MacBook Air M3', '軽量高性能ノートパソコン', 12, 164800, 184800, 'Apple', 'MBA-M3-13', 30, '{"cpu": "M3", "ram": "8GB", "storage": "256GB"}', '{"laptop", "mac", "portable"}', 4.7, 85, 88),
('MacBook Pro M3', 'プロ仕様ノートパソコン', 12, 248800, 268800, 'Apple', 'MBP-M3-14', 25, '{"cpu": "M3 Pro", "ram": "18GB", "storage": "512GB"}', '{"laptop", "mac", "professional"}', 4.8, 92, 91),
('ThinkPad X1', 'ビジネス向け高性能ノート', 12, 198800, 218800, 'Lenovo', 'TPX1-I7-16', 40, '{"cpu": "Intel i7", "ram": "16GB", "storage": "512GB"}', '{"laptop", "windows", "business"}', 4.5, 78, 85),
('Surface Laptop 5', 'Microsoft純正ノート', 12, 158800, 178800, 'Microsoft', 'SFL5-I5-8', 35, '{"cpu": "Intel i5", "ram": "8GB", "storage": "256GB"}', '{"laptop", "windows", "microsoft"}', 4.3, 69, 77),
('Dell XPS 13', 'プレミアムウルトラブック', 12, 168800, 188800, 'Dell', 'XPS13-I7-16', 25, '{"cpu": "Intel i7", "ram": "16GB", "storage": "512GB", "display": "4K"}', '{"laptop", "windows", "premium"}', 4.4, 93, 81),
('ASUS ZenBook', 'スタイリッシュノートPC', 12, 118800, 138800, 'ASUS', 'ZB14-R7-8', 45, '{"cpu": "AMD Ryzen 7", "ram": "8GB", "storage": "512GB"}', '{"laptop", "windows", "portable"}', 4.2, 76, 74),
('HP Spectre x360', '2-in-1コンバーチブル', 12, 178800, 198800, 'HP', 'SPX360-I7-16', 20, '{"cpu": "Intel i7", "ram": "16GB", "storage": "1TB", "touchscreen": true}', '{"laptop", "2in1", "touchscreen"}', 4.3, 65, 79),
('Gaming Laptop Legion', 'ゲーミングノートPC', 12, 229800, 259800, 'Lenovo', 'LGN-RTX4070-32', 15, '{"cpu": "Intel i7", "ram": "32GB", "gpu": "RTX 4070", "storage": "1TB"}', '{"laptop", "gaming", "high_performance"}', 4.6, 87, 88),

('Sony WH-1000XM5', 'ノイズキャンセリングヘッドフォン', 13, 49500, 54500, 'Sony', 'WH1000XM5-B', 75, '{"noise_cancelling": true, "battery_life": "30h", "wireless": true}', '{"headphones", "wireless", "noise-cancelling"}', 4.6, 203, 92),
('AirPods Pro 2', 'Apple純正ワイヤレスイヤホン', 13, 39800, 44800, 'Apple', 'APP2-USB-C', 90, '{"noise_cancelling": true, "battery_life": "6h", "wireless": true}', '{"earphones", "wireless", "apple"}', 4.5, 156, 89),
('Bose QC45', 'プレミアムノイズキャンセリング', 13, 42800, 47800, 'Bose', 'QC45-BLACK', 65, '{"noise_cancelling": true, "battery_life": "24h", "wireless": true}', '{"headphones", "wireless", "premium"}', 4.4, 134, 86),
('Sennheiser HD660S', 'オーディオファイル向けヘッドフォン', 13, 68800, 78800, 'Sennheiser', 'HD660S', 30, '{"impedance": "150ohm", "open_back": true, "wired": true}', '{"headphones", "audiophile", "wired"}', 4.7, 78, 85),
('Audio-Technica ATH-M50x', 'プロフェッショナルモニターヘッドフォン', 13, 19800, 24800, 'Audio-Technica', 'ATH-M50X', 80, '{"closed_back": true, "wired": true, "professional": true}', '{"headphones", "professional", "wired"}', 4.4, 145, 82),
('Beats Studio3', 'スタイリッシュワイヤレスヘッドフォン', 13, 34800, 39800, 'Beats', 'BS3-W-RED', 60, '{"noise_cancelling": true, "battery_life": "22h", "wireless": true}', '{"headphones", "wireless", "stylish"}', 4.2, 123, 78),
('JBL Tune 770NC', '手頃な価格のノイズキャンセリング', 13, 14800, 17800, 'JBL', 'T770NC-BLK', 100, '{"noise_cancelling": true, "battery_life": "44h", "wireless": true}', '{"headphones", "budget", "wireless"}', 4.1, 167, 74),
('Galaxy Buds2 Pro', 'サムスン純正ワイヤレスイヤホン', 13, 32800, 37800, 'Samsung', 'GB2P-BLK', 85, '{"noise_cancelling": true, "battery_life": "5h", "wireless": true}', '{"earphones", "wireless", "samsung"}', 4.3, 98, 79),

('iPad Air', 'バランス型タブレット', 14, 92800, 102800, 'Apple', 'IPA-M2-64', 45, '{"cpu": "M2", "storage": "64GB", "screen": "10.9inch"}', '{"tablet", "apple", "portable"}', 4.4, 98, 83),
('Surface Pro 9', 'ノートPC代替タブレット', 14, 148800, 168800, 'Microsoft', 'SPR9-I5-8', 30, '{"cpu": "Intel i5", "ram": "8GB", "storage": "256GB"}', '{"tablet", "windows", "2in1"}', 4.3, 74, 79),
('iPad Pro 12.9', 'プロ仕様大画面タブレット', 14, 172800, 192800, 'Apple', 'IPP-M2-128', 25, '{"cpu": "M2", "storage": "128GB", "screen": "12.9inch"}', '{"tablet", "apple", "professional"}', 4.6, 87, 86),
('Galaxy Tab S9', 'Android最上位タブレット', 14, 118800, 138800, 'Samsung', 'GTS9-128', 40, '{"storage": "128GB", "screen": "11inch", "s_pen": true}', '{"tablet", "android", "premium"}', 4.3, 76, 81),
('Fire HD 10', 'エンターテインメント特化タブレット', 14, 19980, 24980, 'Amazon', 'FHD10-32', 120, '{"storage": "32GB", "screen": "10.1inch", "alexa": true}', '{"tablet", "budget", "entertainment"}', 3.9, 234, 68),

('Apple Watch Series 9', '最新スマートウォッチ', 15, 59800, 69800, 'Apple', 'AWS9-GPS-41', 55, '{"gps": true, "health_monitoring": true, "battery_life": "18h"}', '{"smartwatch", "health", "apple"}', 4.5, 167, 87),
('Galaxy Watch 6', 'Samsung最新ウォッチ', 15, 47800, 52800, 'Samsung', 'GW6-BT-40', 60, '{"bluetooth": true, "health_monitoring": true, "battery_life": "40h"}', '{"smartwatch", "health", "samsung"}', 4.3, 89, 81),
('Fitbit Charge 5', 'フィットネストラッカー', 15, 24800, 29800, 'Fitbit', 'FC5-BLK-S', 80, '{"fitness_tracking": true, "gps": true, "battery_life": "7days"}', '{"fitness_tracker", "health", "budget"}', 4.2, 145, 76),
('Garmin Forerunner 955', 'ランニング専用ウォッチ', 15, 69800, 79800, 'Garmin', 'FR955-BLK', 35, '{"gps": true, "music": true, "battery_life": "15days"}', '{"running_watch", "gps", "professional"}', 4.6, 78, 84),
('Huawei Watch GT 4', 'スタイリッシュスマートウォッチ', 15, 34800, 39800, 'Huawei', 'WGT4-46MM', 50, '{"health_monitoring": true, "battery_life": "14days", "elegant": true}', '{"smartwatch", "stylish", "long_battery"}', 4.2, 93, 78),

('Canon EOS R6 Mark II', 'フルフレームミラーレス', 16, 348800, 398800, 'Canon', 'CR6M2-BODY', 15, '{"sensor": "24.2MP", "video": "4K", "image_stabilization": true}', '{"camera", "mirrorless", "professional"}', 4.7, 56, 89),
('Sony α7 IV', 'ハイブリッドミラーレス', 16, 298800, 348800, 'Sony', 'A7IV-BODY', 20, '{"sensor": "33MP", "video": "4K", "dual_card_slots": true}', '{"camera", "mirrorless", "hybrid"}', 4.6, 67, 87),
('Fujifilm X-T5', 'APS-C最高峰ミラーレス', 16, 248800, 278800, 'Fujifilm', 'XT5-BODY', 25, '{"sensor": "40.2MP", "film_simulation": true, "weather_sealed": true}', '{"camera", "mirrorless", "apsc"}', 4.5, 78, 85),
('Nikon Z8', 'プロ仕様ミラーレス', 16, 498800, 548800, 'Nikon', 'Z8-BODY', 8, '{"sensor": "45.7MP", "video": "8K", "dual_processors": true}', '{"camera", "mirrorless", "flagship"}', 4.8, 34, 91),
('DJI Mini 4 Pro', 'コンパクトドローン', 16, 118800, 138800, 'DJI', 'MINI4P-RC', 30, '{"weight": "249g", "video": "4K", "obstacle_avoidance": true}', '{"drone", "compact", "beginner"}', 4.4, 89, 82),

('PS5', 'PlayStation 5 ゲーム機', 17, 66980, 76980, 'Sony', 'PS5-STD', 40, '{"cpu": "AMD Zen 2", "gpu": "RDNA 2", "storage": "825GB"}', '{"gaming", "console", "next_gen"}', 4.6, 234, 95),
('Xbox Series X', 'Xbox最上位機種', 17, 59978, 69978, 'Microsoft', 'XSX-1TB', 30, '{"cpu": "AMD Zen 2", "gpu": "RDNA 2", "storage": "1TB"}', '{"gaming", "console", "4k"}', 4.5, 178, 92),
('Nintendo Switch OLED', '有機ELディスプレイ搭載', 17, 37980, 42980, 'Nintendo', 'NSW-OLED', 60, '{"screen": "7inch OLED", "portable": true, "battery_life": "9h"}', '{"gaming", "portable", "family"}', 4.7, 345, 96),
('Steam Deck', 'ポータブルPCゲーミング', 17, 79800, 89800, 'Valve', 'SD-512GB', 25, '{"cpu": "AMD APU", "storage": "512GB", "steam_os": true}', '{"gaming", "portable", "pc"}', 4.4, 123, 88),

-- Fashion - メンズファッション
('プレミアムダウンジャケット', '高品質ダウンジャケット', 18, 29800, 39800, 'UNIQLO', 'PDJ-M-BK-L', 100, '{"material": "90% down", "water_resistant": true, "size": "L"}', '{"outerwear", "winter", "mens"}', 4.3, 156, 78),
('カシミヤコート', 'エレガントなカシミヤコート', 18, 89800, 109800, 'ZARA', 'CSH-CT-BG-M', 25, '{"material": "100% cashmere", "color": "beige", "size": "M"}', '{"outerwear", "luxury", "mens"}', 4.6, 67, 85),
('デニムジャケット', 'カジュアルデニムジャケット', 18, 12800, 15800, 'LEVIS', 'DNM-JK-BL-L', 80, '{"material": "cotton", "color": "blue", "size": "L"}', '{"outerwear", "casual", "mens"}', 4.2, 134, 72),
('ウールスーツ', 'ビジネススーツ上下セット', 18, 59800, 79800, 'AOKi', 'WL-SUT-NV-L', 35, '{"material": "wool", "color": "navy", "size": "L", "2_piece": true}', '{"suit", "business", "formal"}', 4.4, 89, 81),
('ポロシャツ', 'コットンポロシャツ', 18, 5800, 7800, 'Ralph Lauren', 'POLO-CT-WH-M', 120, '{"material": "cotton", "color": "white", "size": "M"}', '{"shirt", "casual", "classic"}', 4.2, 167, 73),
('チノパンツ', 'スリムフィットチノ', 18, 8800, 11800, 'GAP', 'CHINO-SLM-BG-32', 90, '{"fit": "slim", "color": "beige", "size": "32"}', '{"pants", "casual", "versatile"}', 4.1, 145, 70),
('ニットセーター', 'メリノウールセーター', 18, 15800, 19800, 'MUJI', 'KNT-SW-GY-L', 60, '{"material": "merino wool", "color": "gray", "size": "L"}', '{"knitwear", "winter", "comfortable"}', 4.3, 98, 75),
('レザージャケット', '本革ライダースジャケット', 18, 89800, 119800, 'Schott', 'LTH-JK-BK-L', 15, '{"material": "genuine leather", "color": "black", "size": "L"}', '{"outerwear", "leather", "rock"}', 4.7, 45, 87),

-- Fashion - レディースファッション
('シルクブラウス', 'エレガントなシルクブラウス', 19, 15800, 19800, 'ZARA', 'SLK-BLS-WH-M', 45, '{"material": "100% silk", "color": "white", "size": "M"}', '{"blouse", "silk", "elegant", "womens"}', 4.4, 89, 71),
('ニットワンピース', '暖かいニットワンピース', 19, 24800, 29800, 'H&M', 'KNT-OP-GR-S', 60, '{"material": "wool blend", "color": "gray", "size": "S"}', '{"dress", "winter", "womens"}', 4.3, 76, 69),
('リネンシャツ', '涼しいリネンシャツ', 19, 8800, 11800, 'MUJI', 'LNN-SH-WH-M', 90, '{"material": "100% linen", "color": "white", "size": "M"}', '{"shirt", "summer", "womens"}', 4.1, 112, 65),
('フローラルワンピース', '花柄ロングワンピース', 19, 18800, 23800, 'ZARA', 'FLR-OP-M-L', 50, '{"pattern": "floral", "length": "midi", "size": "M"}', '{"dress", "floral", "feminine"}', 4.5, 123, 78),
('デニムスカート', 'Aラインデニムスカート', 19, 9800, 12800, 'GU', 'DNM-SK-BL-M', 80, '{"material": "denim", "color": "blue", "size": "M"}', '{"skirt", "denim", "casual"}', 4.2, 145, 72),
('カーディガン', 'ロングカーディガン', 19, 12800, 16800, 'UNIQLO', 'CDG-LG-BG-M', 70, '{"length": "long", "color": "beige", "size": "M"}', '{"cardigan", "layering", "versatile"}', 4.3, 134, 74),
('プリーツスカート', 'ミディ丈プリーツスカート', 19, 14800, 18800, 'H&M', 'PLT-SK-NV-S', 65, '{"style": "pleated", "color": "navy", "size": "S"}', '{"skirt", "pleated", "elegant"}', 4.4, 98, 76),
('ブレザー', 'テーラードブレザー', 19, 28800, 35800, 'ZARA', 'BLZ-TL-BK-M', 40, '{"style": "tailored", "color": "black", "size": "M"}', '{"blazer", "business", "sophisticated"}', 4.5, 67, 79),

-- 靴
('ランニングシューズ', '軽量高性能ランニングシューズ', 20, 12800, 16800, 'Nike', 'NKE-RUN-001-27', 60, '{"size": "27.0cm", "weight": "250g", "cushioning": "Air Zoom"}', '{"running", "shoes", "lightweight"}', 4.4, 187, 89),
('ビジネスシューズ', '本革ビジネスシューズ', 20, 28800, 35800, 'Cole Haan', 'CH-BIZ-BK-27', 40, '{"material": "genuine leather", "color": "black", "size": "27.0cm"}', '{"business", "leather", "formal"}', 4.5, 98, 82),
('スニーカー', 'カジュアルスニーカー', 20, 9800, 12800, 'Adidas', 'ADI-SNK-WH-26', 85, '{"color": "white", "size": "26.5cm", "style": "casual"}', '{"casual", "sneakers", "daily"}', 4.2, 156, 75),
('ハイヒール', 'エレガントパンプス', 20, 24800, 29800, 'Jimmy Choo', 'JC-PMP-BK-24', 30, '{"heel_height": "8cm", "color": "black", "size": "24.0cm"}', '{"heels", "formal", "elegant"}', 4.6, 78, 84),
('ブーツ', 'レザーアンクルブーツ', 20, 35800, 42800, 'Dr. Martens', 'DM-BT-BK-25', 45, '{"material": "leather", "color": "black", "size": "25.0cm"}', '{"boots", "leather", "sturdy"}', 4.4, 123, 81),
('サンダル', 'コンフォートサンダル', 20, 15800, 19800, 'Birkenstock', 'BK-SD-BN-26', 60, '{"material": "cork", "color": "brown", "size": "26.0cm"}', '{"sandals", "comfort", "summer"}', 4.3, 167, 76),
('ローファー', 'ペニーローファー', 20, 32800, 38800, 'Tods', 'TD-LF-BN-26', 25, '{"material": "leather", "color": "brown", "size": "26.0cm"}', '{"loafers", "casual_formal", "classic"}', 4.5, 56, 83),
('トレッキングシューズ', '防水ハイキングシューズ', 20, 22800, 27800, 'Merrell', 'MR-TRK-GY-27', 40, '{"waterproof": true, "color": "gray", "size": "27.0cm"}', '{"hiking", "outdoor", "waterproof"}', 4.6, 89, 85),

-- アクセサリー
('シルバーネックレス', 'スターリングシルバーチェーン', 21, 8800, 11800, 'Tiffany & Co.', 'TF-NCK-SLV-50', 50, '{"material": "sterling silver", "length": "50cm", "chain_type": "link"}', '{"jewelry", "silver", "classic"}', 4.4, 134, 78),
('ゴールドリング', '18金ゴールドリング', 21, 48800, 58800, 'Cartier', 'CT-RNG-GLD-17', 20, '{"material": "18k gold", "size": "17", "style": "classic"}', '{"jewelry", "gold", "luxury"}', 4.7, 45, 86),
('腕時計', '自動巻き腕時計', 21, 128800, 158800, 'Seiko', 'SK-WTH-AUTO-BK', 15, '{"movement": "automatic", "color": "black", "water_resistance": "100m"}', '{"watch", "automatic", "premium"}', 4.6, 67, 88),
('レザーベルト', '本革ビジネスベルト', 21, 12800, 15800, 'Coach', 'CCH-BLT-BK-90', 40, '{"material": "genuine leather", "color": "black", "size": "90cm"}', '{"belt", "leather", "business"}', 4.3, 123, 75),
('サングラス', 'UVカット偏光サングラス', 21, 18800, 23800, 'Ray-Ban', 'RB-SG-BK-UV', 60, '{"uv_protection": true, "polarized": true, "color": "black"}', '{"sunglasses", "uv_protection", "fashion"}', 4.5, 156, 82),
('スカーフ', 'シルクスカーフ', 21, 15800, 19800, 'Hermès', 'HM-SCF-SLK-90', 25, '{"material": "100% silk", "size": "90x90cm", "pattern": "geometric"}', '{"scarf", "silk", "luxury"}', 4.6, 78, 84),

-- バッグ
('トートバッグ', 'レザートートバッグ', 22, 38800, 45800, 'Michael Kors', 'MK-TOT-BK-L', 30, '{"material": "genuine leather", "color": "black", "size": "large"}', '{"bag", "tote", "leather"}', 4.4, 98, 80),
('リュックサック', 'ビジネスリュック', 22, 22800, 27800, 'Samsonite', 'SM-BP-BK-30L', 50, '{"capacity": "30L", "laptop_compartment": true, "color": "black"}', '{"backpack", "business", "laptop"}', 4.3, 145, 77),
('ショルダーバッグ', 'クロスボディバッグ', 22, 18800, 23800, 'Coach', 'CCH-CB-BN-S', 40, '{"style": "crossbody", "color": "brown", "size": "small"}', '{"bag", "crossbody", "compact"}', 4.2, 134, 74),
('クラッチバッグ', 'イブニングクラッチ', 22, 28800, 35800, 'Kate Spade', 'KS-CLT-BK-S', 25, '{"occasion": "evening", "color": "black", "size": "small"}', '{"bag", "clutch", "formal"}', 4.5, 67, 82),
('ボストンバッグ', 'トラベルボストンバッグ', 22, 45800, 55800, 'Louis Vuitton', 'LV-BST-BN-L', 15, '{"material": "leather", "color": "brown", "size": "large"}', '{"bag", "travel", "luxury"}', 4.7, 34, 87),

-- Home & Kitchen - より多くの製品を追加
('北欧スタイルソファ', 'モダンな北欧デザインソファ', 23, 89800, 109800, 'IKEA', 'NRD-SFA-GY-3P', 12, '{"seats": 3, "material": "fabric", "color": "gray"}', '{"sofa", "nordic", "furniture"}', 4.2, 67, 65),
('エルゴノミックチェア', '人間工学に基づいたオフィスチェア', 23, 68800, 78800, 'Herman Miller', 'HM-ERG-BK', 20, '{"adjustable": true, "lumbar_support": true, "color": "black"}', '{"chair", "ergonomic", "office"}', 4.7, 89, 88),
('ダイニングテーブル', '無垢材ダイニングテーブル', 23, 128800, 148800, 'Karimoku', 'KRM-DT-OAK-4', 8, '{"material": "oak", "seats": 4, "finish": "natural"}', '{"table", "dining", "wood"}', 4.6, 45, 79),
('本棚', '5段オープンシェルフ', 23, 24800, 29800, 'MUJI', 'MJ-SHF-OAK-5', 30, '{"shelves": 5, "material": "oak", "style": "open"}', '{"storage", "bookshelf", "minimalist"}', 4.3, 123, 72),
('ベッドフレーム', 'クイーンサイズベッドフレーム', 23, 89800, 109800, 'Nitori', 'NT-BF-QN-WH', 15, '{"size": "queen", "color": "white", "headboard": true}', '{"bed", "queen_size", "modern"}', 4.4, 78, 76),
('コーヒーテーブル', 'ガラストップコーヒーテーブル', 23, 45800, 55800, 'West Elm', 'WE-CT-GL-120', 20, '{"material": "glass_metal", "size": "120x60cm", "modern": true}', '{"table", "coffee", "glass"}', 4.2, 89, 74),
('ワードローブ', '3ドアワードローブ', 23, 78800, 98800, 'IKEA', 'IK-WR-3D-WH', 12, '{"doors": 3, "color": "white", "mirror": true}', '{"storage", "wardrobe", "bedroom"}', 4.1, 65, 71),
('テレビボード', '180cmテレビスタンド', 23, 35800, 42800, 'Lowya', 'LW-TV-180-WN', 25, '{"width": "180cm", "color": "walnut", "storage": true}', '{"tv_stand", "entertainment", "storage"}', 4.3, 98, 75),

('ステンレス製鍋セット', '高品質調理器具セット', 24, 25800, 32800, 'T-fal', 'SS-POT-SET-5', 80, '{"pieces": 5, "material": "stainless_steel", "induction_compatible": true}', '{"cookware", "kitchen", "stainless"}', 4.5, 142, 83),
('ブレンダー', '高性能ミキサーブレンダー', 24, 18800, 23800, 'Vitamix', 'VTX-BLD-001', 35, '{"power": "1200W", "capacity": "2L", "speed_settings": 10}', '{"blender", "kitchen", "healthy"}', 4.4, 76, 78),
('コーヒーメーカー', '全自動コーヒーメーカー', 24, 45800, 55800, 'Nespresso', 'NSP-CM-BK', 25, '{"type": "automatic", "milk_frother": true, "color": "black"}', '{"coffee", "automatic", "premium"}', 4.3, 134, 81),
('電子レンジ', 'オーブンレンジ', 24, 38800, 45800, 'Panasonic', 'PAN-OR-26L', 40, '{"capacity": "26L", "convection": true, "auto_cook": true}', '{"microwave", "oven", "versatile"}', 4.4, 167, 79),
('冷蔵庫', '3ドア冷蔵庫', 24, 158800, 178800, 'Sharp', 'SH-RF-350L', 10, '{"capacity": "350L", "doors": 3, "energy_efficient": true}', '{"refrigerator", "large", "efficient"}', 4.5, 89, 85),
('食器洗い機', '卓上型食器洗い機', 24, 68800, 78800, 'Panasonic', 'PAN-DW-6P', 20, '{"capacity": "6_place_settings", "compact": true, "eco_mode": true}', '{"dishwasher", "compact", "eco"}', 4.3, 76, 77),
('炊飯器', '圧力IH炊飯器', 24, 28800, 35800, 'Zojirushi', 'ZJ-RC-5.5', 50, '{"capacity": "5.5合", "pressure_ih": true, "multiple_menus": true}', '{"rice_cooker", "pressure", "premium"}', 4.6, 145, 82),
('エアフライヤー', 'ノンフライヤー', 24, 15800, 19800, 'Philips', 'PH-AF-3L', 60, '{"capacity": "3L", "oil_free": true, "digital_display": true}', '{"air_fryer", "healthy", "convenient"}', 4.2, 123, 76),

-- 寝具
('羽毛布団', 'プレミアム羽毛掛け布団', 25, 35800, 42800, '西川', 'NK-FD-DL-WH', 40, '{"size": "double", "down_percentage": "93%", "warmth_level": "medium"}', '{"bedding", "down", "premium"}', 4.5, 89, 83),
('枕', 'メモリーフォーム枕', 25, 8800, 11800, 'テンピュール', 'TMP-PIL-MF', 60, '{"material": "memory_foam", "size": "standard", "washable_cover": true}', '{"bedding", "pillow", "ergonomic"}', 4.4, 145, 78),
('ベッドシーツセット', 'オーガニックコットンシーツ', 25, 12800, 15800, 'MUJI', 'MJ-SH-OC-DL', 80, '{"material": "organic_cotton", "size": "double", "thread_count": "200"}', '{"bedding", "organic", "comfortable"}', 4.3, 134, 76),
('マットレス', 'ポケットコイルマットレス', 25, 89800, 109800, 'Simmons', 'SM-MT-PC-QN', 15, '{"size": "queen", "coil_type": "pocket", "firmness": "medium"}', '{"bedding", "mattress", "supportive"}', 4.6, 67, 85),

-- 照明
('ペンダントライト', 'モダンペンダント照明', 26, 18800, 23800, 'IKEA', 'IK-PL-MD-BK', 50, '{"style": "modern", "color": "black", "bulb_type": "LED"}', '{"lighting", "pendant", "modern"}', 4.2, 98, 74),
('フロアランプ', 'アーク型フロアランプ', 26, 28800, 35800, 'West Elm', 'WE-FL-ARC-BR', 30, '{"style": "arc", "color": "brass", "adjustable": true}', '{"lighting", "floor_lamp", "adjustable"}', 4.4, 76, 79),
('テーブルランプ', 'セラミックテーブルランプ', 26, 15800, 19800, 'Pottery Barn', 'PB-TL-CER-WH', 40, '{"material": "ceramic", "color": "white", "shade_included": true}', '{"lighting", "table_lamp", "ceramic"}', 4.3, 89, 75),
('LEDシーリングライト', 'リモコン付きシーリングライト', 26, 22800, 27800, 'Panasonic', 'PAN-CL-LED-12', 35, '{"type": "LED", "remote_control": true, "dimmer": true}', '{"lighting", "ceiling", "smart"}', 4.5, 123, 81),

-- 収納用品
('収納ボックス', 'スタッキング収納ボックス', 27, 3800, 4800, 'MUJI', 'MJ-SB-STK-M', 100, '{"stackable": true, "size": "medium", "clear": true}', '{"storage", "box", "stackable"}', 4.1, 167, 72),
('ハンガーラック', 'キャスター付きハンガーラック', 27, 8800, 11800, 'Nitori', 'NT-HR-CST-150', 60, '{"height": "150cm", "wheels": true, "adjustable": true}', '{"storage", "hanger", "mobile"}', 4.2, 134, 73),
('シューズラック', '5段シューズラック', 27, 12800, 15800, 'IKEA', 'IK-SR-5T-WH', 45, '{"tiers": 5, "color": "white", "capacity": "15_pairs"}', '{"storage", "shoes", "organized"}', 4.3, 98, 75),

-- Books - 技術書
('プログラミング入門', 'Go言語でWebアプリケーション開発', 31, 3200, null, '技術評論社', 'BK-GO-WEB-001', 200, '{"pages": 400, "level": "beginner", "language": "japanese"}', '{"programming", "go", "web", "beginner"}', 4.1, 78, 72),
('AI機械学習実践', 'PythonではじめるAI開発', 31, 4200, null, 'オライリー', 'BK-AI-PY-001', 150, '{"pages": 520, "level": "intermediate", "language": "japanese"}', '{"ai", "python", "machine_learning"}', 4.3, 89, 85),
('クラウド設計パターン', 'AWS実践アーキテクチャ', 31, 3800, null, '翔泳社', 'BK-AWS-ARC-001', 120, '{"pages": 360, "level": "advanced", "language": "japanese"}', '{"aws", "cloud", "architecture"}', 4.4, 67, 81),
('Web開発の教科書', 'React + Node.js実践', 31, 3600, null, 'SBクリエイティブ', 'BK-WEB-REACT-001', 180, '{"pages": 480, "level": "intermediate", "language": "japanese"}', '{"web", "react", "nodejs"}', 4.2, 134, 79),
('データベース設計', 'SQL実践ガイド', 31, 3400, null, '翔泳社', 'BK-DB-SQL-001', 160, '{"pages": 350, "level": "beginner", "language": "japanese"}', '{"database", "sql", "design"}', 4.3, 112, 77),

-- Books - ビジネス書
('投資の心理学', '行動経済学から学ぶ投資術', 30, 2800, null, '日経BP', 'BK-INV-PSY-001', 150, '{"pages": 320, "genre": "finance", "language": "japanese"}', '{"investment", "psychology", "finance"}', 4.3, 94, 68),
('マーケティング戦略', 'デジタル時代のブランド構築', 30, 3200, null, 'ダイヤモンド社', 'BK-MKT-DIG-001', 180, '{"pages": 280, "genre": "marketing", "language": "japanese"}', '{"marketing", "digital", "strategy"}', 4.2, 76, 74),
('リーダーシップ論', '次世代リーダーの条件', 30, 2600, null, '東洋経済新報社', 'BK-LDR-NXT-001', 200, '{"pages": 260, "genre": "leadership", "language": "japanese"}', '{"leadership", "management", "business"}', 4.1, 123, 69),
('起業の科学', 'スタートアップサイエンス', 30, 3400, null, '日経BP', 'BK-STP-SCI-001', 120, '{"pages": 380, "genre": "startup", "language": "japanese"}', '{"startup", "entrepreneurship", "science"}', 4.4, 89, 76),
('問題解決力', 'ロジカルシンキング実践', 30, 2400, null, 'PHP研究所', 'BK-LGC-THK-001', 250, '{"pages": 220, "genre": "thinking", "language": "japanese"}', '{"logic", "problem_solving", "thinking"}', 4.0, 167, 71),

-- Books - 小説・文学
('推理小説傑作選', '現代ミステリー名作集', 28, 1800, null, '新潮社', 'BK-MST-COL-001', 300, '{"pages": 450, "genre": "mystery", "language": "japanese"}', '{"mystery", "novel", "collection"}', 4.5, 234, 82),
('SF小説', '未来世界の物語', 28, 2200, null, 'ハヤカワ文庫', 'BK-SF-FUT-001', 250, '{"pages": 380, "genre": "sci_fi", "language": "japanese"}', '{"sci_fi", "novel", "future"}', 4.3, 156, 78),
('恋愛小説', '青春ラブストーリー', 28, 1600, null, '集英社', 'BK-ROM-YTH-001', 400, '{"pages": 320, "genre": "romance", "language": "japanese"}', '{"romance", "novel", "youth"}', 4.2, 189, 75),
('歴史小説', '戦国時代の英雄たち', 28, 2000, null, '文藝春秋', 'BK-HST-WAR-001', 200, '{"pages": 520, "genre": "historical", "language": "japanese"}', '{"historical", "novel", "samurai"}', 4.4, 98, 80),

-- Books - 漫画
('鬼滅の刃 全巻セット', '人気漫画コンプリートセット', 32, 12800, 15800, '集英社', 'MG-KMT-SET-23', 50, '{"volumes": 23, "genre": "action", "complete_set": true}', '{"manga", "action", "complete"}', 4.8, 456, 95),
('ワンピース 最新巻', '冒険漫画最新巻', 32, 550, null, '集英社', 'MG-OP-108', 200, '{"volume": 108, "genre": "adventure", "series": "ongoing"}', '{"manga", "adventure", "popular"}', 4.7, 234, 89),
('進撃の巨人 全巻', 'ダークファンタジー完結作品', 32, 15800, 18800, '講談社', 'MG-SNK-SET-34', 40, '{"volumes": 34, "genre": "dark_fantasy", "complete_set": true}', '{"manga", "fantasy", "complete"}', 4.6, 345, 88),
('呪術廻戦', '現代バトル漫画', 32, 528, null, '集英社', 'MG-JJK-25', 180, '{"volume": 25, "genre": "battle", "series": "ongoing"}', '{"manga", "battle", "supernatural"}', 4.5, 298, 86),
('スパイファミリー', 'コメディアクション漫画', 32, 594, null, '集英社', 'MG-SPY-12', 150, '{"volume": 12, "genre": "comedy", "series": "ongoing"}', '{"manga", "comedy", "family"}', 4.4, 267, 84),

-- Sports & Fitness - ランニング・フィットネス
('ヨガマット', 'プレミアムヨガマット', 35, 6800, 8800, 'Manduka', 'YGA-MAT-PRP-6MM', 40, '{"thickness": "6mm", "material": "PVC", "color": "purple"}', '{"yoga", "fitness", "mat"}', 4.6, 156, 76),
('ダンベルセット', '可変重量ダンベル', 35, 28800, 35800, 'Bowflex', 'BWF-DB-24KG', 25, '{"weight_range": "2-24kg", "adjustable": true, "space_saving": true}', '{"weights", "strength", "adjustable"}', 4.5, 98, 84),
('エクササイズバイク', '室内用フィットネスバイク', 35, 89800, 109800, 'Peloton', 'PLT-BK-001', 15, '{"resistance_levels": 32, "bluetooth": true, "display": "22inch"}', '{"cardio", "indoor", "connected"}', 4.4, 67, 79),
('トレッドミル', '家庭用ランニングマシン', 35, 128800, 148800, 'Johnson Health Tech', 'JHT-TM-001', 12, '{"max_speed": "16kmh", "incline": true, "heart_rate_monitor": true}', '{"cardio", "running", "home"}', 4.3, 45, 77),
('プルアップバー', 'ドアフレーム懸垂バー', 35, 3800, 4800, 'Perfect Pushup', 'PP-PUB-001', 80, '{"door_mount": true, "weight_limit": "136kg", "padding": true}', '{"strength", "pull_up", "compact"}', 4.2, 134, 73),
('ケトルベル', '16kgケトルベル', 35, 8800, 11800, 'CAP Barbell', 'CAP-KB-16KG', 60, '{"weight": "16kg", "material": "cast_iron", "grip": "wide_handle"}', '{"strength", "functional", "kettlebell"}', 4.4, 89, 78),
('フォームローラー', '筋膜リリースローラー', 35, 4800, 5800, 'TriggerPoint', 'TP-FR-001', 70, '{"density": "medium", "texture": "grooved", "length": "60cm"}', '{"recovery", "massage", "flexibility"}', 4.3, 123, 75),

-- Sports & Fitness - アウトドア
('テント', '4人用キャンプテント', 36, 38800, 45800, 'Coleman', 'CLM-TNT-4P', 30, '{"capacity": 4, "waterproof": true, "easy_setup": true}', '{"camping", "outdoor", "family"}', 4.3, 89, 77),
('登山リュック', '大容量ハイキングバックパック', 36, 24800, 29800, 'Osprey', 'OSP-BP-50L', 40, '{"capacity": "50L", "waterproof": true, "ergonomic": true}', '{"hiking", "backpack", "outdoor"}', 4.6, 76, 82),
('寝袋', '3シーズン用寝袋', 36, 15800, 19800, 'Marmot', 'MR-SB-3S', 45, '{"temperature_rating": "-5C", "weight": "1.2kg", "compressible": true}', '{"camping", "sleeping_bag", "3season"}', 4.4, 67, 79),
('キャンプチェア', '折りたたみアウトドアチェア', 36, 8800, 11800, 'Helinox', 'HX-CH-001', 60, '{"weight": "1kg", "load_capacity": "145kg", "packable": true}', '{"camping", "chair", "lightweight"}', 4.5, 98, 80),
('LEDランタン', '充電式LEDランタン', 36, 6800, 8800, 'Goal Zero', 'GZ-LNT-LED', 50, '{"brightness": "400lm", "battery_life": "48h", "usb_charging": true}', '{"camping", "lighting", "rechargeable"}', 4.3, 134, 76),
('クーラーボックス', '35L保冷ボックス', 36, 18800, 23800, 'YETI', 'YT-CB-35L', 25, '{"capacity": "35L", "ice_retention": "5days", "bear_resistant": true}', '{"camping", "cooler", "premium"}', 4.6, 56, 83),

-- Beauty - スキンケア
('保湿クリーム', 'オーガニック保湿クリーム', 38, 4800, 5800, 'SK-II', 'SKII-MOIST-50ML', 90, '{"volume": "50ml", "organic": true, "skin_type": "all"}', '{"skincare", "organic", "moisturizer"}', 4.7, 234, 91),
('美容液', 'ヒアルロン酸配合美容液', 38, 8800, 11800, 'Shiseido', 'SHI-SER-30ML', 60, '{"volume": "30ml", "ingredient": "hyaluronic_acid", "anti_aging": true}', '{"skincare", "serum", "anti_aging"}', 4.6, 167, 87),
('日焼け止め', 'SPF50+ PA++++', 38, 2800, null, 'Anessa', 'ANS-SS-60ML', 120, '{"volume": "60ml", "spf": 50, "waterproof": true}', '{"skincare", "sunscreen", "protection"}', 4.4, 198, 83),
('洗顔フォーム', 'センシティブスキン洗顔', 38, 2200, null, 'Cetaphil', 'CTP-CL-150ML', 100, '{"volume": "150ml", "gentle": true, "fragrance_free": true}', '{"skincare", "cleanser", "sensitive"}', 4.3, 145, 78),
('化粧水', '高保湿化粧水', 38, 3800, null, 'Hada Labo', 'HL-TN-170ML', 80, '{"volume": "170ml", "hyaluronic_acid": true, "alcohol_free": true}', '{"skincare", "toner", "hydrating"}', 4.5, 189, 82),
('アイクリーム', '目元専用クリーム', 38, 6800, 8800, 'Kiehls', 'KHL-EC-15ML', 70, '{"volume": "15ml", "caffeine": true, "anti_aging": true}', '{"skincare", "eye_cream", "anti_aging"}', 4.4, 123, 80),

-- Beauty - メイクアップ
('リップスティック', 'マットフィニッシュリップ', 39, 3200, null, 'MAC', 'MAC-LIP-RED-01', 120, '{"color": "classic_red", "finish": "matte", "long_lasting": true}', '{"makeup", "lipstick", "matte"}', 4.5, 167, 82),
('アイシャドウパレット', '12色アイシャドウ', 39, 6800, 8800, 'Urban Decay', 'UD-ES-12COL', 80, '{"colors": 12, "finish": "mixed", "pigmented": true}', '{"makeup", "eyeshadow", "palette"}', 4.4, 145, 79),
('ファンデーション', 'リキッドファンデーション', 39, 5800, null, 'NARS', 'NRS-FD-30ML', 100, '{"volume": "30ml", "coverage": "medium", "finish": "natural"}', '{"makeup", "foundation", "natural"}', 4.3, 123, 76),
('マスカラ', 'ボリューム&カールマスカラ', 39, 3800, null, 'Maybelline', 'MYB-MSC-BK', 150, '{"color": "black", "waterproof": true, "volume": true}', '{"makeup", "mascara", "waterproof"}', 4.2, 198, 74),
('アイライナー', 'リキッドアイライナー', 39, 2800, null, 'Kat Von D', 'KVD-EL-BK', 130, '{"color": "black", "precision": true, "long_wearing": true}', '{"makeup", "eyeliner", "precision"}', 4.4, 156, 77),
('チーク', 'パウダーチーク', 39, 3400, null, 'Tarte', 'TRT-CH-PCH', 90, '{"color": "peach", "buildable": true, "natural_finish": true}', '{"makeup", "blush", "natural"}', 4.3, 134, 75),

-- Beauty - ヘアケア
('シャンプー', 'スカルプケアシャンプー', 40, 3200, null, 'Kerastase', 'KRS-SHP-300ML', 80, '{"volume": "300ml", "scalp_care": true, "sulfate_free": true}', '{"haircare", "shampoo", "scalp"}', 4.4, 167, 79),
('コンディショナー', 'リペアコンディショナー', 40, 3400, null, 'Olaplex', 'OLP-CND-250ML', 70, '{"volume": "250ml", "bond_building": true, "damaged_hair": true}', '{"haircare", "conditioner", "repair"}', 4.5, 145, 81),
('ヘアオイル', 'アルガンヘアオイル', 40, 4800, 5800, 'Moroccan Oil', 'MO-OIL-50ML', 60, '{"volume": "50ml", "argan_oil": true, "frizz_control": true}', '{"haircare", "oil", "argan"}', 4.6, 123, 83),
('ヘアドライヤー', 'イオンドライヤー', 40, 28800, 35800, 'Dyson', 'DYS-HD-001', 30, '{"ion_technology": true, "heat_protection": true, "fast_drying": true}', '{"haircare", "dryer", "ion"}', 4.7, 89, 86),

-- Food - オーガニック食品
('オーガニック米', '無農薬コシヒカリ', 42, 3200, null, 'JAS認定農場', 'ORG-RC-5KG', 50, '{"weight": "5kg", "organic": true, "variety": "koshihikari"}', '{"food", "organic", "rice"}', 4.5, 134, 82),
('オーガニック野菜セット', '季節の有機野菜詰め合わせ', 42, 2800, null, '自然農法', 'ORG-VG-SET', 30, '{"variety": "seasonal", "organic": true, "pesticide_free": true}', '{"food", "organic", "vegetables"}', 4.6, 89, 85),
('オーガニックハチミツ', '純粋非加熱ハチミツ', 42, 2400, null, '山田養蜂場', 'ORG-HN-300G', 80, '{"weight": "300g", "raw": true, "organic": true}', '{"food", "organic", "honey"}', 4.4, 167, 80),
('オーガニックナッツ', 'ミックスナッツ', 42, 1800, null, 'ナチュラルハウス', 'ORG-NT-200G', 100, '{"weight": "200g", "mixed": true, "unsalted": true}', '{"food", "organic", "nuts"}', 4.3, 145, 76),

-- Food - 飲料
('プレミアムコーヒー豆', 'スペシャルティコーヒー', 43, 2800, null, 'ブルーボトル', 'BB-CF-200G', 60, '{"weight": "200g", "single_origin": true, "roast": "medium"}', '{"beverage", "coffee", "specialty"}', 4.6, 123, 84),
('緑茶', '静岡産煎茶', 43, 1200, null, '伊藤園', 'IT-GT-100G', 120, '{"weight": "100g", "origin": "shizuoka", "grade": "premium"}', '{"beverage", "tea", "green"}', 4.3, 189, 78),
('炭酸水', 'プレミアム炭酸水', 43, 180, null, 'San Pellegrino', 'SP-SW-500ML', 200, '{"volume": "500ml", "natural": true, "sparkling": true}', '{"beverage", "water", "sparkling"}', 4.2, 234, 72),
('フルーツジュース', '100%オレンジジュース', 43, 350, null, 'Tropicana', 'TP-OJ-1L', 150, '{"volume": "1L", "100_percent": true, "no_sugar_added": true}', '{"beverage", "juice", "orange"}', 4.1, 178, 70),

-- Food - スナック・お菓子
('プレミアムチョコレート', 'ベルギー産チョコレート', 44, 1800, null, 'Godiva', 'GDV-CH-200G', 80, '{"weight": "200g", "belgium": true, "dark_chocolate": true}', '{"snack", "chocolate", "premium"}', 4.7, 156, 88),
('おかき', '手焼きおかき詰め合わせ', 44, 1200, null, '播磨屋', 'HRM-OK-SET', 100, '{"variety": "assorted", "hand_baked": true, "traditional": true}', '{"snack", "rice_cracker", "traditional"}', 4.4, 134, 79),
('ポテトチップス', '厚切りポテトチップス', 44, 280, null, 'Calbee', 'CB-PC-120G', 200, '{"weight": "120g", "thick_cut": true, "sea_salt": true}', '{"snack", "chips", "potato"}', 4.0, 267, 71),
('ドライフルーツ', 'ミックスドライフルーツ', 44, 980, null, 'サンスウィート', 'SS-DF-150G', 120, '{"weight": "150g", "mixed": true, "no_sugar_added": true}', '{"snack", "dried_fruit", "healthy"}', 4.2, 189, 74),

-- Health - サプリメント・ビタミン
('マルチビタミン', '1日1粒マルチビタミン', 45, 2800, null, 'DHC', 'DHC-MV-90T', 150, '{"tablets": 90, "daily_dose": 1, "comprehensive": true}', '{"supplement", "vitamin", "daily"}', 4.3, 234, 77),
('オメガ3', 'フィッシュオイルサプリ', 45, 3200, null, 'Nature Made', 'NM-O3-120C', 100, '{"capsules": 120, "epa_dha": true, "fish_oil": true}', '{"supplement", "omega3", "heart_health"}', 4.4, 167, 80),
('プロテイン', 'ホエイプロテインパウダー', 45, 4800, null, 'Optimum Nutrition', 'ON-WP-2LB', 80, '{"weight": "2lb", "whey": true, "vanilla": true}', '{"supplement", "protein", "fitness"}', 4.5, 145, 82),
('コラーゲン', '美容コラーゲンサプリ', 45, 3800, null, '資生堂', 'SHI-COL-30D', 120, '{"days_supply": 30, "beauty": true, "marine_collagen": true}', '{"supplement", "collagen", "beauty"}', 4.2, 189, 75),
('ビタミンD', '高濃度ビタミンD3', 45, 1800, null, 'ナウフーズ', 'NOW-VD3-120C', 180, '{"capsules": 120, "high_potency": true, "bone_health": true}', '{"supplement", "vitamin_d", "bone"}', 4.4, 123, 78),

-- Toys - ボードゲーム
('人生ゲーム', 'クラシック人生ゲーム', 47, 3200, null, 'タカラトミー', 'TT-JG-001', 60, '{"players": "2-6", "age": "6+", "classic": true}', '{"board_game", "family", "classic"}', 4.3, 156, 78),
('カタン', '開拓者たちのカタン', 47, 4800, null, 'ジーピー', 'GP-CTN-001', 40, '{"players": "3-4", "age": "10+", "strategy": true}', '{"board_game", "strategy", "popular"}', 4.6, 123, 85),
('モノポリー', 'モノポリークラシック', 47, 2800, null, 'ハズブロ', 'HSB-MNP-001', 80, '{"players": "2-8", "age": "8+", "real_estate": true}', '{"board_game", "classic", "economics"}', 4.2, 189, 76),
('ドミニオン', 'デッキ構築ゲーム', 47, 5800, null, 'ホビージャパン', 'HJ-DMN-001', 35, '{"players": "2-4", "age": "13+", "deck_building": true}', '{"board_game", "strategy", "cards"}', 4.7, 89, 87),

-- Toys - ビデオゲーム
('ゼルダの伝説', 'ティアーズ オブ ザ キングダム', 48, 7678, null, '任天堂', 'NTD-ZLD-TOK', 100, '{"platform": "Nintendo Switch", "genre": "adventure", "rating": "E10+"}', '{"video_game", "adventure", "zelda"}', 4.8, 567, 96),
('ポケモン', 'スカーレット・バイオレット', 48, 6578, null, '任天堂', 'NTD-PKM-SV', 120, '{"platform": "Nintendo Switch", "genre": "rpg", "rating": "E"}', '{"video_game", "rpg", "pokemon"}', 4.5, 445, 91),
('マリオカート8', 'デラックス', 48, 6578, null, '任天堂', 'NTD-MK8-DX', 150, '{"platform": "Nintendo Switch", "genre": "racing", "rating": "E"}', '{"video_game", "racing", "mario"}', 4.7, 678, 94),
('ファイナルファンタジーXVI', 'FF16', 48, 8778, null, 'スクウェア・エニックス', 'SE-FF16-001', 80, '{"platform": "PS5", "genre": "rpg", "rating": "M"}', '{"video_game", "rpg", "final_fantasy"}', 4.4, 234, 87),

-- Toys - 知育玩具
('レゴクラシック', '基本ブロックセット', 49, 4800, null, 'LEGO', 'LG-CLS-001', 80, '{"pieces": 484, "age": "4-99", "creativity": true}', '{"educational_toy", "building", "creativity"}', 4.5, 234, 84),
('くもん式ひらがなカード', 'ひらがな学習カード', 49, 1200, null, 'くもん出版', 'KMN-HRG-001', 150, '{"cards": 46, "age": "2-6", "language": "japanese"}', '{"educational_toy", "language", "cards"}', 4.3, 167, 76),
('知育パズル', '日本地図パズル', 49, 2800, null, 'コクヨ', 'KKY-PZL-JPN', 100, '{"pieces": 47, "age": "5+", "geography": true}', '{"educational_toy", "puzzle", "geography"}', 4.4, 123, 79),
('プログラミングロボット', 'キッズ向けロボット', 49, 12800, 15800, 'Makeblock', 'MB-ROB-001', 40, '{"programmable": true, "age": "8+", "stem": true}', '{"educational_toy", "robot", "programming"}', 4.6, 89, 86),

-- Toys - アクションフィギュア
('ガンダムプラモデル', 'リアルグレード ガンダム', 50, 3200, null, 'バンダイ', 'BND-RG-001', 60, '{"scale": "1/144", "grade": "RG", "articulated": true}', '{"action_figure", "gundam", "model"}', 4.7, 145, 88),
('フィギュア', '鬼滅の刃 炭治郎フィギュア', 50, 8800, null, 'グッドスマイルカンパニー', 'GSC-KMT-001', 50, '{"character": "tanjirou", "scale": "1/8", "detailed": true}', '{"action_figure", "anime", "collectible"}', 4.6, 178, 85),
('ミニカー', 'トミカ消防車', 50, 550, null, 'タカラトミー', 'TT-TMK-001', 200, '{"vehicle": "fire_truck", "die_cast": true, "collectible": true}', '{"action_figure", "car", "miniature"}', 4.2, 234, 73),

-- Health - 医療機器
('血圧計', 'デジタル血圧計', 51, 8800, 11800, 'オムロン', 'OMR-BP-001', 60, '{"type": "digital", "memory": "60readings", "cuff_size": "standard"}', '{"medical_device", "blood_pressure", "health"}', 4.4, 134, 80),
('体温計', '非接触体温計', 51, 4800, 5800, 'テルモ', 'TRM-TMP-001', 80, '{"type": "non_contact", "response_time": "1sec", "fever_alarm": true}', '{"medical_device", "thermometer", "contactless"}', 4.3, 167, 78),
('体重計', 'スマート体重計', 51, 12800, 15800, 'タニタ', 'TNT-SC-001', 50, '{"smart": true, "body_composition": true, "app_sync": true}', '{"medical_device", "scale", "smart"}', 4.5, 98, 82),
('パルスオキシメーター', '血中酸素濃度計', 51, 3800, 4800, 'コニカミノルタ', 'KM-POX-001', 100, '{"measures": "SpO2", "portable": true, "digital_display": true}', '{"medical_device", "oxygen", "pulse"}', 4.2, 145, 76),

-- Health - ウェルネス
('アロマディフューザー', '超音波式アロマディフューザー', 52, 6800, 8800, '無印良品', 'MJ-ARD-001', 70, '{"type": "ultrasonic", "capacity": "500ml", "timer": true}', '{"wellness", "aromatherapy", "relaxation"}', 4.4, 123, 79),
('マッサージクッション', '首肩マッサージクッション', 52, 15800, 19800, 'ルルド', 'LRD-MSC-001', 40, '{"massage_type": "shiatsu", "heat": true, "portable": true}', '{"wellness", "massage", "relaxation"}', 4.3, 89, 77),
('空気清浄機', 'HEPAフィルター空気清浄機', 52, 28800, 35800, 'シャープ', 'SH-AP-001', 30, '{"filter": "HEPA", "coverage": "23畳", "plasma_cluster": true}', '{"wellness", "air_purifier", "health"}', 4.6, 67, 84),
('加湿器', '超音波式加湿器', 52, 8800, 11800, 'パナソニック', 'PAN-HMD-001', 60, '{"type": "ultrasonic", "capacity": "4L", "auto_shutoff": true}', '{"wellness", "humidifier", "comfort"}', 4.2, 134, 75),

-- Automotive - カー用品
('ドライブレコーダー', '前後2カメラドライブレコーダー', 53, 18800, 23800, 'コムテック', 'CMT-DR-001', 40, '{"cameras": 2, "resolution": "Full HD", "night_vision": true}', '{"car_accessory", "dashcam", "safety"}', 4.5, 98, 83),
('カーナビ', 'ポータブルナビゲーション', 53, 28800, 35800, 'ガーミン', 'GMN-NAV-001', 25, '{"screen": "7inch", "maps": "japan", "traffic": true}', '{"car_accessory", "navigation", "gps"}', 4.3, 76, 79),
('タイヤ', 'エコタイヤ', 53, 15800, 19800, 'ブリヂストン', 'BS-TR-195', 60, '{"size": "195/65R15", "eco": true, "fuel_efficient": true}', '{"car_accessory", "tire", "eco"}', 4.4, 89, 81),
('カーシャンプー', 'プレミアムカーシャンプー', 53, 2800, null, 'シュアラスター', 'SL-CS-001', 120, '{"volume": "1L", "wax": true, "water_repellent": true}', '{"car_accessory", "shampoo", "cleaning"}', 4.2, 145, 74),
('シートカバー', 'レザー調シートカバー', 53, 12800, 15800, 'クラッツィオ', 'CLZ-SC-001', 50, '{"material": "leather_like", "universal": true, "washable": true}', '{"car_accessory", "seat_cover", "interior"}', 4.1, 123, 72),

-- Automotive - メンテナンス用品
('エンジンオイル', '高性能エンジンオイル', 54, 4800, 5800, 'カストロール', 'CTR-OIL-4L', 80, '{"volume": "4L", "viscosity": "5W-30", "synthetic": true}', '{"maintenance", "oil", "engine"}', 4.4, 167, 80),
('タイヤワックス', 'タイヤ艶出しワックス', 54, 1800, null, 'ソフト99', 'S99-TW-001', 150, '{"volume": "300ml", "shine": true, "protection": true}', '{"maintenance", "wax", "tire"}', 4.1, 189, 73),
('バッテリー', 'カーバッテリー', 54, 8800, 11800, 'パナソニック', 'PAN-BAT-001', 40, '{"capacity": "55B24L", "maintenance_free": true, "long_life": true}', '{"maintenance", "battery", "electrical"}', 4.3, 98, 77),
('ブレーキフルード', 'DOT4ブレーキフルード', 54, 1200, null, 'ワコーズ', 'WKS-BF-001', 100, '{"type": "DOT4", "volume": "500ml", "high_boiling": true}', '{"maintenance", "brake_fluid", "safety"}', 4.2, 134, 75),

-- Magazine - より多くの雑誌
('ファッション雑誌', 'VOGUE JAPAN', 34, 980, null, 'コンデナスト・ジャパン', 'VG-JPN-2024', 200, '{"category": "fashion", "frequency": "monthly", "target": "women"}', '{"magazine", "fashion", "women"}', 4.2, 156, 76),
('車雑誌', 'CAR TOP', 34, 680, null, '交通タイムス社', 'CT-MAG-2024', 150, '{"category": "automotive", "frequency": "monthly", "target": "car_enthusiasts"}', '{"magazine", "automotive", "cars"}', 4.1, 123, 74),
('料理雑誌', 'オレンジページ', 34, 580, null, 'オレンジページ', 'OP-MAG-2024', 180, '{"category": "cooking", "frequency": "bi_weekly", "target": "homemaker"}', '{"magazine", "cooking", "lifestyle"}', 4.3, 189, 78),
('ゲーム雑誌', 'ファミ通', 34, 680, null, 'KADOKAWA', 'FMT-2024', 120, '{"category": "gaming", "frequency": "weekly", "target": "gamers"}', '{"magazine", "gaming", "entertainment"}', 4.0, 234, 72),
('IT雑誌', '日経コンピュータ', 34, 1800, null, '日経BP', 'NC-MAG-2024', 100, '{"category": "technology", "frequency": "bi_weekly", "target": "it_professionals"}', '{"magazine", "technology", "business"}', 4.4, 89, 81),

-- さらに多くの電子機器を追加
('Bluetoothスピーカー', 'ポータブルBluetoothスピーカー', 1, 8800, 11800, 'JBL', 'JBL-BT-001', 80, '{"wireless": true, "battery_life": "12h", "waterproof": "IPX7"}', '{"speaker", "bluetooth", "portable"}', 4.3, 167, 81),
('ワイヤレス充電器', 'Qi対応ワイヤレス充電パッド', 1, 3800, 4800, 'Anker', 'ANK-WC-001', 120, '{"wireless": true, "fast_charging": true, "qi_compatible": true}', '{"charger", "wireless", "convenient"}', 4.2, 189, 78),
('モバイルバッテリー', '20000mAhモバイルバッテリー', 1, 5800, 7800, 'Anker', 'ANK-MB-20K', 100, '{"capacity": "20000mAh", "fast_charging": true, "multiple_ports": true}', '{"battery", "portable", "power"}', 4.4, 145, 83),
('USB-Cハブ', '7-in-1 USB-Cハブ', 1, 6800, 8800, 'HyperDrive', 'HD-HUB-7IN1', 60, '{"ports": 7, "usb_c": true, "hdmi": true}', '{"hub", "usb_c", "connectivity"}', 4.1, 123, 76),
('Webカメラ', '4K Webカメラ', 1, 12800, 15800, 'Logitech', 'LGT-CAM-4K', 50, '{"resolution": "4K", "auto_focus": true, "built_in_mic": true}', '{"camera", "webcam", "streaming"}', 4.5, 98, 85),

-- より多くのキッチン用品
('電気ケトル', 'ステンレス電気ケトル', 24, 8800, 11800, 'ティファール', 'TF-KTL-1.2L', 70, '{"capacity": "1.2L", "auto_shutoff": true, "temperature_control": true}', '{"kettle", "electric", "convenient"}', 4.3, 134, 79),
('トースター', '2枚焼きポップアップトースター', 24, 6800, 8800, 'バルミューダ', 'BLM-TST-001', 50, '{"slices": 2, "steam": true, "precise_control": true}', '{"toaster", "bread", "kitchen"}', 4.6, 89, 86),
('ハンドブレンダー', 'マルチハンドブレンダー', 24, 12800, 15800, 'ブラウン', 'BRN-HB-001', 60, '{"attachments": 3, "variable_speed": true, "dishwasher_safe": true}', '{"blender", "hand", "versatile"}', 4.4, 123, 82),
('フードプロセッサー', '多機能フードプロセッサー', 24, 18800, 23800, 'クイジナート', 'CZN-FP-001', 40, '{"capacity": "3L", "attachments": 5, "powerful_motor": true}', '{"food_processor", "cooking", "efficient"}', 4.5, 78, 84),
('真空パック器', '家庭用真空パック機', 24, 15800, 19800, 'FoodSaver', 'FS-VS-001', 35, '{"vacuum": true, "sealing": true, "food_preservation": true}', '{"vacuum_sealer", "preservation", "storage"}', 4.2, 67, 77),

-- より多くのファッションアイテム
('カーディガン', 'カシミヤカーディガン', 2, 35800, 42800, 'UNIQLO', 'UNI-CDG-CSH', 40, '{"material": "cashmere", "color": "gray", "size": "M"}', '{"cardigan", "luxury", "soft"}', 4.5, 89, 83),
('マフラー', 'ウールマフラー', 2, 8800, 11800, 'ザラ', 'ZR-SCF-WL', 60, '{"material": "wool", "color": "navy", "warm": true}', '{"scarf", "wool", "winter"}', 4.2, 134, 76),
('手袋', 'レザーグローブ', 2, 12800, 15800, 'Coach', 'CCH-GLV-LTH', 50, '{"material": "leather", "lined": true, "touchscreen": true}', '{"gloves", "leather", "smart"}', 4.4, 78, 81),
('帽子', 'ニットビーニー', 2, 3800, 4800, 'パタゴニア', 'PTG-BN-001', 80, '{"material": "wool", "warm": true, "sustainable": true}', '{"hat", "beanie", "outdoor"}', 4.1, 156, 74),
('靴下', 'メリノウール靴下', 2, 2800, null, 'Smartwool', 'SW-SCK-MW', 120, '{"material": "merino_wool", "odor_resistant": true, "moisture_wicking": true}', '{"socks", "wool", "performance"}', 4.3, 189, 78)
ON CONFLICT (sku) DO NOTHING;
