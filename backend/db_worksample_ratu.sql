-- ==========================================================
-- Database Schema: db_worksample_ratu
-- PostgreSQL Schema for Senior Fullstack Work Sample
-- Candidate: Ratu Syahirah Khairunnisa
-- ==========================================================

-- 1. Create Database (Run this query first in DBeaver/PostgreSQL if not exists)
-- CREATE DATABASE db_worksample_ratu;

-- 2. Connect to database
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
    FOR i IN 11..5000 LOOP
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
