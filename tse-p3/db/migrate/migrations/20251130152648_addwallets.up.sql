-- "add-wallets" Up Migration
-- executed when this migration is applied


CREATE TABLE IF NOT EXISTS wallets (
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT NOT NULL,
    symbol    TEXT NOT NULL,
    amount    TEXT NOT NULL,

    -- FK: when a user is deleted, delete their wallets
    CONSTRAINT wallets_user_fk
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    -- Unique per user per symbol
    CONSTRAINT wallets_user_symbol_uniq
        UNIQUE (user_id, symbol)
);