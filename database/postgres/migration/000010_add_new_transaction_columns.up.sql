CREATE TYPE trade_type AS ENUM ('SELL','BUY');
CREATE TYPE transaction_status AS ENUM ('COMPLETED','INCOMPLETED');


ALTER TABLE transactions ADD COLUMN ticker VARCHAR(20) NOT NULL;
ALTER TABLE transactions ADD COLUMN trading_date timestamp NOT NULL;
ALTER TABLE transactions ADD COLUMN trade trade_type NOT NULL;
ALTER TABLE transactions ADD COLUMN volume INT  NOT NULL;
ALTER TABLE transactions ADD COLUMN order_price BIGINT  NOT NULL;
ALTER TABLE transactions ADD COLUMN match_volume INT  NOT NULL;
ALTER TABLE transactions ADD COLUMN match_price BIGINT NOT NULL;
ALTER TABLE transactions ADD COLUMN match_value BIGINT NOT NULL;
ALTER TABLE transactions ADD COLUMN fee INT  NOT NULL;
ALTER TABLE transactions ADD COLUMN tax INT NOT NULL;
ALTER TABLE transactions ADD COLUMN "cost" BIGINT NOT NULL;
ALTER TABLE transactions ADD COLUMN cost_of_goods_sold BIGINT NOT NULL;
ALTER TABLE transactions ADD COLUMN "return" BIGINT NOT NULL;
ALTER TABLE transactions ADD COLUMN "status" transaction_status NOT NULL;