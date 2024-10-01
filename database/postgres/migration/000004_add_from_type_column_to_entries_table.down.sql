ALTER TABLE account_entries
RENAME TO entries;

ALTER TABLE entries
DROP COLUMN from_type;

DROP TYPE entry_from_type;
