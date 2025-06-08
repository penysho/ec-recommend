-- Extended sample data for EC recommendation system testing

-- Clean up existing order-related data for regeneration
DELETE FROM recommendation_logs;
DELETE FROM customer_activities;
DELETE FROM cart_items;
DELETE FROM wishlist_items;
DELETE FROM product_reviews;
DELETE FROM order_items;
DELETE FROM orders;

-- Reset customer order statistics
UPDATE customers SET
    total_spent = 0,
    order_count = 0,
    is_premium = false;

-- Generate more realistic order data with multiple orders per customer
WITH customer_orders AS (
    SELECT
        c.id as customer_id,
        c.email,
        generate_series(1, CASE
            -- Premium customers (high spending potential)
            WHEN c.price_range_max > 150000 THEN 20 + floor(random() * 20)::integer
            -- High-medium customers
            WHEN c.price_range_max > 80000 THEN 15 + floor(random() * 15)::integer
            -- Medium customers
            WHEN c.price_range_max > 40000 THEN 8 + floor(random() * 12)::integer
            -- Low-medium customers
            WHEN c.price_range_max > 20000 THEN 5 + floor(random() * 10)::integer
            -- Budget customers
            ELSE 2 + floor(random() * 6)::integer
        END) as order_number
    FROM customers c
),
order_data AS (
    INSERT INTO orders (customer_id, order_number, status, subtotal, tax_amount, shipping_fee, total_amount, payment_method, ordered_at, delivered_at)
    SELECT
        co.customer_id,
        'ORD-' || LPAD((ROW_NUMBER() OVER())::text, 5, '0'),
        CASE WHEN RANDOM() < 0.9 THEN 'delivered' ELSE 'shipped' END,
        (CASE
            WHEN RANDOM() < 0.3 THEN RANDOM() * 10000 + 1000
            WHEN RANDOM() < 0.6 THEN RANDOM() * 30000 + 5000
            WHEN RANDOM() < 0.85 THEN RANDOM() * 80000 + 10000
            ELSE RANDOM() * 200000 + 50000
        END)::numeric(10,2),
        ((CASE
            WHEN RANDOM() < 0.3 THEN RANDOM() * 10000 + 1000
            WHEN RANDOM() < 0.6 THEN RANDOM() * 30000 + 5000
            WHEN RANDOM() < 0.85 THEN RANDOM() * 80000 + 10000
            ELSE RANDOM() * 200000 + 50000
        END) * 0.1)::numeric(10,2),
        CASE WHEN RANDOM() < 0.3 THEN 500 ELSE 0 END,
        ((CASE
            WHEN RANDOM() < 0.3 THEN RANDOM() * 10000 + 1000
            WHEN RANDOM() < 0.6 THEN RANDOM() * 30000 + 5000
            WHEN RANDOM() < 0.85 THEN RANDOM() * 80000 + 10000
            ELSE RANDOM() * 200000 + 50000
        END) * 1.1 + CASE WHEN RANDOM() < 0.3 THEN 500 ELSE 0 END)::numeric(10,2),
        CASE
            WHEN RANDOM() < 0.7 THEN 'credit_card'
            WHEN RANDOM() < 0.9 THEN 'bank_transfer'
            ELSE 'mobile_payment'
        END,
        CURRENT_TIMESTAMP - INTERVAL '730 days' + (RANDOM() * INTERVAL '730 days'),
        CURRENT_TIMESTAMP - INTERVAL '725 days' + (RANDOM() * INTERVAL '725 days')
    FROM customer_orders co
    ON CONFLICT (order_number) DO NOTHING
    RETURNING id, customer_id, order_number
)
-- Generate order items for each order with multiple items per order
INSERT INTO order_items (order_id, product_id, quantity, unit_price, total_price)
WITH order_item_generation AS (
    SELECT
        od.id as order_id,
        od.customer_id,
        gs.item_number,
        CASE
            WHEN RANDOM() < 0.4 THEN 1
            WHEN RANDOM() < 0.6 THEN 2
            WHEN RANDOM() < 0.8 THEN 3
            WHEN RANDOM() < 0.9 THEN 4
            WHEN RANDOM() < 0.95 THEN 5
            ELSE 6
        END as items_in_order
    FROM order_data od
    CROSS JOIN generate_series(1, 6) gs(item_number)
),
filtered_items AS (
    SELECT
        oig.order_id,
        oig.customer_id,
        oig.item_number,
        p.id as product_id,
        p.price,
        CASE
            WHEN RANDOM() < 0.6 THEN 1
            WHEN RANDOM() < 0.8 THEN 2
            WHEN RANDOM() < 0.95 THEN 3
            ELSE 4
        END as quantity
    FROM order_item_generation oig
    CROSS JOIN LATERAL (
        SELECT id, price
        FROM products
        ORDER BY RANDOM()
        LIMIT 1
    ) p
    WHERE oig.item_number <= oig.items_in_order
)
SELECT
    fi.order_id,
    fi.product_id,
    fi.quantity,
    fi.price,
    fi.price * fi.quantity
FROM filtered_items fi;

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
WHERE RANDOM() < 0.6 AND o.status = 'delivered'
ON CONFLICT (product_id, customer_id, order_id) DO NOTHING;

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
  CROSS JOIN generate_series(1, CASE
     WHEN c.price_range_max > 100000 THEN 80
     WHEN c.price_range_max > 50000 THEN 50
     WHEN c.price_range_max > 20000 THEN 30
     ELSE 15
 END) gs;

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
)
ON CONFLICT (customer_id, product_id) DO NOTHING;

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
)
ON CONFLICT (customer_id, product_id) DO NOTHING;

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
