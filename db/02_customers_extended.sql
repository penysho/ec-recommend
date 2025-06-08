-- Extended customers data generation for EC recommendation system

-- Generate 1000 customers with diverse characteristics
DO $$
DECLARE
    i INTEGER;
    first_names_male TEXT[] := ARRAY['太郎', '次郎', '三郎', '健太', '大輔', '翔太', '和也', '信一', '博', '竜也', '健一', '雅之', '直樹', '智也', '拓也', '昌平', '裕介', '雄大', '慎一', '康弘', '誠', '洋', '修', '勇', '実', '隆', '正', '茂', '豊', '敏', '武', '清', '進', '治', '浩', '哲', '純', '宏', '剛', '光'];
    first_names_female TEXT[] := ARRAY['花子', '美咲', '愛子', '麻衣', '由美', '千尋', '理恵', '涼子', '美穂', '恵美', '真理', '香織', '美奈', '智子', '礼子', '優子', '久美', '典子', '美幸', '知恵', '恵', '裕子', '直子', '幸子', '悦子', '京子', '洋子', '綾子', '和子', '節子', '文子', '良子', '明子', '清子', '春子', '美代子', '信子', '光子', '弘子', '雅子'];
    last_names TEXT[] := ARRAY['田中', '佐藤', '山田', '渡辺', '鈴木', '高橋', '伊藤', '小林', '中村', '加藤', '林', '森', '吉田', '松本', '井上', '木村', '山崎', '佐々木', '石井', '藤田', '近藤', '後藤', '清水', '前田', '山本', '新井', '池田', '橋本', '福田', '西村', '岡田', '長谷川', '村上', '近江', '斎藤', '菊地', '安田', '原田', '青木', '武田', '上田', '杉山', '千葉', '村田', '河野', '酒井', '今井', '山口', '石川', '工藤'];
    email_domains TEXT[] := ARRAY['example.com', 'test.jp', 'demo.co.jp', 'sample.net', 'mail.com'];
    phone_prefixes TEXT[] := ARRAY['090', '080', '070'];
    genders TEXT[] := ARRAY['male', 'female'];
    prefectures TEXT[] := ARRAY['Tokyo', 'Osaka', 'Kanagawa', 'Aichi', 'Saitama', 'Chiba', 'Hyogo', 'Hokkaido', 'Fukuoka', 'Shizuoka', 'Hiroshima', 'Kyoto', 'Ibaraki', 'Niigata', 'Miyagi'];
    cities TEXT[] := ARRAY['Central', 'North', 'South', 'East', 'West', 'Shibuya', 'Shinjuku', 'Ikebukuro', 'Ginza', 'Harajuku', 'Namba', 'Umeda', 'Tennoji', 'Yokohama', 'Kawasaki', 'Chiba', 'Funabashi'];
    brands TEXT[] := ARRAY['Apple', 'Sony', 'Nike', 'ZARA', 'UNIQLO', 'MUJI', 'H&M', 'IKEA', 'Shiseido', 'MAC', 'BMW', 'Nintendo', 'SK-II', 'GAP', 'Coach', 'Adidas', 'Panasonic', 'Canon', 'T-fal', 'Coleman'];
    lifestyle_tags TEXT[] := ARRAY['tech_enthusiast', 'fashion_lover', 'fitness', 'gamer', 'bookworm', 'outdoor', 'beauty', 'minimalist', 'premium', 'urban', 'professional', 'student', 'family', 'trendy', 'health_conscious', 'eco_friendly', 'luxury', 'creative', 'social', 'practical'];

    selected_first_name TEXT;
    selected_last_name TEXT;
    selected_gender TEXT;
    email_local TEXT;
    generated_email TEXT;
    phone_number TEXT;
    birth_date DATE;
    selected_prefecture TEXT;
    selected_city TEXT;
    price_min INTEGER;
    price_max INTEGER;
    selected_brands TEXT[];
    selected_lifestyle TEXT[];
    category_array TEXT;

BEGIN
    FOR i IN 1..1000 LOOP
        -- Select random gender
        selected_gender := genders[1 + floor(random() * array_length(genders, 1))];

        -- Select appropriate first name based on gender
        IF selected_gender = 'male' THEN
            selected_first_name := first_names_male[1 + floor(random() * array_length(first_names_male, 1))];
        ELSE
            selected_first_name := first_names_female[1 + floor(random() * array_length(first_names_female, 1))];
        END IF;

        -- Select random last name
        selected_last_name := last_names[1 + floor(random() * array_length(last_names, 1))];

        -- Generate unique email
        email_local := lower(selected_last_name) || i::text;
        generated_email := email_local || '@' || email_domains[1 + floor(random() * array_length(email_domains, 1))];

        -- Generate phone number
        phone_number := phone_prefixes[1 + floor(random() * array_length(phone_prefixes, 1))] ||
                       '-' || lpad((1000 + floor(random() * 9000))::text, 4, '0') ||
                       '-' || lpad((1000 + floor(random() * 9000))::text, 4, '0');

        -- Generate birth date (ages 18-65)
        birth_date := CURRENT_DATE - INTERVAL '18 years' - (random() * INTERVAL '47 years');

        -- Select location
        selected_prefecture := prefectures[1 + floor(random() * array_length(prefectures, 1))];
        selected_city := cities[1 + floor(random() * array_length(cities, 1))];

        -- Generate price range based on age and lifestyle
        CASE
            WHEN extract(year from age(birth_date)) < 25 THEN
                price_min := 500 + floor(random() * 1500);
                price_max := 15000 + floor(random() * 35000);
            WHEN extract(year from age(birth_date)) < 35 THEN
                price_min := 1000 + floor(random() * 3000);
                price_max := 25000 + floor(random() * 75000);
            WHEN extract(year from age(birth_date)) < 50 THEN
                price_min := 2000 + floor(random() * 8000);
                price_max := 50000 + floor(random() * 150000);
            ELSE
                price_min := 5000 + floor(random() * 15000);
                price_max := 80000 + floor(random() * 220000);
        END CASE;

        -- Generate preferred categories (2-4 categories)
        category_array := '{' ||
            (1 + floor(random() * 10))::text || ',' ||
            (1 + floor(random() * 10))::text || ',' ||
            (1 + floor(random() * 10))::text ||
            CASE WHEN random() < 0.5 THEN ',' || (1 + floor(random() * 10))::text ELSE '' END ||
            '}';

                -- Generate preferred brands array (2-4 brands)
        selected_brands := ARRAY[
            brands[1 + floor(random() * array_length(brands, 1))],
            brands[1 + floor(random() * array_length(brands, 1))]
        ]::text[];

        IF random() < 0.7 THEN
            selected_brands := selected_brands || brands[1 + floor(random() * array_length(brands, 1))];
        END IF;

        IF random() < 0.4 THEN
            selected_brands := selected_brands || brands[1 + floor(random() * array_length(brands, 1))];
        END IF;

        -- Generate lifestyle tags array (2-4 tags)
        selected_lifestyle := ARRAY[
            lifestyle_tags[1 + floor(random() * array_length(lifestyle_tags, 1))],
            lifestyle_tags[1 + floor(random() * array_length(lifestyle_tags, 1))]
        ]::text[];

        IF random() < 0.8 THEN
            selected_lifestyle := selected_lifestyle || lifestyle_tags[1 + floor(random() * array_length(lifestyle_tags, 1))];
        END IF;

        IF random() < 0.5 THEN
            selected_lifestyle := selected_lifestyle || lifestyle_tags[1 + floor(random() * array_length(lifestyle_tags, 1))];
        END IF;

        -- Insert customer
        INSERT INTO customers (
            email,
            first_name,
            last_name,
            phone,
            date_of_birth,
            gender,
            preferred_categories,
            price_range_min,
            price_range_max,
            preferred_brands,
            location,
            lifestyle_tags
        ) VALUES (
            generated_email,
            selected_first_name,
            selected_last_name,
            phone_number,
            birth_date,
            selected_gender,
            category_array::integer[],
            price_min,
            price_max,
            selected_brands,
            ('{"prefecture": "' || selected_prefecture || '", "city": "' || selected_city || '"}')::jsonb,
            selected_lifestyle
        ) ON CONFLICT (email) DO NOTHING;

    END LOOP;
END $$;

-- Display summary of customer data
SELECT
    'Total customers' as metric,
    count(*)::text as value
FROM customers
UNION ALL
SELECT
    'Gender distribution - Male',
    count(*)::text
FROM customers WHERE gender = 'male'
UNION ALL
SELECT
    'Gender distribution - Female',
    count(*)::text
FROM customers WHERE gender = 'female'
UNION ALL
SELECT
    'Age groups - 18-25',
    count(*)::text
FROM customers WHERE extract(year from age(date_of_birth)) BETWEEN 18 AND 25
UNION ALL
SELECT
    'Age groups - 26-35',
    count(*)::text
FROM customers WHERE extract(year from age(date_of_birth)) BETWEEN 26 AND 35
UNION ALL
SELECT
    'Age groups - 36-50',
    count(*)::text
FROM customers WHERE extract(year from age(date_of_birth)) BETWEEN 36 AND 50
UNION ALL
SELECT
    'Age groups - 51+',
    count(*)::text
FROM customers WHERE extract(year from age(date_of_birth)) > 50;
