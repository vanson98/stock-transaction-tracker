CREATE TABLE "investments" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint,
  "stock_code" varchar NOT NULL,
  "company_name" varchar,
  "total_buy_amount" int NOT NULL DEFAULT 0,
  "total_money_buy" numeric(10,3) NOT NULL,
  "capital_cost" numeric(10,3) NOT NULL,
  "market_price" numeric(10,3) NOT NULL,
  "total_sell_amount" int NOT NULL DEFAULT 0,
  "total_money_sell" numeric(10,3) NOT NULL DEFAULT 0,
  "current_volume" int NOT NULL DEFAULT 0,
  "description" varchar,
  "status" varchar
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "channel_name" varchar NOT NULL,
  "owner" varchar NOT NULL,
  "balance" numeric(10,3) NOT NULL,
  "buy_fee" numeric(10,3) NOT NULL,
  "sell_free" numeric(10,3) NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "acount_id" bigint,
  "amount" numeric(10,3) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "investment_id" bigint,
  "capital_cost" numeric(10,3) NOT NULL,
  "price" numeric(10,3) NOT NULL,
  "amount" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "type" varchar NOT NULL,
  "fee" numeric(10,3) NOT NULL
);

CREATE INDEX ON "investments" ("stock_code");

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("acount_id");

CREATE INDEX ON "transactions" ("investment_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be possitive or negative';

COMMENT ON COLUMN "transactions"."amount" IS 'must be possitive';

ALTER TABLE "investments" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("acount_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("investment_id") REFERENCES "investments" ("id");
