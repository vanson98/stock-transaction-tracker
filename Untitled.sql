

select * from schema_migrations
select * from investments
select * from accounts
select * from transactions
ORDER BY id desc;
SELECT * FROM entries
where account_id=10;

update schema_migrations
SET "version" = 11, dirty=FALSE
=========
ALTER TABLE investments RENAME COLUMN stock_code TO ticker;
ALTER TABLE investments RENAME COLUMN total_buy_amount TO buy_volume;
ALTER TABLE investments RENAME COLUMN total_money_buy TO buy_value;
ALTER TABLE investments RENAME COLUMN total_sell_amount TO sell_volume;
ALTER TABLE investments RENAME COLUMN total_money_sell TO sell_value;
ALTER TABLE investments ADD COLUMN fee INT NOT NULL;
ALTER TABLE investments ADD COLUMN tax INT NOT NULL;
===============================================

ALTER TABLE investments RENAME COLUMN ticker  TO stock_code;
ALTER TABLE investments RENAME COLUMN buy_volume  TO total_buy_amount;
ALTER TABLE investments RENAME COLUMN buy_value  TO total_money_buy;
ALTER TABLE investments RENAME COLUMN sell_volume  TO total_sell_amount;
ALTER TABLE investments RENAME COLUMN sell_value  TO total_money_sell;
ALTER TABLE investments DROP COLUMN fee;
ALTER TABLE investments DROP COLUMN tax;

