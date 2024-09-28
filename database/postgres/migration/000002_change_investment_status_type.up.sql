CREATE TYPE investment_status AS ENUM ('inactive','active','buyout');

ALTER TABLE investments
ALTER COLUMN status TYPE investment_status USING status::investment_status,
ALTER COLUMN status SET NOT NULL,
ALTER COLUMN status SET DEFAULT 'inactive';