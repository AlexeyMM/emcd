UPDATE coins_networks
SET decimals = 8;

UPDATE coins_networks
SET hash_divisor_power_of_ten = 12
WHERE coin_id = 'btc'
  AND network_id = 'btc';

UPDATE coins_networks
SET hash_divisor_power_of_ten = 12
WHERE coin_id = 'bch'
  AND network_id = 'bch';

UPDATE coins_networks
SET hash_divisor_power_of_ten = 9
WHERE coin_id = 'ltc'
  AND network_id = 'ltc';

UPDATE coins_networks
SET hash_divisor_power_of_ten = 12
WHERE coin_id = 'dash'
  AND network_id = 'dash';

UPDATE coins_networks
SET hash_divisor_power_of_ten = 9
WHERE coin_id = 'doge'
  AND network_id = 'doge';

UPDATE coins_networks
SET hash_divisor_power_of_ten = 12
WHERE coin_id = 'etc'
  AND network_id = 'etc';

UPDATE coins_networks
SET withdrawals_disabled_description = NULL
WHERE coin_id IN ('eth', 'etc');
