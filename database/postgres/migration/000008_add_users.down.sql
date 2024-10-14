ALTER TABLE accounts DROP CONSTRAINT "owner_channel_key";
ALTER TABLE accounts DROP CONSTRAINT "accounts_owner_fkey";
DROP TABLE users;