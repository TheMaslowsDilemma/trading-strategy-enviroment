-- "initial" Up Migration
-- executed when this migration is applied


CREATE TABLE IF NOT EXISTS users (
	id	BIGSERIAL PRIMARY KEY,
	name	TEXT UNIQUE NOT NULL,
	trader_id	BIGINT UNIQUE NOT NULL,
	password_hash	TEXT NOT NULL,
	data_subscriptions	JSONB,
	created_at	TIMESTAMPTZ DEFAULT NOW()
);