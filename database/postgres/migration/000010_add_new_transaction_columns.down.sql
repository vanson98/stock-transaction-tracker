ALTER TABLE transactions DROP COLUMN "status";
ALTER TABLE transactions DROP COLUMN "return";
ALTER TABLE transactions DROP COLUMN cost_of_goods_sold;
ALTER TABLE transactions DROP COLUMN "cost";
ALTER TABLE transactions DROP COLUMN tax;
ALTER TABLE transactions DROP COLUMN fee;
ALTER TABLE transactions DROP COLUMN match_value;
ALTER TABLE transactions DROP COLUMN match_price;
ALTER TABLE transactions DROP COLUMN match_volume;
ALTER TABLE transactions DROP COLUMN order_price;
ALTER TABLE transactions DROP COLUMN volume;
ALTER TABLE transactions DROP COLUMN trade;
ALTER TABLE transactions DROP COLUMN trading_date;
ALTER TABLE transactions DROP COLUMN ticker;

DROP TYPE transaction_status;
DROP TYPE trade_type;