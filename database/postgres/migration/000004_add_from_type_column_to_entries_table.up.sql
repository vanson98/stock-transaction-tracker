CREATE TYPE entry_from_type AS ENUM ('TM','IT');

ALTER TABLE entries
ADD COLUMN from_type entry_from_type NOT NULL;    

ALTER TABLE entries
RENAME TO account_entries;