

select * from schema_migrations
select * from investments
select * from accounts
select * from transactions
SELECT * FROM entries
where account_id=200

update schema_migrations
SET "version" = 11, dirty=FALSE
=========



select 
a.id, 
a.channel_name, 
a.balance,
a.currency,
a."owner",
SUM(
	case
	WHEN amount > 0 and e.type='TM' then amount
	ELSE 0
	END
	) as transfer_in,
SUM(
	CASE 
	WHEN amount < 0 and e.type='TM' THEN amount
	ELSE 0 
	END
) AS transfer_out
from accounts as a
left join entries as e on a.id = e.account_id 
where a.id = 165
GROUP BY a.id,  a.channel_name, a.balance, a.currency, a."owner"
LIMIT 1;



===============================================
select * from entries
INSERT INTO entries(accou)

INSERT INTO transactions 
