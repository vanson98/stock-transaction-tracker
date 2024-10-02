
ALTER TABLE entries
RENAME COLUMN "type" TO "from_type";

ALTER TYPE entry_type RENAME TO entry_from_type;

ALTER TABLE entries 
RENAME TO account_entries;