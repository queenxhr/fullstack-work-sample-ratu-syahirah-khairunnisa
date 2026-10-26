-- ==========================================================
-- Database Schema: db_worksample_ratu
-- PostgreSQL Schema for Senior Fullstack Work Sample
-- Candidate: Ratu Syahirah Khairunnisa
-- ==========================================================

-- LANGKAH 1: Buat database dulu (jalankan di DBeaver sebagai superuser)
-- CREATE DATABASE db_worksample_ratu;

-- LANGKAH 2: Hubungkan ke database, lalu jalankan semua query di bawah ini
-- \c db_worksample_ratu;

-- ==========================================================
-- Table 1: products (Supports 5,000+ Items Virtualized Table)
-- ==========================================================
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    price NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    stock INT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'IN_STOCK',
    version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_products_stock ON products(stock);

-- ==========================================================
-- Table 2: orders (Supports Concurrency & Race Condition Tests)
-- ==========================================================
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(64) UNIQUE NOT NULL,
    product_id BIGINT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price NUMERIC(12, 2) NOT NULL,
    total_amount NUMERIC(12, 2) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'CONFIRMED',
    handling_method VARCHAR(32) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_orders_products
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_orders_product_id ON orders(product_id);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);

-- ==========================================================
-- Table 3: webhook_events (Supports Idempotency & Resiliency)
-- ==========================================================
CREATE TABLE IF NOT EXISTS webhook_events (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(100) UNIQUE NOT NULL,
    idempotency_key VARCHAR(100) UNIQUE NOT NULL,
    payment_id VARCHAR(100) NOT NULL,
    order_id VARCHAR(100) NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    raw_payload TEXT,
    processed_count INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_webhook_payment_id ON webhook_events(payment_id);
CREATE INDEX IF NOT EXISTS idx_webhook_order_id ON webhook_events(order_id);

-- ==========================================================
-- Seed Data 1: Product Khusus Flash Sale (Produk #1)
-- ==========================================================
INSERT INTO products (id, sku, name, category, price, stock, status, version) VALUES
(1, 'SKU-00001', '🔥 [FLASH SALE EXCLUSIVE] Sony PlayStation 5 Pro Limited (1 Remaining Unit)', 'Gaming & PC', 8999000.00, 1, 'LOW_STOCK', 1)
ON CONFLICT (id) DO UPDATE SET 
    stock = EXCLUDED.stock,
    status = EXCLUDED.status;

-- FIX: Setelah INSERT eksplisit id=1, BIGSERIAL sequence TIDAK otomatis advance.
-- Akibatnya 9 insert berikut (tanpa explicit id) semua dapat nextval=1 → PK conflict
-- → seluruh batch gagal diam-diam → total hanya 4991 bukan 5000.
-- Solusi: reset sequence ke nilai MAX(id) yang ada sekarang.
SELECT setval(pg_get_serial_sequence('products', 'id'), COALESCE((SELECT MAX(id) FROM products), 1));

-- ==========================================================
-- Seed Data 2: Sample Dummy Initial Products
-- ==========================================================
INSERT INTO products (sku, name, category, price, stock, status, version) VALUES
('SKU-00002', 'Apple MacBook Pro M3 Max 16-inch 36GB', 'Laptops', 38999000.00, 24, 'IN_STOCK', 1),
('SKU-00003', 'Sony WH-1000XM5 Wireless Noise Cancelling', 'Audio & Wearables', 4999000.00, 50, 'IN_STOCK', 1),
('SKU-00004', 'Logitech MX Master 3S Performance Wireless Mouse', 'Accessories', 1650000.00, 120, 'IN_STOCK', 1),
('SKU-00005', 'Samsung Galaxy S24 Ultra 512GB Titanium Gray', 'Smartphones', 21999000.00, 18, 'IN_STOCK', 1),
('SKU-00006', 'Asus ROG Swift OLED 27-inch 240Hz Gaming Monitor', 'Gaming & PC', 14500000.00, 8, 'IN_STOCK', 1),
('SKU-00007', 'Dell XPS 15 9530 Intel Core i9 32GB RAM', 'Laptops', 32500000.00, 12, 'IN_STOCK', 1),
('SKU-00008', 'Bose QuietComfort Ultra Headphones', 'Audio & Wearables', 5500000.00, 35, 'IN_STOCK', 1),
('SKU-00009', 'Razer BlackWidow V4 Pro Mechanical Keyboard', 'Gaming & PC', 3899000.00, 4, 'LOW_STOCK', 1),
('SKU-00010', 'LG DualUp 28-inch Ergonomic Productivity Monitor', 'Accessories', 9800000.00, 0, 'OUT_OF_STOCK', 1)
ON CONFLICT (sku) DO NOTHING;

-- ==========================================================
-- Script Otomatis: Generate 5.000+ Produk Lengkap di PostgreSQL
-- (Jalankan blok PL/pgSQL di bawah ini di DBeaver untuk mengisi 5.000 produk seketika)
-- ==========================================================
DO $$
DECLARE
    i INT;
    v_categories TEXT[] := ARRAY['Smartphones', 'Laptops', 'Audio & Wearables', 'Gaming & PC', 'Accessories', 'Home Appliances'];
    v_brands TEXT[] := ARRAY['Apple', 'Sony', 'Logitech', 'Samsung', 'Asus', 'Dell', 'Bose', 'Xiaomi', 'Razer', 'LG'];
    v_types TEXT[] := ARRAY['Pro Max', 'Wireless Ultra', 'Gaming Edition', 'Studio Monitor', 'Ergo Mechanical', 'OLED Smart', 'Thunderbolt Dock', 'Noise Cancelling'];
    v_cat TEXT;
    v_name TEXT;
    v_price NUMERIC(12, 2);
    v_stock INT;
    v_status TEXT;
BEGIN
    FOR i IN 11..5010 LOOP  -- range lebih lebar agar total tetap 5000 meski ada conflict
        v_cat := v_categories[1 + floor(random() * array_length(v_categories, 1))::int];
        v_name := v_brands[1 + floor(random() * array_length(v_brands, 1))::int] || ' ' || 
                  v_types[1 + floor(random() * array_length(v_types, 1))::int] || ' - Model #' || LPAD(i::text, 4, '0');
        v_price := (1500 + floor(random() * 85000)) * 1000;
        v_stock := floor(random() * 150)::int;
        
        IF v_stock = 0 THEN
            v_status := 'OUT_OF_STOCK';
        ELSIF v_stock <= 5 THEN
            v_status := 'LOW_STOCK';
        ELSE
            v_status := 'IN_STOCK';
        END IF;

        INSERT INTO products (sku, name, category, price, stock, status, version)
        VALUES (
            'SKU-' || LPAD(i::text, 5, '0'),
            v_name,
            v_cat,
            v_price,
            v_stock,
            v_status,
            1
        ) ON CONFLICT (sku) DO NOTHING;
    END LOOP;
END $$;

-- ==========================================================
-- Verifikasi Data Setelah Import
-- ==========================================================
SELECT 'products' AS tabel, COUNT(*) AS total FROM products
UNION ALL
SELECT 'orders', COUNT(*) FROM orders
UNION ALL
SELECT 'webhook_events', COUNT(*) FROM webhook_events;

-- Cek produk flash sale (ID=1, stock=1 untuk demo concurrency)
SELECT id, sku, name, stock, status FROM products WHERE id = 1;

-- ==========================================================
-- Catatan Konsep: JWT Stateless + Redis JTI Blacklist (Soal #2)
-- ==========================================================
-- Autentikasi bersifat stateless - TIDAK memerlukan tabel users di DB
-- Credential validasi dilakukan hardcoded untuk kebutuhan demo assessment:
--   Email: ratu.syahirah@example.com  Password: password123
--   Email: admin@worksample.com       Password: password123
--   Email: tester@worksample.com      Password: password123
--
-- Di production nyata, alurnya adalah:
--   1. Lookup email di tabel users → bcrypt.CompareHashAndPassword()
--   2. Terbitkan JWT (exp 15 menit) dengan JTI = UUID v4 unik
--   3. Saat logout: masukkan JTI ke Redis BLACKLIST dengan TTL = sisa exp token
--   4. Setiap request ke protected endpoint: cek JTI di Redis → tolak jika ada
--
-- Endpoint Auth:
--   POST /api/v1/auth/login    → terbitkan JWT
--   POST /api/v1/auth/logout   → revoke JTI ke blacklist
--   GET  /api/v1/auth/validate → validasi token + cek blacklist
-- ==========================================================
