ALTER TABLE coins_networks ADD COLUMN priority INT NOT NULL DEFAULT 0
;

UPDATE coins_networks SET priority = 10 WHERE coin_id = 'usdt' and network_id = 'trc20'
;
UPDATE coins_networks SET priority = 10 WHERE coin_id = 'usdc' and network_id = 'trc20'
;

UPDATE coins_networks SET priority = 20 WHERE coin_id = 'usdt' and network_id = 'bep20'
;
UPDATE coins_networks SET priority = 20 WHERE coin_id = 'usdc' and network_id = 'bep20'
;

UPDATE coins_networks SET priority = 30 WHERE coin_id = 'usdt' and network_id = 'avax'
;
UPDATE coins_networks SET priority = 30 WHERE coin_id = 'usdc' and network_id = 'avax'
;

UPDATE coins_networks SET priority = 10 WHERE coin_id != 'usdt' and coin_id != 'usdc'
;
