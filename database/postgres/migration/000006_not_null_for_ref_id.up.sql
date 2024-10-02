ALTER TABLE account_entries
RENAME COLUMN acount_id TO account_id;

ALTER TABLE account_entries
ALTER COLUMN account_id SET NOT NULL;

ALTER TABLE investments
ALTER COLUMN account_id SET NOT NULL;

ALTER TABLE transactions 
ALTER COLUMN investment_id SET NOT NULL;