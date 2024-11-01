ALTER TABLE transactions ADD COLUMN capital_cost int8;
ALTER TABLE transactions ADD COLUMN price int8;
ALTER TABLE transactions ADD COLUMN amount int4;
ALTER TABLE transactions ADD COLUMN created_at timestamptz;
ALTER TABLE transactions ADD COLUMN "type" VARCHAR;
ALTER TABLE transactions ADD COLUMN fee int8;