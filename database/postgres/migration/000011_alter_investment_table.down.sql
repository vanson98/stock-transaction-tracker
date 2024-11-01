ALTER TABLE investments RENAME COLUMN ticker  TO stock_code;
ALTER TABLE investments RENAME COLUMN buy_volume  TO total_buy_amount;
ALTER TABLE investments RENAME COLUMN buy_value  TO total_money_buy;
ALTER TABLE investments RENAME COLUMN sell_volume  TO total_sell_amount;
ALTER TABLE investments RENAME COLUMN sell_value  TO total_money_sell;
ALTER TABLE investments DROP COLUMN fee;
ALTER TABLE investments DROP COLUMN tax;
ALTER TABLE investments DROP COLUMN updated_date;