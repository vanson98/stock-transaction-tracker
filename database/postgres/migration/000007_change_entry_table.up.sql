ALTER TABLE account_entries
RENAME TO entries;

ALTER TABLE entries
RENAME COLUMN from_type TO "type";

ALTER TYPE entry_from_type RENAME TO entry_type;