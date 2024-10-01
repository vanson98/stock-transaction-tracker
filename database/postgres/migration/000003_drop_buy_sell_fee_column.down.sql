
ALTER TABLE accounts
ADD COLUMN "buy_fee" numeric(10,3) NOT NULL DEFAULT 0,
ADD COLUMN "sell_free" numeric(10,3) NOT NULL DEFAULT 0;


  