-- "add-wallets" Down Migration
-- executed when this migration is rolled back


ALTER TABLE wallets
    DROP CONSTRAINT IF EXISTS wallets_user_symbol_uniq;

ALTER TABLE wallets
    DROP CONSTRAINT IF EXISTS wallets_user_fk;

DROP TABLE IF EXISTS wallets;