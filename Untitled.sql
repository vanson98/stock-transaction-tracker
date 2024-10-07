

select * from schema_migrations
select * from investments
select * from accounts
ORDER BY id desc;
SELECT * FROM entries;


update schema_migrations
SET "version" = 7, dirty=FALSE
=========
BEGIN;

INSERT INTO entries(account_id,amount,"type")
VALUES(2,10,'IT')
RETURNING *;

SELECT id, channel_name, owner, balance, currency, created_at FROM accounts
WHERE id=2 LIMIT 1
FOR UPDATE;

UPDATE accounts
SET balance = 510
WHERE id = 2
RETURNING *;

ROLLBACK;
========

ALTER TABLE entries 
RENAME TO account_entries;

ALTER TABLE entries
RENAME COLUMN "type" TO "from_type";

ALTER TYPE entry_type RENAME TO entry_from_type;