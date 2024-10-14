

select * from schema_migrations
select * from investments
select * from accounts
ORDER BY id desc;
SELECT * FROM entries
where account_id=10;

update schema_migrations
SET "version" = 7, dirty=FALSE
=========
DELETE FROM accounts 
DELETE FROM entries

CREATE TABLE users(
	username VARCHAR PRIMARY KEY,
	hashed_password VARCHAR NOT NULL,
	full_name VARCHAR NOT NULL,
	email varchar UNIQUE NOT NULL,
	password_changed_at TIMESTAMPTZ NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
	created_at TIMESTAMPTZ NOT NULL DEFAULT(now())
);

ALTER TABLE accounts ADD FOREIGN KEY ("owner") REFERENCES users (username);
-- CREATE UNIQUE INDEX owner_channel_key ON public.accounts USING btree (owner, channel_name);
ALTER TABLE accounts ADD CONSTRAINT "owner_channel_key" UNIQUE ("owner",channel_name);
========
ALTER TABLE accounts DROP CONSTRAINT "owner_channel_key";
ALTER TABLE accounts DROP CONSTRAINT accounts_owner_fkey;
DROP TABLE users;
