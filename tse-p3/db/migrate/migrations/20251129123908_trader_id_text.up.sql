-- "trader_id_text" Up Migration
-- executed when this migration is applied


ALTER TABLE users 
ALTER COLUMN trader_id TYPE TEXT 
USING trader_id::text;

