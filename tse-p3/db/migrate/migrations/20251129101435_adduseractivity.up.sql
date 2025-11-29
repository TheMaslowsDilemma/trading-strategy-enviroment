-- "add-user-activity" Up Migration
-- executed when this migration is applied

ALTER TABLE users
ADD COLUMN active BOOLEAN NOT NULL DEFAULT false;
