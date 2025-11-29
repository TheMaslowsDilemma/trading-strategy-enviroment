-- "add-user-activity" Down Migration
-- executed when this migration is rolled back

alter table users
drop column active;

