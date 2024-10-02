ALTER TABLE account_entries
ALTER COLUMN account_id DROP NOT  NULL;

ALTER TABLE investments
ALTER COLUMN account_id DROP NOT  NULL;

ALTER TABLE transactions 
ALTER COLUMN investment_id DROP NOT  NULL;

ALTER TABLE account_entries
RENAME COLUMN account_id TO acount_id;