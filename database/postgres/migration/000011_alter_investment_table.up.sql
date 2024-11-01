ALTER TABLE investments RENAME COLUMN stock_code TO ticker;
ALTER TABLE investments RENAME COLUMN total_buy_amount TO buy_volume;
ALTER TABLE investments RENAME COLUMN total_money_buy TO buy_value;
ALTER TABLE investments RENAME COLUMN total_sell_amount TO sell_volume;
ALTER TABLE investments RENAME COLUMN total_money_sell TO sell_value;
ALTER TABLE investments ADD COLUMN fee INT NOT NULL;
ALTER TABLE investments ADD COLUMN tax INT NOT NULL;
ALTER TABLE investments ADD COLUMN updated_date TIMESTAMP;
