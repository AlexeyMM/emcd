ALTER TABLE coins
	ADD COLUMN legacy_coin_id INT NOT NULL DEFAULT 0;

UPDATE coins
SET legacy_coin_id = 1
WHERE id = 'btc';

UPDATE coins
SET legacy_coin_id = 2
WHERE id = 'bch';

UPDATE coins
SET legacy_coin_id = 4
WHERE id = 'ltc';

UPDATE coins
SET legacy_coin_id = 5
WHERE id = 'dash';

UPDATE coins
SET legacy_coin_id = 6
WHERE id = 'eth';

UPDATE coins
SET legacy_coin_id = 7
WHERE id = 'etc';

UPDATE coins
SET legacy_coin_id = 8
WHERE id = 'doge';

UPDATE coins
SET legacy_coin_id = 9
WHERE id = 'bnb';

UPDATE coins
SET legacy_coin_id = 10
WHERE id = 'usdt';

UPDATE coins
SET legacy_coin_id = 11
WHERE id = 'usdc';

UPDATE coins
SET legacy_coin_id = 12
WHERE id = 'trx';

UPDATE coins
SET legacy_coin_id = 13
WHERE id = 'ton';

UPDATE coins
SET legacy_coin_id = 14
WHERE id = 'kas';
