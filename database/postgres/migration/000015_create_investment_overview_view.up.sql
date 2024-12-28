CREATE VIEW investment_overview
AS
SELECT  
i.id, a.id AS account_id ,a.channel_name, 
i.ticker, i.buy_value, i.buy_volume, i.capital_cost, 
i.current_volume, i.market_price, 
i.sell_value, i.sell_volume, i.fee, i.tax, i.status, 
(
	CASE 
	WHEN i.status = 'active' THEN 
		TRIM_SCALE(TRUNC((CAST(i.market_price-i.capital_cost AS NUMERIC) /i.capital_cost)*100, 2)) 	
	WHEN i.status = 'sellout' THEN
		TRUNC(CAST(i.sell_value - i.buy_value - i.fee- i.tax AS NUMERIC)*100/i.buy_value, 2)
	ELSE 0 END 	
)::NUMERIC AS profit
FROM investments AS i
JOIN accounts AS a ON i.account_id = a.id;
