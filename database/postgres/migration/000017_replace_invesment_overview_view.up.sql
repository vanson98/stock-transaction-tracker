DROP VIEW investment_overview;
CREATE VIEW investment_overview
AS
 SELECT i.id,
    a.id AS account_id,
    a.channel_name,
    i.ticker,
    i.buy_value,
    i.buy_volume,
    i.capital_cost,
    i.current_volume,
    i.market_price,
    i.sell_value,
    i.sell_volume,
    i.fee,
    i.tax,
    i.status,
        CASE
            WHEN i.status = 'active'::investment_status THEN trim_scale(trunc((i.market_price - i.capital_cost)::numeric / i.capital_cost::numeric * 100::numeric, 2))::float
            WHEN i.status = 'sellout'::investment_status THEN trunc((i.sell_value - i.buy_value - i.fee - i.tax)::numeric * 100::numeric / i.buy_value::numeric, 2)::float
            ELSE 0::float
        END AS profit
   FROM investments i
     JOIN accounts a ON i.account_id = a.id;