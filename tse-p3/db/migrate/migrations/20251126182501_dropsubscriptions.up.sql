-- "drop-subscriptions" Up Migration
-- executed when this migration is applied

ALTER TABLE users
DROP COLUMN IF EXISTS data_subscriptions;
