ALTER TABLE coins_networks
	ADD COLUMN minpay_mining DOUBLE PRECISION NOT NULL DEFAULT 0;

UPDATE coins_networks
SET minpay_mining = 0.00001
WHERE coin_id = 'btc'
  AND network_id = 'btc';

UPDATE coins_networks
SET minpay_mining = 0.00001
WHERE coin_id = 'bch'
  AND network_id = 'bch';

UPDATE coins_networks
SET minpay_mining = 0.00001
WHERE coin_id = 'ltc'
  AND network_id = 'ltc';

UPDATE coins_networks
SET minpay_mining = 0.00001
WHERE coin_id = 'dash'
  AND network_id = 'dash';

UPDATE coins_networks
SET minpay_mining = 1
WHERE coin_id = 'doge'
  AND network_id = 'doge';

UPDATE coins_networks
SET minpay_mining = 0.1
WHERE coin_id = 'etc'
  AND network_id = 'etc';
