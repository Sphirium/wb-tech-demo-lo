-- Основной заказ
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    order_uid UUID NOT NULL UNIQUE,
    track_number VARCHAR(64) NOT NULL UNIQUE,
    entry VARCHAR(10) NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature TEXT,
    customer_id VARCHAR(128) NOT NULL,
    delivery_service VARCHAR(64) NOT NULL,
    shardkey VARCHAR(2) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Доставка
CREATE TABLE IF NOT EXISTS delivery (
    order_id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    phone VARCHAR(15) NOT NULL,
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(64) NOT NULL,
    address VARCHAR(64) NOT NULL,
    region VARCHAR(64) NOT NULL,
    email VARCHAR(128),
    FOREIGN KEY (order_id) REFERENCES orders(order_uid) ON DELETE CASCADE
);

-- Оплата
CREATE TABLE IF NOT EXISTS payment (
    transaction UUID PRIMARY KEY,
    order_id UUID NOT NULL UNIQUE,
    request_id VARCHAR(64),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(32) NOT NULL,
    amount INT NOT NULL,
    payment_dt TIMESTAMP NOT NULL,
    bank VARCHAR(20) NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(order_uid) ON DELETE CASCADE
);

-- Товары
CREATE TABLE IF NOT EXISTS items (
    chrt_id BIGINT PRIMARY KEY,
    order_id UUID NOT NULL,
    track_number VARCHAR(64) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    sale INT NOT NULL CHECK(sale BETWEEN 0 AND 100),
    size VARCHAR(16) NOT NULL,
    total_price INT NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(64) NOT NULL,
    status INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(order_uid) ON DELETE CASCADE
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_items_order_id ON items(order_id);
CREATE INDEX IF NOT EXISTS idx_delivery_order_id ON delivery(order_id);
CREATE INDEX IF NOT EXISTS idx_payment_order_id ON payment(order_id);